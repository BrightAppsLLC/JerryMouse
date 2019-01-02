package servers

import (
	"fmt"
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
	routeToHandler map[string]RealtimeHandler
}

// NewLowLevelServer -
func NewLowLevelServer(handlers []LowLevelHandler) IServer {
	return &LowLevelServer{
		handlers:       handlers,
		routeToHandler: map[string]RealtimeHandler{},
	}
}

// Run -
func (thisRef *LowLevelServer) Run(ipPort string, enableCORS bool) error {
	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		return err
	}

	var router = mux.NewRouter()
	thisRef.PrepareRoutes(router)
	thisRef.RunOnExistingListenerAndRouter(listener, router, enableCORS)

	return nil
}

// PrepareRoutes -
func (thisRef *LowLevelServer) PrepareRoutes(router *mux.Router) {
	for _, handler := range thisRef.handlers {
		fmt.Println(fmt.Sprintf("LLS: %s - for %s", handler.Route, handler.Verb))
		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb)
	}
}

// RunOnExistingListenerAndRouter -
func (thisRef *LowLevelServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	if enableCORS {
		corsSetterHandler := cors.Default().Handler(router)
		log.Fatal(http.Serve(listener, corsSetterHandler))
	} else {
		log.Fatal(http.Serve(listener, router))
	}
}
