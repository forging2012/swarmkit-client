package api

import (
	"fmt"
	"net/http"

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

// GET /nodes/{nodeid:.*}
func inspectNode(c *context, w http.ResponseWriter, r *http.Request) {
	var nodeid = mux.Vars(r)["nodeid"]
	fmt.Println(nodeid)
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
