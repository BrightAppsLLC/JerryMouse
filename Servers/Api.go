package Servers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// Api - Creates a new ApiServer
func Api() *ApiServer {
	return &ApiServer{
		lowLevelHandlers: []LowLevelHandler{},
		jsonHandlers:     []JsonHandler{},
		shouldEnableCors: false,
	}
}

// SetLowLevelHandlers - Sets ApiServer handlers
func (has *ApiServer) SetLowLevelHandlers(handlers []LowLevelHandler) {
	has.lowLevelHandlers = handlers
}

// SetJsonHandlers - Sets ApiServer handlers
func (has *ApiServer) SetJsonHandlers(handlers []JsonHandler) {
	has.jsonHandlers = handlers
}

func (has *ApiServer) SetRealtimeHandlers(handlers []RealtimeHandler) {
	has.realtimeHandlers = handlers
}

func (has *ApiServer) SendToAllRealtimePeers(data []byte) {
	for _, conn := range realtimePeers {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (has *ApiServer) EnableCORS() {
	has.shouldEnableCors = true
}

var jsonRouteToHandler map[string]JsonHandler
var realtimeRouteToHandler map[string]RealtimeHandler
var realtimePeers []*websocket.Conn

// Run - Runs ApiServer
func (has ApiServer) Run(ipPort string) {
	router := mux.NewRouter()

	// Low level
	for _, handler := range has.lowLevelHandlers {
		router.HandleFunc(handler.Route, handler.Handler).Methods(handler.Verb)
	}

	// JSON
	jsonRouteToHandler = map[string]JsonHandler{}
	for _, handler := range has.jsonHandlers {
		jsonRouteToHandler[handler.Route] = handler
		router.HandleFunc(handler.Route, helpers_LowLevelRequestDelegate).Methods("POST")
	}

	// Realtime
	realtimeRouteToHandler = map[string]RealtimeHandler{}
	realtimePeers = []*websocket.Conn{}
	for _, handler := range has.realtimeHandlers {
		realtimeRouteToHandler[handler.Route] = handler
		router.HandleFunc(handler.Route, helpers_RealtimeRequestDelegate)
	}

	// TRACE
	// fmt.Println(fmt.Sprintf("%s -> Ready.", Utils.CallStack()))

	if has.shouldEnableCors {
		corsSetterHandler := cors.Default().Handler(router)
		log.Fatal(http.ListenAndServe(ipPort, corsSetterHandler))
	} else {
		log.Fatal(http.ListenAndServe(ipPort, router))
	}
}

func helpers_LowLevelRequestDelegate(rw http.ResponseWriter, r *http.Request) {
	var handler JsonHandler = jsonRouteToHandler[r.URL.Path]

	// Get JSON fields
	//var jsonData JsonData
	//_ = json.NewDecoder(r.Body).Decode(&jsonData)

	// TRACE
	// if false {
	// 	reqAsJSON, _ := json.Marshal(req)
	// 	fmt.Println(fmt.Sprintf("%s -> %s", Utils.CallStack(), string(reqAsJSON)))
	// }

	//jsonData.ToObject(handler.JsonObject)

	// Pass Object
	//var response JsonResponse = handler.Handler(handler.JsonObject)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(rw, "can't read body", http.StatusBadRequest)
		return
	}

	var response JsonResponse = handler.Handler(body)
	//var response JsonResponse = handler.Handler(jsonData)
	json.NewEncoder(rw).Encode(response)
}

func helpers_RemoveRealtimePeer(peer *websocket.Conn) {
	index := -1
	for i, conn := range realtimePeers {
		if conn == peer {
			index = i
			break
		}
	}
	if index != -1 {
		realtimePeers = append(realtimePeers[:index], realtimePeers[index+1:]...)
	}
}

func helpers_RealtimeRequestDelegate(rw http.ResponseWriter, r *http.Request) {
	r.Header["Origin"] = nil

	var handler RealtimeHandler = realtimeRouteToHandler[r.URL.Path]

	var upgrader = websocket.Upgrader{}
	ws, err := upgrader.Upgrade(rw, r, nil) // Upgrade the connection to a websocket
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}

	realtimePeers = append(realtimePeers, ws)

	helpers_RealtimeCommunicationHandler(ws, &handler)
}

func helpers_RealtimeCommunicationHandler(ws *websocket.Conn, handler *RealtimeHandler) {
	var inputChannel = make(chan []byte)
	var outputChannel = make(chan []byte)

	var once sync.Once
	closeInputChannel := func() {
		close(inputChannel)
	}

	// GO channel DOX: senders close; receivers check for closed.

	go func() {
		for {
			data, ok := <-outputChannel
			if !ok {
				// `outputChannel` from hook closed, means we have to close connection
				ws.Close()
				helpers_RemoveRealtimePeer(ws)
				return
			}

			err := ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				// TRACE
				// if false {
				// log.Print("WebSocket-Write-Error:", err) // TODO: LOG
				// }

				// we can't write means connection is closed, meaans we have to close chanels
				once.Do(closeInputChannel)
				return
			}
		}
	}()

	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				// TRACE
				// if false {
				// log.Println("WebSocket-Read-Error: ", err) // TODO: LOG
				// }

				// we can't read means connection is closed, meaans we have to close chanels
				once.Do(closeInputChannel)
				return
			}

			inputChannel <- data
		}
	}()

	go handler.Handler(inputChannel, outputChannel) //, doneChannel)
}

// func helpers_LowLevelRequestDelegate2(rw http.ResponseWriter, r *http.Request) {
// 	var handler JsonHandler = jsonRouteToHandler[r.URL.Path]

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

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// WS Client
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~

func NewRealtimeClient() *RealtimeClient {
	return &RealtimeClient{
		Address: "",
		Peers:   []string{},
	}
}

func (rtc *RealtimeClient) ConnectToPeer(url url.URL, realtimeHandler *RealtimeHandler) {
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	helpers_RealtimeCommunicationHandler(c, realtimeHandler)
}

// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
// JSON Helpers
// ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~
func (jsonData *JsonData) ToObject(objectInstance interface{}) {
	// Do JSON to Object Mapping
	objectValue := reflect.ValueOf(objectInstance).Elem()
	for i := 0; i < objectValue.NumField(); i++ {
		field := objectValue.Field(i)
		fieldName := objectValue.Type().Field(i).Name

		if valueToCopy, ok := (*jsonData)[fieldName]; ok {
			if !field.CanInterface() {
				continue
			}
			switch field.Interface().(type) {
			case string:
				valueToCopyAsString := reflect.ValueOf(valueToCopy).String()
				objectValue.Field(i).SetString(valueToCopyAsString)
				break
			case int:
				valueToCopyAsInt := int64(reflect.ValueOf(valueToCopy).Float())
				objectValue.Field(i).SetInt(valueToCopyAsInt)
				break
			case float64:
				valueToCopyAsFloat := reflect.ValueOf(valueToCopy).Float()
				objectValue.Field(i).SetFloat(valueToCopyAsFloat)
				break
			default:
			}
		}
	}
}
