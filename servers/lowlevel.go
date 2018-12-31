package servers

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// LowLevelRequestHandler -
type LowLevelRequestHandler func(rw http.ResponseWriter, r *http.Request)

// LowLevelHandler -
type LowLevelHandler struct {
	Route   string
	Handler LowLevelRequestHandler
	Verb    string
}

// LowLevelServer -
type LowLevelServer struct {
	handlers       []LowLevelHandler
	enableCORS     bool
	routeToHandler map[string]RealtimeHandler
}

// NewLowLevelServer -
func NewLowLevelServer(enableCORS bool, handlers []LowLevelHandler) IServer {
	return &LowLevelServer{
		handlers:       handlers,
		enableCORS:     enableCORS,
		routeToHandler: map[string]RealtimeHandler{},
	}
}

// Run -
func (thisRef *LowLevelServer) Run(ipPort string) error {
	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		return err
	}

	thisRef.RunOnExistingListener(listener)

	return nil
}

// RunOnExistingListener -
func (thisRef *LowLevelServer) RunOnExistingListener(listener net.Listener) {
	router := mux.NewRouter()

	for _, handler := range thisRef.handlers {
		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb)
	}

	if thisRef.enableCORS {
		corsSetterHandler := cors.Default().Handler(router)
		log.Fatal(http.Serve(listener, corsSetterHandler))
	} else {
		log.Fatal(http.Serve(listener, router))
	}
}
