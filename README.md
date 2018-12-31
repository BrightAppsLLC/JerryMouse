
### ![](https://raw.github.com/brightappsllc/JerryMouse/master/logo.png)

### Build API Servers
```go
package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brightappsllc/JerryMouse/Servers"
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
	Field4 string
}

func jsonRequestHandler(data interface{}) Servers.JSONResponse {
	dataAsJson, ok := data.(*IncomingJson)
	if !ok {
		return Servers.JSONResponse{HasError: true, ErrorMessage: "Invalid Params"}
	}

	// Input params seem ok, Process & Set Fields
	var response Servers.JSONResponse
	response.Data = dataAsJson

	return response
}

func streamTelemetryRequestHandler(inChannel chan []byte, outChannel chan []byte) { //, done chan bool) {
	// DOX:
	// `close(outChannel)` will close the connection
	// if error when reading on `inChannel` means connection was closed, do not send data

	go func() {
		for {
			data, ok := <-inChannel
			if !ok {
				close(outChannel)
				break
			} else {
				println("RECV: " + string(data))
			}
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)

			dataToSend := "Async Hi From Server @ " + time.Now().Format(time.RFC3339)

			select {
			case outChannel <- []byte(dataToSend):
				// message sent - all looks ok
				println("SEND: " + dataToSend)
			default:
				// message not sent - connection was closed")
				close(outChannel)
				break
			}
		}
	}()
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
		Servers.JsonHandler{
			Route:      "/MyRestEndopint",
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
```

##### Test With
```shell
curl localhost:9999 -H "Content-Type: application/json" -X POST -d '{"Field1": "value1", "Fieldd2": 1, "Field3": 0.1, "Field4": "another string"}'
```

```shell
curl localhost:9999/SayHello -X GET
```

```shell
curl localhost:9999/EchoBack  -X POST -d 'Helloooooooooooo'
```

For the WebSockets example, open the `sample.html` page

