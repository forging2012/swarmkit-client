package api

import (
	"crypto/tls"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Dispatcher is a meta http.Handler. It acts as an http.Handler and forwards
// requests to another http.Handler that can be changed at runtime.
type dispatcher struct {
	handler http.Handler
}

// SetHandler changes the underlying handler.
func (d *dispatcher) SetHandler(handler http.Handler) {
	d.handler = handler
}

// ServeHTTP forwards requests to the underlying handler.
func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if d.handler == nil {
		httpError(w, "No dispatcher defined", http.StatusInternalServerError)
	}
	d.handler.ServeHTTP(w, r)
}

// Server is a swarmkit API server.
type Server struct {
	tlsConfig  *tls.Config
	dispatcher *dispatcher
	host       string
}

// NewServer creates an api.Server.
func NewServer(host string, tlsConfig *tls.Config) *Server {
	return &Server{
		tlsConfig:  tlsConfig,
		dispatcher: &dispatcher{},
		host:       host,
	}
}

// SetHandler is used to overwrite the HTTP handler for the API.
// This can be the api router or a reverse proxy.
func (s *Server) SetHandler(handler http.Handler) {
	s.dispatcher.SetHandler(handler)
}

func newListener(proto, addr string, tlsConfig *tls.Config) (net.Listener, error) {
	l, err := net.Listen(proto, addr)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		tlsConfig.NextProtos = []string{"http/1.1"}
		l = tls.NewListener(l, tlsConfig)
	}
	return l, nil
}

// ListenAndServe provide http rest api
func (s *Server) ListenAndServe() (err error) {

	var (
		l      net.Listener
		server = &http.Server{
			Addr:    s.host,
			Handler: s.dispatcher,
		}
	)
	if l, err = newListener("tcp", s.host, s.tlsConfig); err != nil {
		return
	}

	log.Infof("Start API server at: %s", s.host)
	return server.Serve(l)
}
