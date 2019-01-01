package servers

import (
	"html/template"
	"net"

	"github.com/gorilla/mux"
)

// IServer -
type IServer interface {
	Run(ipPort string) error
	RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router)
}

// JSONData -
type JSONData map[string]interface{}

// ServerAppContext -
type ServerAppContext struct {
	Templates *template.Template
	// TemplateSet *isokit.TemplateSet
}
