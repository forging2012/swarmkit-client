package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/shenshouer/swarmkit-client/swarmkit"
	"golang.org/x/crypto/bcrypt"
	ct "golang.org/x/net/context"
)

// GET /clusters
func listClusters(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err             error
		listClusterResp *api.ListClustersResponse
	)

	if listClusterResp, err = c.swarmkitAPI.ListClusters(ct.TODO(), &api.ListClustersRequest{}); err != nil {
		errResponse(w, r, err, c)
		return
	}
	c.render.JSON(w, http.StatusOK, listClusterResp.Clusters)
}

// GET /clusters/{clusterid:.*}
func inspectClusters(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		cluster   *api.Cluster
		clusterid = mux.Vars(r)["clusterid"]
	)

	if cluster, err = swarmkit.GetCluster(ct.TODO(), c.swarmkitAPI, clusterid); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, cluster)
}

// POST /clusters/{clusterid:.*}/update
//{
//  autoaccept: [],     // Roles to automatically issue certificates for
//  secret:[],          // Secret required to join the cluster
//  taskhistory:0,      // Number of historic task entries to retain per slot or node
//  certexpiry:"",      // Duration node certificates will be valid for
//  heartbeatperiod:""  // Duration Period when heartbeat is expected to receive from agent
//}
func updateClusters(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		cluster           *api.Cluster
		clusterid         = mux.Vars(r)["clusterid"]
		updateClusterResp *api.UpdateClusterResponse
		cInfo             = &struct {
			Autoaccept      []string `json:"autoaccept"`
			Secret          []string `json:"secret"`
			Taskhistory     int64    `json:"taskhistory"`
			Certexpiry      string   `json:"certexpiry"`
			Heartbeatperiod string   `json:"heartbeatperiod"`
		}{}
	)

	if err = DecoderRequest(r, cInfo); err != nil {
		errResponse(w, r, err, c)
		return
	}

	if cluster, err = swarmkit.GetCluster(ct.TODO(), c.swarmkitAPI, clusterid); err != nil {
		errResponse(w, r, err, c)
		return
	}

	spec := &cluster.Spec
	if len(cInfo.Autoaccept) > 0 {
		// We are getting a whitelist, so make all of the autoaccepts false
		for _, policy := range spec.AcceptancePolicy.Policies {
			policy.Autoaccept = false
		}
		autoaccept := cInfo.Autoaccept
		// For each of the roles handed to us by the client, make them true
		for _, role := range autoaccept {
			// Convert the role into a proto role
			apiRole, err := ca.FormatRole("swarm-" + role)
			if err != nil {
				err = fmt.Errorf("unrecognized role %s", role)
				errResponse(w, r, err, c)
				return
			}
			// Attempt to find this role inside of the current policies
			found := false
			for _, policy := range spec.AcceptancePolicy.Policies {
				if policy.Role == apiRole {
					// We found a matching policy, let's update it
					policy.Autoaccept = true
					found = true
				}
			}
			// We didn't find this policy, create it
			if !found {
				newPolicy := &api.AcceptancePolicy_RoleAdmissionPolicy{
					Role:       apiRole,
					Autoaccept: true,
				}
				spec.AcceptancePolicy.Policies = append(spec.AcceptancePolicy.Policies, newPolicy)
			}
		}
	}

	if len(cInfo.Secret) > 0 {
		// Using the defaut bcrypt cost
		hashedSecret, err := bcrypt.GenerateFromPassword([]byte(cInfo.Secret[0]), 0)
		if err != nil {
			errResponse(w, r, err, c)
		}
		for _, policy := range spec.AcceptancePolicy.Policies {
			policy.Secret = &api.AcceptancePolicy_RoleAdmissionPolicy_HashedSecret{
				Data: hashedSecret,
				Alg:  "bcrypt",
			}
		}
	}

	if len(cInfo.Certexpiry) > 1 {
		duration, err := ParseString(cInfo.Certexpiry)
		if err != nil {
			errResponse(w, r, err, c)
			return
		}
		ceProtoPeriod := ptypes.DurationProto(duration.Duration())
		spec.CAConfig.NodeCertExpiry = ceProtoPeriod

	}

	if cInfo.Taskhistory > 0 {
		spec.Orchestration.TaskHistoryRetentionLimit = cInfo.Taskhistory
	}

	if len(strings.TrimSpace(cInfo.Heartbeatperiod)) > 1 {
		duration, err := ParseString(cInfo.Certexpiry)
		if err != nil {
			errResponse(w, r, err, c)
			return
		}

		spec.Dispatcher.HeartbeatPeriod = uint64(duration.Duration())
	}

	if updateClusterResp, err = c.swarmkitAPI.UpdateCluster(ct.TODO(), &api.UpdateClusterRequest{
		ClusterID:      cluster.ID,
		ClusterVersion: &cluster.Meta.Version,
		Spec:           spec,
	}); err != nil {
		errResponse(w, r, err, c)
		return
	}

	c.render.JSON(w, http.StatusOK, updateClusterResp.Cluster.ID)
}
