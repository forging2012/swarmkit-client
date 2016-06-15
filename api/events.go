package api

import (
	"io"
	"net/http"
	"sync"
	"time"
)

// EventsHandler broadcasts events to multiple client listeners.
type eventsHandler struct {
	sync.RWMutex
	ws map[string]io.Writer
	cs map[string]chan struct{}
}

// NewEventsHandler creates a new EventsHandler for a cluster.
// The new eventsHandler is initialized with no writers or channels.
func newEventsHandler() *eventsHandler {
	return &eventsHandler{
		ws: make(map[string]io.Writer),
		cs: make(map[string]chan struct{}),
	}
}

// Add adds the writer and a new channel for the remote address.
func (eh *eventsHandler) Add(remoteAddr string, w io.Writer) {
	eh.Lock()
	eh.ws[remoteAddr] = w
	eh.cs[remoteAddr] = make(chan struct{})
	eh.Unlock()
}

// Wait waits on a signal from the remote address.
func (eh *eventsHandler) Wait(remoteAddr string, until int64) {

	timer := time.NewTimer(0)
	timer.Stop()
	if until > 0 {
		dur := time.Unix(until, 0).Sub(time.Now())
		timer = time.NewTimer(dur)
	}

	// subscribe to http client close event
	w := eh.ws[remoteAddr]
	var closeNotify <-chan bool
	if closeNotifier, ok := w.(http.CloseNotifier); ok {
		closeNotify = closeNotifier.CloseNotify()
	}

	select {
	case <-eh.cs[remoteAddr]:
	case <-closeNotify:
	case <-timer.C: // `--until` timeout
		close(eh.cs[remoteAddr])
	}
	eh.cleanupHandler(remoteAddr)
}

func (eh *eventsHandler) cleanupHandler(remoteAddr string) {
	eh.Lock()
	// the maps are expected to have the same keys
	delete(eh.cs, remoteAddr)
	delete(eh.ws, remoteAddr)
	eh.Unlock()

}
