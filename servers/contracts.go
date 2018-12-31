package servers

import (
	"html/template"
	"net/http"
)

// Server -
type Server interface {
	Run(ipPort string)
}

// EmptyObject -
type EmptyObject struct{}

// LowLevelRequestHandler -
type LowLevelRequestHandler func(rw http.ResponseWriter, r *http.Request)

// RealtimeRequestHandler -
type RealtimeRequestHandler func(inChannel chan []byte, outChannel chan []byte)

// LowLevelHandler -
type LowLevelHandler struct {
	Route   string
	Handler LowLevelRequestHandler
	Verb    string
}

// JSONData -
type JSONData map[string]interface{}

// RealtimeHandler -
type RealtimeHandler struct {
	Route   string
	Handler RealtimeRequestHandler
}

// RealtimeClient -
type RealtimeClient struct {
	Address string
	Peers   []string
}

// ServerAppContext -
type ServerAppContext struct {
	Templates *template.Template
	// TemplateSet *isokit.TemplateSet
}
