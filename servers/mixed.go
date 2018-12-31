package servers

import (
	"net"
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

	for _, server := range thisRef.servers {
		server.RunOnExistingListener(listener)
	}

	return nil
}

// RunOnExistingListener -
func (thisRef *MixedServer) RunOnExistingListener(listener net.Listener) {
	for _, server := range thisRef.servers {
		server.RunOnExistingListener(listener)
	}
}
