package api

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/gorilla/mux"
	"github.com/shenshouer/swarmkit-client/swarmkit"
	ct "golang.org/x/net/context"
)

// GET /networks
func listNetworks(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err              error
		listNetworksResp *api.ListNetworksResponse
	)
	if listNetworksResp, err = c.swarmkitAPI.ListNetworks(ct.TODO(), &api.ListNetworksRequest{}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, listNetworksResp.Networks)
}

// GET /networks/{networkid:.*}
func inspectNetworks(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		network   *api.Network
		networkid = mux.Vars(r)["networkid"]
	)
	if network, err = swarmkit.GetNetwork(ct.TODO(), c.swarmkitAPI, networkid); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, network)
}

// POST /networks/creat
//{
//	name: "",
//	driver: "",
//	opts: {},
//	ipam_driver: "",
//	subnet: [],
//	gateway:[],
//	ip_range:[],
//
//}
func createNetworks(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		driver   *api.Driver
		ipamOpts *api.IPAMOptions
		nwInfo   = &networkInfo{}
	)

	if err = DecoderRequest(r, nwInfo); err != nil {
		errResponse(w, r, err, c)
		return
	}
	// parse api.Driver
	if len(strings.TrimSpace(nwInfo.Name)) == 0 {
		err = errors.New("name is required")
		errResponse(w, r, err, c)
		return
	}

	if len(strings.TrimSpace(nwInfo.Driver)) > 1 {
		driver = new(api.Driver)
		driver.Name = nwInfo.Name

		if len(nwInfo.Opts) > 0 {
			driver.Options = nwInfo.Opts

		}
	}

	// parse api.IPAMOptions
	if ipamOpts, err = processIPAMOptions(nwInfo); err != nil {
		errResponse(w, r, err, c)
		return
	}

	spec := &api.NetworkSpec{
		Annotations: api.Annotations{
			Name: nwInfo.Name,
		},
		DriverConfig: driver,
		IPAM:         ipamOpts,
	}

	var cnResp *api.CreateNetworkResponse
	if cnResp, err = c.swarmkitAPI.CreateNetwork(ct.TODO(), &api.CreateNetworkRequest{Spec: spec}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, cnResp.Network.ID)
}

// DELETE /networks/{networkid:.*}
func removeNetworks(c *context, w http.ResponseWriter, r *http.Request) {}

type networkInfo struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Opts       map[string]string `json:"opts"`
	IpamDriver string            `json:"ipam_driver"`
	Subnet     []string          `json:"subnet"`
	Gateway    []string          `json:"gateway"`
	IPRange    []string          `json:"ip_range"`
}

func processIPAMOptions(nwInfo *networkInfo) (ipamOpts *api.IPAMOptions, err error) {
	if len(strings.TrimSpace(nwInfo.IpamDriver)) > 1 {
		ipamOpts = &api.IPAMOptions{
			Driver: &api.Driver{
				Name: nwInfo.IpamDriver,
			},
		}
	}

	subnets := nwInfo.Subnet
	if subnets == nil || len(subnets) == 0 {
		return
	}

	ipamConfigs := make([]*api.IPAMConfig, 0, len(subnets))
	for _, s := range nwInfo.Subnet {
		_, ipNet, err := net.ParseCIDR(s)
		if err != nil {
			return nil, err
		}

		family := api.IPAMConfig_IPV6
		if ipNet.IP.To4() != nil {
			family = api.IPAMConfig_IPV4
		}

		var gateway string
		gateways := nwInfo.Gateway
		for i, g := range gateways {
			if ipNet.Contains(net.ParseIP(g)) {
				gateways = append(gateways[:i], gateways[i+1:]...)
				gateway = g
				break
			}
		}

		var iprange string
		ranges := nwInfo.IPRange
		for i, r := range ranges {
			_, rangeNet, err := net.ParseCIDR(r)
			if err != nil {
				return nil, err
			}

			if ipNet.Contains(rangeNet.IP) {
				ranges = append(ranges[:i], ranges[i+1:]...)
				iprange = r
				break
			}
		}

		ipamConfigs = append(ipamConfigs, &api.IPAMConfig{
			Family:  family,
			Subnet:  s,
			Gateway: gateway,
			Range:   iprange,
		})
	}

	if ipamOpts == nil {
		ipamOpts = &api.IPAMOptions{}
	}

	ipamOpts.Configs = ipamConfigs
	return
}
