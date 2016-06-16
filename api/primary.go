package api

import (
	"crypto/tls"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarmkit/api"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Primary router context, used by handlers.
type context struct {
	swarmkitAPI   api.ControlClient
	eventsHandler *eventsHandler
	tlsConfig     *tls.Config
	render        *render.Render
	// apiVersion    string
	// statusHandler StatusHandler
}

type handler func(c *context, w http.ResponseWriter, r *http.Request)

var routes = map[string]map[string]handler{
	http.MethodGet: {
		"/nodes":                   listNodes,
		"/nodes/{nodeid:.*}":       inspectNode,
		"/services":                listService,
		"/services/{serviceid:.*}": inspectService,
	},
	http.MethodPost: {
		"/nodes/accept":               acceptNode,
		"/nodes/{nodeid:.*}/activate": activateNode,
		"/services/create":            createService,
	},
	http.MethodDelete: {
		"/nodes/{nodeid:.*}": removeNode,
	},
}

func writeCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

// NewPrimary creates a new API router.
func NewPrimary(swarmkitAPI api.ControlClient, tlsConfig *tls.Config, enableCors bool) *mux.Router {
	r := mux.NewRouter()
	context := &context{
		swarmkitAPI: swarmkitAPI,
		tlsConfig:   tlsConfig,
		render:      render.New(),
	}

	setupPrimaryRouter(r, context, enableCors)
	return r
}

func setupPrimaryRouter(r *mux.Router, context *context, enableCors bool) {
	for method, mappings := range routes {
		for route, fct := range mappings {
			log.WithFields(log.Fields{"method": method, "route": route}).Debug("Registering HTTP route")

			localRoute := route
			localFct := fct

			wrap := func(w http.ResponseWriter, r *http.Request) {
				log.WithFields(log.Fields{"method": r.Method, "uri": r.RequestURI}).Debug("HTTP request received")
				if enableCors {
					writeCorsHeaders(w, r)
				}
				localFct(context, w, r)
			}

			localMethod := method
			r.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)

			if enableCors {
				optionsMethod := "OPTIONS"
				localFct = optionsHandler

				wrap := func(w http.ResponseWriter, r *http.Request) {
					log.WithFields(log.Fields{"method": optionsMethod, "uri": r.RequestURI}).
						Debug("HTTP request received")
					if enableCors {
						writeCorsHeaders(w, r)
					}
					localFct(context, w, r)
				}

				r.Path(localRoute).Methods(optionsMethod).HandlerFunc(wrap)
			}
		}
	}
}
