package servers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// LowLevelServer -
type LowLevelServer struct {
	handlers       []LowLevelHandler
	enableCors     bool
	routeToHandler map[string]RealtimeHandler
}

// NewLowLevelServer -
func NewLowLevelServer() *LowLevelServer {
	return &LowLevelServer{
		handlers:       []LowLevelHandler{},
		enableCors:     false,
		routeToHandler: map[string]RealtimeHandler{},
	}
}

// SetHandlers -
func (thisRef *LowLevelServer) SetHandlers(handlers []LowLevelHandler) {
	thisRef.handlers = handlers
}

// EnableCORS -
func (thisRef *LowLevelServer) EnableCORS() {
	thisRef.enableCors = true
}

// Run - Runs LowLevelServer
func (thisRef LowLevelServer) Run(ipPort string) {
	router := mux.NewRouter()

	for _, handler := range thisRef.handlers {
		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb)
	}

	if thisRef.enableCors {
		corsSetterHandler := cors.Default().Handler(router)
		log.Fatal(http.ListenAndServe(ipPort, corsSetterHandler))
	} else {
		log.Fatal(http.ListenAndServe(ipPort, router))
	}
}
