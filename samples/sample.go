package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	jmS "github.com/brightappsllc/JerryMouse/servers"
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
	Field4 bool
}

func jsonRequestHandler(data []byte) jmS.JSONResponse {
	var incomingJson = &IncomingJson{}

	err := json.Unmarshal(data, incomingJson)
	if err != nil {
		return jmS.JSONResponse{
			HasError: true,
			Message:  err.Error(),
		}
	}

	// Input params seem ok, Process & Set Fields
	var response jmS.JSONResponse
	response.Data = incomingJson

	return response
}

func streamTelemetryRequestHandler(inChannel chan []byte, outChannel chan []byte) { //, done chan bool) {
	// DOX:
	//    To indicate DONE close the `outChannel`
	//    If error when reading on `inChannel` means connection was closed, do not send data

	var wg sync.WaitGroup
	var rw sync.RWMutex
	var outChannelClosed = false // writing to a closed channel will panic

	wg.Add(1)
	go func() {
		for {
			data, readOk := <-inChannel
			if !readOk {
				break
			} else {
				println("RECV: " + string(data))
			}
		}

		rw.Lock()
		outChannelClosed = true
		close(outChannel)
		rw.Unlock()

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			time.Sleep(1 * time.Second)

			dataToSend := "Async Hi From Server @ " + time.Now().Format(time.RFC3339)

			var haveToStop = false
			rw.Lock()
			if !outChannelClosed {
				select {
				case outChannel <- []byte(dataToSend):
					println("SEND: " + dataToSend)
				default:
					haveToStop = true
					break
				}
			} else {
				haveToStop = true
			}

			if haveToStop {
				break
			}
			rw.Unlock()
		}

		wg.Done()
	}()

	wg.Wait()
}

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// Run Server
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
func main() {
	const port = 9999

	fmt.Println(fmt.Sprintf("curl localhost:%d/SayHello", port))
	fmt.Println(fmt.Sprintf("curl localhost:%d/EchoBack -X POST -d 'should-send-back-the-same'", port))
	fmt.Println(fmt.Sprintf("curl localhost:%d/MyRestEndopint -H \"Content-Type: application/json\" -X POST -d '{\"Field1\":\"val1\", \"Field2\": 0, \"Field3\": 1.0, \"Field4\": true}'", port))
	fmt.Println("ALSO: open the 'sample.html' to see data streaming")

	jmS.NewMixedServer([]jmS.IServer{
		jmS.NewHTTPServer([]jmS.HTTPHandler{
			{
				Route:   "/SayHello",
				Verb:    "GET",
				Handler: sayHelloRequestHandler,
			},
			{
				Route:   "/EchoBack",
				Verb:    "POST",
				Handler: echoBackRequestHandler,
			},
		}),
		jmS.NewJSONServer([]jmS.JSONHandler{
			{
				Route:    "/MyRestEndopint",
				Template: &IncomingJson{},
				Handler:  jsonRequestHandler,
			},
		}),
		jmS.NewWebScoketsServer([]jmS.WebScoketsHandler{
			jmS.WebScoketsHandler{
				Route:   "/StreamTelemetry",
				Handler: streamTelemetryRequestHandler,
			},
		}),
	}).Run(fmt.Sprintf(":%d", port), false)
}
