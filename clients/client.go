package clients

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

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
