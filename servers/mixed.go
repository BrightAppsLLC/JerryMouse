package servers

import (
	"net"
	"sync"

	"github.com/gorilla/mux"
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
func (thisRef *MixedServer) Run(ipPort string) error {
	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		return err
	}

	thisRef.RunOnExistingListenerAndRouter(listener, mux.NewRouter())

	return nil
}

// RunOnExistingListenerAndRouter -
func (thisRef *MixedServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router) {
	var wg sync.WaitGroup

	for _, server := range thisRef.servers {
		wg.Add(1)

		go func(s IServer) {
			s.RunOnExistingListenerAndRouter(listener, router)
			wg.Done()
		}(server)
	}

	wg.Wait()
}
