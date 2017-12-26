
### ![](https://raw.github.com/nic0lae/JerryMouse/master/logo.png)

### Build API Server
```go
import (
	"io/ioutil"
	"net/http"

	"github.com/nic0lae/JerryMouse/Servers"
)

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// Define Handlers
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
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
	apiServer := Servers.Api()
	apiServer.SetJsonHandlers([]Servers.JsonHandler{
		Servers.JsonHandler{
			Route:      "/",
			Handler:    jsonRequestHandler,
			JsonObject: &IncomingJson{},
		},
	})
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
	apiServer.Run(":9999")
}
```

##### Test with
```shell
curl localhost:9999 -H "Content-Type: application/json" -X POST -d '{"Field1": "value1", "Fieldd2": 1, "Field3": 0.1}'
```

```shell
curl localhost:9999/SayHello -X GET
```

```shell
curl localhost:9999/EchoBack  -X POST -d 'Helloooooooooooo'
```
