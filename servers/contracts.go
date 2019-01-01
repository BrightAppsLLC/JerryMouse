package servers

import (
	"net"

	"github.com/gorilla/mux"
)

// IServer -
type IServer interface {
	Run(ipPort string) error
	RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router)
}
