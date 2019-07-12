package servers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

// JSONResponse -
type JSONResponse struct {
	HasError     bool        `json:"hasError"`
	ErrorMessage string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// JSONRequestHandler -
//type JSONRequestHandler func(data interface{}) JsonResponse
type JSONRequestHandler func(data []byte) JSONResponse

// JSONHandler -
type JSONHandler struct {
	Route      string
	Handler    JSONRequestHandler
	JSONObject interface{}
}

// JSONServer -
type JSONServer struct {
	handlers       []JSONHandler
	routeToHandler map[string]JSONHandler
	lowLevelServer IServer
}

// NewJSONServer -
func NewJSONServer(handlers []JSONHandler) *JSONServer {

	var thisRef = &JSONServer{
		handlers:       handlers,
		routeToHandler: map[string]JSONHandler{},
		lowLevelServer: nil,
	}

	var lowLevelRequestHelper = func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var jsonHandler = thisRef.routeToHandler[r.URL.Path]

		// Pass Object
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(rw, "can't read body", http.StatusBadRequest)
			return
		}

		var jsonResponse = jsonHandler.Handler(body)
		json.NewEncoder(rw).Encode(jsonResponse)
	}

	var lowLevelHandlers = []LowLevelHandler{}

	for _, handler := range thisRef.handlers {
		thisRef.routeToHandler[handler.Route] = handler

		lowLevelHandlers = append(lowLevelHandlers, LowLevelHandler{
			Route:   handler.Route,
			Handler: lowLevelRequestHelper,
			Verb:    "POST",
		})
	}

	thisRef.lowLevelServer = NewLowLevelServer(lowLevelHandlers)

	return thisRef
}

// Run - Implement `IServer`
func (thisRef *JSONServer) Run(ipPort string, enableCORS bool) error {
	return thisRef.lowLevelServer.Run(ipPort, enableCORS)
}

// PrepareRoutes - Implement `IServer`
func (thisRef *JSONServer) PrepareRoutes(router *mux.Router) {
	thisRef.lowLevelServer.PrepareRoutes(router)
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *JSONServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	thisRef.lowLevelServer.RunOnExistingListenerAndRouter(listener, router, enableCORS)
}

// func (jsonData *JsonData) ToObject(objectInstance interface{}) {
// 	// Do JSON to Object Mapping
// 	objectValue := reflect.ValueOf(objectInstance).Elem()
// 	for i := 0; i < objectValue.NumField(); i++ {
// 		field := objectValue.Field(i)
// 		fieldName := objectValue.Type().Field(i).Name

// 		if valueToCopy, ok := (*jsonData)[fieldName]; ok {
// 			if !field.CanInterface() {
// 				continue
// 			}
// 			switch field.Interface().(type) {
// 			case string:
// 				valueToCopyAsString := reflect.ValueOf(valueToCopy).String()
// 				objectValue.Field(i).SetString(valueToCopyAsString)
// 				break
// 			case int:
// 				valueToCopyAsInt := int64(reflect.ValueOf(valueToCopy).Float())
// 				objectValue.Field(i).SetInt(valueToCopyAsInt)
// 				break
// 			case float64:
// 				valueToCopyAsFloat := reflect.ValueOf(valueToCopy).Float()
// 				objectValue.Field(i).SetFloat(valueToCopyAsFloat)
// 				break
// 			default:
// 			}
// 		}
// 	}
// }

// Get JSON fields
//var jsonData JsonData
//_ = json.NewDecoder(r.Body).Decode(&jsonData)

// TRACE
// if false {
// 	reqAsJSON, _ := json.Marshal(req)
// 	fmt.Println(fmt.Sprintf("%s -> %s", Utils.CallStack(), string(reqAsJSON)))
// }

//jsonData.ToObject(jsonHandler.JsonObject)

// Pass Object
//var response JsonResponse = jsonHandler.Handler(jsonData)
