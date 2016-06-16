package api

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarmkit/api"
	ct "golang.org/x/net/context"
)

// GET /services
func listService(c *context, w http.ResponseWriter, r *http.Request) {
	sresp, err := c.swarmkitAPI.ListServices(ct.TODO(), &api.ListServicesRequest{})
	if err != nil {
		c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err.Error()})
		return
	}

	c.render.JSON(w, http.StatusOK, sresp.Services)
}

// POST /services/create
// {
//    name:"redis",
//    image:"redis:3.0.5",
//    labels:{"com.docker.test":"test"},        // service label (key=value)
//    mode:"",                                  // one of replicated, global
//    replicas: 1,                              // number of replicas for the service (only works in replicated service mode)
//    image: "redis:3.0.5",                     // container image
//    args: [],                                 // container args
//    env: [],                                  // container env
//    ports: [],                                // ports
//    network:"",                               // network name
//    memory-reservation: "",                   // amount of reserved memory (e.g. 512m)
//    memory-limit: "",                         // memory limit (e.g. 512m)
//    cpu-reservation:"",                       // number of CPU cores reserved (e.g. 0.5)
//    cpu-limit:"",                             // CPU cores limit (e.g. 0.5)
//    update-parallelism:0,                     // task update parallelism (0 = all at once)
//    update-delay:"0s",                        // delay between task updates (0s = none)
//    restart-condition:"any",                  // condition to restart the task (any, failure, none)
//    restart-delay:"5s",                       // delay between task restarts
//    restart-max-attempts:0,                   // maximum number of restart attempts (0 = unlimited)
//    restart-window:"0s",                      // time window to evaluate restart attempts (0 = unbound)
//    constraint:[],                            // Placement constraint (node.labels.key==value)
//    bind:[],                                  // define a bind mount
//    volume:[],                                // define a volume mount
// }
func createService(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		cspec = &createSpec{}
	)
	if err = DecoderRequest(r, cspec); err != nil {
		err = fmt.Errorf("Parse params for create service error:%v", err)
		log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
		c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
		return
	}

	if len(strings.TrimSpace(cspec.Name)) == 0 || len(strings.TrimSpace(cspec.Image)) == 0 {
		err = fmt.Errorf("--name and --image are mandatory")
		log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
		c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
		return
	}

	spec := &api.ServiceSpec{
		Mode: &api.ServiceSpec_Replicated{
			Replicated: &api.ReplicatedService{
				Replicas: 1,
			},
		},
		Task: api.TaskSpec{
			Runtime: &api.TaskSpec_Container{
				Container: &api.ContainerSpec{},
			},
		},
	}

	if err = merge(cspec, spec, c.swarmkitAPI); err != nil {
		log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
		c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
		return
	}

	var csResp *api.CreateServiceResponse
	if csResp, err = c.swarmkitAPI.CreateService(ct.TODO(), &api.CreateServiceRequest{Spec: spec}); err != nil {
		log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
		c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
		return
	}

	c.render.JSON(w, http.StatusOK, csResp.Service)
}
