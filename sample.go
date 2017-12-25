package main

import (
	"io/ioutil"
	"net/http"

	"github.com/nic0lae/JerryMouse/ApiServers"
)

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// Define Handlers
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
type IncomingJson struct {
	Field1 string
	Field2 int
	Field3 float64
}

func jsonRequestHandler(data interface{}) ApiServers.JsonResponse {
	var response ApiServers.JsonResponse

	dataAsJson, ok := data.(*IncomingJson)

	if !ok {
		response.Error = "Invalid Params"
		response.Data = dataAsJson
	} else {
		// Process & Set Fields
		response.Error = ""
		response.Data = dataAsJson
	}

	return response
}

func sayHelloRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello !"))
}

func echoBackRequestHandler(rw http.ResponseWriter, r *http.Request) {
	data, ok := ioutil.ReadAll(r.Body)
	if ok != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else {
		rw.Write(data)
	}
}

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// Run Server
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
func main() {
	apiServer := ApiServers.HttpApi()
	apiServer.SetJsonHandlers([]ApiServers.JsonHandler{
		ApiServers.JsonHandler{
			Route:      "/",
			Handler:    jsonRequestHandler,
			JsonObject: &IncomingJson{},
		},
	})
	apiServer.SetLowLevelHandlers([]ApiServers.LowLevelHandler{
		ApiServers.LowLevelHandler{
			Route:   "/SayHello",
			Handler: sayHelloRequestHandler,
			Verb:    "GET",
		},
		ApiServers.LowLevelHandler{
			Route:   "/EchoBack",
			Handler: echoBackRequestHandler,
			Verb:    "POST",
		},
	})
	apiServer.Run(":9999")
}
