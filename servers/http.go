package servers

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// HTTPRequestHandler -
type HTTPRequestHandler func(rw http.ResponseWriter, r *http.Request)

// HTTPHandler -
type HTTPHandler struct {
	Route   string
	Verb    string
	Handler HTTPRequestHandler
}

// HTTPServer -
type HTTPServer struct {
	handlers []HTTPHandler
}

// NewHTTPServer -
func NewHTTPServer(handlers []HTTPHandler) IServer {
	return &HTTPServer{
		handlers: handlers,
	}
}

// Run - Implement `IServer`
func (thisRef *HTTPServer) Run(ipPort string, enableCORS bool) error {
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	thisRef.PrepareRoutes(router)
	thisRef.RunOnExistingListenerAndRouter(listener, router, enableCORS)

	return nil
}

// PrepareRoutes - Implement `IServer`
func (thisRef *HTTPServer) PrepareRoutes(router *mux.Router) {
	for _, handler := range thisRef.handlers {
		log.Debug(fmt.Sprintf("%s - for %s", handler.Route, handler.Verb))
		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb).Name(handler.Route)
	}
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *HTTPServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	if enableCORS {
		corsSetterHandler := cors.Default().Handler(router)
		err := http.Serve(listener, corsSetterHandler)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err := http.Serve(listener, router)
		if err != nil {
			log.Fatal(err)
		}
	}
}
