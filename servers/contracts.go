package servers

import (
	"html/template"
	"net"
)

// IServer -
type IServer interface {
	Run(ipPort string) error
	RunOnExistingListener(listener net.Listener)
}

// JSONData -
type JSONData map[string]interface{}

// ServerAppContext -
type ServerAppContext struct {
	Templates *template.Template
	// TemplateSet *isokit.TemplateSet
}
