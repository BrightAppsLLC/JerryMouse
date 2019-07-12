package servers

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// RealtimeRequestHandler -
type RealtimeRequestHandler func(inChannel chan []byte, outChannel chan []byte)

// RealtimeHandler -
type RealtimeHandler struct {
	Route   string
	Handler RealtimeRequestHandler
}

// RealtimeServer -
type RealtimeServer struct {
	handlers       []RealtimeHandler
	routeToHandler map[string]RealtimeHandler
	lowLevelServer IServer
	peers          []*websocket.Conn
	peersSync      sync.RWMutex
}

// NewRealtimeServer -
func NewRealtimeServer(handlers []RealtimeHandler) IServer {

	var thisRef = &RealtimeServer{
		handlers:       handlers,
		routeToHandler: map[string]RealtimeHandler{},
		lowLevelServer: nil,
		peers:          []*websocket.Conn{},
		peersSync:      sync.RWMutex{},
	}

	var lowLevelRequestHelper = func(rw http.ResponseWriter, r *http.Request) {
		r.Header["Origin"] = nil

		var handler RealtimeHandler = thisRef.routeToHandler[r.URL.Path]

		var upgrader = websocket.Upgrader{}
		ws, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Print("upgrade: ", err)
			return
		}

		thisRef.setupCommunication(ws, &handler)
	}

	var lowLevelHandlers = []LowLevelHandler{}

	for _, handler := range thisRef.handlers {
		thisRef.routeToHandler[handler.Route] = handler

		lowLevelHandlers = append(lowLevelHandlers, LowLevelHandler{
			Route:   handler.Route,
			Handler: lowLevelRequestHelper,
			Verb:    "GET",
		})
	}

	thisRef.lowLevelServer = NewLowLevelServer(lowLevelHandlers)

	return thisRef
}

// Run - Implement `IServer`
func (thisRef *RealtimeServer) Run(ipPort string, enableCORS bool) error {
	return thisRef.lowLevelServer.Run(ipPort, enableCORS)
}

// PrepareRoutes - Implement `IServer`
func (thisRef *RealtimeServer) PrepareRoutes(router *mux.Router) {
	thisRef.lowLevelServer.PrepareRoutes(router)
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *RealtimeServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	thisRef.lowLevelServer.RunOnExistingListenerAndRouter(listener, router, enableCORS)
}

func (thisRef *RealtimeServer) setupCommunication(ws *websocket.Conn, handler *RealtimeHandler) {
	thisRef.addPeer(ws)

	var inChannel = make(chan []byte)
	var outChannel = make(chan []byte)

	var once sync.Once
	closeInChannel := func() {
		close(inChannel)
	}

	var wg sync.WaitGroup

	// outChannel -> PEER
	wg.Add(1)
	go func() {
		fmt.Println("SEND-TO-PEER - START")

		for {
			data, readOk := <-outChannel
			if !readOk {
				break
			}

			err := ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				break
			}
		}

		fmt.Println("SEND-TO-PEER - END")
		once.Do(closeInChannel)
		wg.Done()
	}()

	// PEER -> inChannel
	wg.Add(1)
	go func() {
		fmt.Println("READ-FROM-PEER - START")

		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				break
			}

			var haveToStop = false
			select {
			case inChannel <- []byte(data):
			default:
				haveToStop = true
				break
			}

			if haveToStop {
				break
			}
		}

		fmt.Println("READ-FROM-PEER - END")
		once.Do(closeInChannel)
		wg.Done()
	}()

	go handler.Handler(inChannel, outChannel)

	wg.Wait()
	fmt.Println("setupCommunication - DONE")
	thisRef.removePeer(ws)
}

// SendToAllPeers -
func (thisRef *RealtimeServer) SendToAllPeers(data []byte) {
	thisRef.peersSync.RLock()
	defer thisRef.peersSync.RUnlock()

	for _, conn := range thisRef.peers {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (thisRef *RealtimeServer) addPeer(peer *websocket.Conn) {
	thisRef.peersSync.Lock()
	defer thisRef.peersSync.Unlock()

	thisRef.peers = append(thisRef.peers, peer)
}

func (thisRef *RealtimeServer) removePeer(peer *websocket.Conn) {
	thisRef.peersSync.Lock()
	defer thisRef.peersSync.Unlock()

	index := -1
	for i, conn := range thisRef.peers {
		if conn == peer {
			index = i
			break
		}
	}
	if index != -1 {
		thisRef.peers = append(thisRef.peers[:index], thisRef.peers[index+1:]...)
	}

	peer.Close()
}
