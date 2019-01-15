package servers

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// MixedServer -
type MixedServer struct {
	servers []IServer
}

// NewMixedServer -
func NewMixedServer(servers []IServer) IServer {
	return &MixedServer{
		servers: servers,
	}
}

// Run -
func (thisRef *MixedServer) Run(ipPort string, enableCORS bool) error {
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	var router = mux.NewRouter()
	thisRef.PrepareRoutes(router)
	thisRef.RunOnExistingListenerAndRouter(listener, router, enableCORS)

	return nil
}

// PrepareRoutes -
func (thisRef *MixedServer) PrepareRoutes(router *mux.Router) {
	for _, server := range thisRef.servers {
		server.PrepareRoutes(router)
	}
}

// RunOnExistingListenerAndRouter -
func (thisRef *MixedServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	if enableCORS {
		corsSetterHandler := cors.Default().Handler(router)
		log.Fatal(http.Serve(listener, corsSetterHandler))
	} else {
		log.Fatal(http.Serve(listener, router))
	}

	// var wg sync.WaitGroup

	// for _, server := range thisRef.servers {
	// 	wg.Add(1)

	// 	go func(s IServer) {
	// 		s.RunOnExistingListenerAndRouter(listener, router, enableCORS)
	// 		wg.Done()
	// 	}(server)
	// }

	// wg.Wait()
}
