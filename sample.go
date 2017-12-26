package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nic0lae/JerryMouse/Servers"
)

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// Define Handlers
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
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

type IncomingJson struct {
	Field1 string
	Field2 int
	Field3 float64
}

func jsonRequestHandler(data interface{}) Servers.JsonResponse {
	var response Servers.JsonResponse

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

func streamTelemetryRequestHandler(inChannel chan []byte, outChannel chan []byte) {
	// DOX:
	// `close(outChannel)` will close the connection
	// recieve err on `inChannel` means connection was closed

	for {
		data, ok := <-inChannel
		if !ok {
			// connection was closed
			return
		} else {
			log.Print("data:", data)
		}
	}
}

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// Run Server
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
func main() {
	apiServer := Servers.Api()

	apiServer.SetLowLevelHandlers([]Servers.LowLevelHandler{
		Servers.LowLevelHandler{
			Route:   "/SayHello",
			Handler: sayHelloRequestHandler,
			Verb:    "GET",
		},
		Servers.LowLevelHandler{
			Route:   "/EchoBack",
			Handler: echoBackRequestHandler,
			Verb:    "POST",
		},
	})

	apiServer.SetJsonHandlers([]Servers.JsonHandler{
		Servers.JsonHandler{
			Route:      "/",
			Handler:    jsonRequestHandler,
			JsonObject: &IncomingJson{},
		},
	})

	apiServer.SetRealtimeHandlers([]Servers.RealtimeHandler{
		Servers.RealtimeHandler{
			Route:   "/StreamTelemetry",
			Handler: streamTelemetryRequestHandler,
		},
	})

	apiServer.Run(":9999")
}
