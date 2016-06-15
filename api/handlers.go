package api

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarmkit/api"
	"github.com/gorilla/mux"
	ct "golang.org/x/net/context"
)

// Node OP

// GET /nodes
func listNodes(c *context, w http.ResponseWriter, r *http.Request) {
	lsNodeRes, err := c.swarmkitAPI.ListNodes(ct.TODO(), &api.ListNodesRequest{})
	if err != nil {
		log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
		c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
		return
	}
	c.render.JSON(w, http.StatusOK, lsNodeRes.Nodes)
}

// GET /nodes/{nodeid:.*}?all=1
//    all:0 only display running
//		  1 display all
//	  default 0
func inspectNode(c *context, w http.ResponseWriter, r *http.Request) {
	var (
		node   *api.Node
		nodeid = mux.Vars(r)["nodeid"]
		allStr = r.URL.Query().Get("all")
		ctx    = ct.TODO()
		all    = false
	)
	if len(strings.TrimSpace(allStr)) != 0 && allStr == "1" {
		all = true
	}

	// GetNode to match via full ID.
	rg, err := c.swarmkitAPI.GetNode(ctx, &api.GetNodeRequest{NodeID: nodeid})
	if err != nil {
		// If any error (including NotFound), ListServices to match via ID prefix and full name.
		rl, err := c.swarmkitAPI.ListNodes(ctx,
			&api.ListNodesRequest{
				Filters: &api.ListNodesRequest_Filters{
					Names: []string{nodeid},
				},
			},
		)
		if err != nil {
			log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
			c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
			return
		}

		if len(rl.Nodes) == 0 {
			err = fmt.Errorf("node %s not found", nodeid)
			log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
			c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
			return
		}

		if l := len(rl.Nodes); l > 1 {
			err = fmt.Errorf("node %s is ambiguous (%d matches found)", nodeid, l)
			log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
			c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
			return
		}

		node = rl.Nodes[0]
	}
	node = rg.Node

	// TODO(aluzzardi): This should be implemented as a ListOptions filter.
	ltRes, err := c.swarmkitAPI.ListTasks(ct.TODO(), &api.ListTasksRequest{})
	if err != nil {
		log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
		c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
		return
	}

	tasks := []*api.Task{}
	for _, t := range ltRes.Tasks {
		if t.NodeID == node.ID {
			if !all && t.DesiredState > api.TaskStateRunning {
				continue
			}

			tasks = append(tasks, t)
		}
	}

	c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"node": node, "tasks": tasks})
}

// POST /nodes/accept
func acceptNode(c *context, w http.ResponseWriter, r *http.Request) {

}

// DELETE /nodes/{nodeid:.*}
func removeNode(c *context, w http.ResponseWriter, r *http.Request) {}

// POST /nodes/{nodeid:.*}/activate
func activateNode(c *context, w http.ResponseWriter, r *http.Request) {}

func optionsHandler(c *context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Emit an HTTP error and log it.
func httpError(w http.ResponseWriter, err string, status int) {
	log.WithField("status", status).Errorf("HTTP error: %v", err)
	http.Error(w, err, status)
}
