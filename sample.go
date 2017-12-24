package main

import (
	"github.com/nic0lae/JerryMouse/ApiServers"
)

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// Define Handler
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
type IncomingJson struct {
	Field1 string
	Field2 string
	Field3 string
}

func jsonRequestHandler(data interface{}) ApiServers.JsonResponse {
	var response ApiServers.JsonResponse

	dataAsJson, ok := data.(IncomingJson)

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
	apiServer.Run(":9999")
}
