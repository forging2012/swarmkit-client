package api

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarmkit/api"
	"github.com/gorilla/mux"
	"github.com/shenshouer/swarmkit-client/swarmkit"
	ct "golang.org/x/net/context"
)

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
		err    error
		node   *api.Node
		nodeid = mux.Vars(r)["nodeid"]
		allStr = r.URL.Query().Get("all")
		all    = false
	)
	if len(strings.TrimSpace(allStr)) != 0 && allStr == "1" {
		all = true
	}

	if node, err = swarmkit.GetNode(ct.TODO(), c.swarmkitAPI, nodeid); err != nil {
		errResponse(w, r, err, c)
		return
	}

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
