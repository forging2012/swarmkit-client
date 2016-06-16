package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func errResponse(w http.ResponseWriter, r *http.Request, err error, c *context) {
	log.WithFields(log.Fields{"method": r.Method, "route": r.RequestURI}).Errorln(err)
	c.render.JSON(w, http.StatusBadRequest, map[string]interface{}{"msg": err})
}
