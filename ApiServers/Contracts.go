package ApiServers

import (
	"net/http"
)

type LowLevelRequestHandler func(rw http.ResponseWriter, r *http.Request)
type JsonRequestHandler func(data interface{}) JsonResponse

type LowLevelHandler struct {
	Route   string
	Handler LowLevelRequestHandler
	Verb    string
}

type JsonHandler struct {
	Route      string
	Handler    JsonRequestHandler
	JsonObject interface{}
}

type JsonData map[string]interface{}
type JsonResponse struct {
	Error string `json:"Error"`
	Data  interface{}
}

type HttpApiServer struct {
	lowLevelHandlers []LowLevelHandler
	jsonHandlers     []JsonHandler
}
