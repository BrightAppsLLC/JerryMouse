package ApiServers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// HttpApi - Creates a new HttpApiServer
func HttpApi() *HttpApiServer {
	return &HttpApiServer{
		lowLevelHandlers: []LowLevelHandler{},
		jsonHandlers:     []JsonHandler{},
	}
}

// SetLowLevelHandlers - Sets HttpApiServer handlers
func (has *HttpApiServer) SetLowLevelHandlers(handlers []LowLevelHandler) {
	has.lowLevelHandlers = handlers
}

// SetJsonHandlers - Sets HttpApiServer handlers
func (has *HttpApiServer) SetJsonHandlers(handlers []JsonHandler) {
	has.jsonHandlers = handlers
}

var routeToDelegate map[string]JsonHandler

// Run - Runs HttpApiServer
func (has HttpApiServer) Run(ipPort string) {
	router := mux.NewRouter()

	for _, handler := range has.lowLevelHandlers {
		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb)
	}

	routeToDelegate = map[string]JsonHandler{}

	for _, handler := range has.jsonHandlers {
		routeToDelegate[handler.Route] = handler
		router.HandleFunc(handler.Route, helpers_LowLevelRequestDelegate).Methods("POST")
	}

	log.Fatal(http.ListenAndServe(ipPort, router))
}

func helpers_LowLevelRequestDelegate(rw http.ResponseWriter, r *http.Request) {
	var handler JsonHandler = routeToDelegate[r.URL.Path]

	// Get JSON fields
	var jsonData JsonData
	_ = json.NewDecoder(r.Body).Decode(&jsonData)

	// Do JSON to Object Mapping
	objectValue := reflect.ValueOf(handler.JsonObject).Elem()
	for i := 0; i < objectValue.NumField(); i++ {
		fieldName := objectValue.Type().Field(i).Name

		if valueToCopy, ok := jsonData[fieldName]; ok {
			valueToCopyAsString := reflect.ValueOf(valueToCopy).String()
			objectValue.Field(i).SetString(valueToCopyAsString)
		}
	}

	// Pass Object
	var response JsonResponse = handler.Handler(handler.JsonObject)
	json.NewEncoder(rw).Encode(response)
}

// func helpers_LowLevelRequestDelegate2(rw http.ResponseWriter, r *http.Request) {
// 	var handler JsonHandler = routeToDelegate[r.URL.Path]

// 	data, _ := ioutil.ReadAll(r.Body)
// 	dataAsString := string(data)
// 	println(dataAsString)

// 	var dataAsJson JsonData
// 	json.Unmarshal(data, &dataAsJson)

// 	var response JsonResponse = handler.Handler(dataAsJson)

// 	response.Data["error"] = response.Error
// 	responseAsByteArray, _ := json.Marshal(response.Data)

// 	rw.Write(responseAsByteArray)
// 	// json.NewEncoder(rw).Encode(responseAsByteArray)
// }
