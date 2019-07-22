package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/brightappsllc/baqe-crypto/ssh"
	"github.com/gorilla/mux"

	golog "github.com/brightappsllc/golog"
	gologC "github.com/brightappsllc/golog/contracts"

	reflectionHelpers "github.com/brightappsllc/gohelpers/reflection"
)

// SSHTunnelServer -
type SSHTunnelServer struct {
	sshServerConfig *ssh.ServerConfig
	server          IServer
	router          *mux.Router
}

// NewSSHTunnelServer -
func NewSSHTunnelServer(sshServerConfig *ssh.ServerConfig, server IServer) IServer {
	return &SSHTunnelServer{
		sshServerConfig: sshServerConfig,
		server:          server,
		router:          mux.NewRouter(),
	}
}

// Run - Implement `IServer`
func (thisRef *SSHTunnelServer) Run(ipPort string, enableCORS bool) error {

	//
	// BASED-ON: https://godoc.org/github.com/brightappsllc/baqe-crypto/ssh#example-NewServerConn
	//

	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	thisRef.PrepareRoutes(thisRef.router)
	thisRef.RunOnExistingListenerAndRouter(listener, thisRef.router, enableCORS)

	return nil
}

// PrepareRoutes - Implement `IServer`
func (thisRef *SSHTunnelServer) PrepareRoutes(router *mux.Router) {
	thisRef.server.PrepareRoutes(router)
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *SSHTunnelServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			golog.Instance().LogErrorWithFields(gologC.Fields{
				"method":  reflectionHelpers.GetThisFuncName(),
				"message": fmt.Sprintf("JM-SSH: failed to accept incoming connection: %s", err),
			})

			continue
		}

		go thisRef.runSSH(connection)
	}
}

type customResponseWriter struct {
	http.ResponseWriter
	sshChannel ssh.Channel
}

func (thisRef *customResponseWriter) Write(data []byte) (int, error) {

	// golog.Instance().LogDebugWithFields(gologC.Fields{
	// 	"method":  reflectionHelpers.GetThisFuncName(),
	// 	"message": string(data),
	// })

	return thisRef.sshChannel.Write(data)
}

func (thisRef *SSHTunnelServer) runSSH(connection net.Conn) {
	// Before use, a handshake must be performed on the incoming connection
	sshServerConnection, chans, reqs, err := ssh.NewServerConn(connection, thisRef.sshServerConfig)
	if err != nil {
		golog.Instance().LogErrorWithFields(gologC.Fields{
			"method":  reflectionHelpers.GetThisFuncName(),
			"message": fmt.Sprintf("JM-SSH: failed to handshake: %s", err),
		})

		return
	}

	golog.Instance().LogInfoWithFields(gologC.Fields{
		"method":  reflectionHelpers.GetThisFuncName(),
		"message": fmt.Sprintf("JM-SSH: Connection %s", sshServerConnection.RemoteAddr().String()),
	})

	// golog.Instance().LogDebugWithFields(gologC.Fields{
	// 	"method":  reflectionHelpers.GetThisFuncName(),
	// 	"message": fmt.Sprintf("logged in with key %s", sshServerConn.),
	// })

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level protocol intended.
		// In the case of a shell, the type is "session" and ServerShell may be used
		// to present a simple terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, _, err := newChannel.Accept() // requests
		if err != nil {
			golog.Instance().LogErrorWithFields(gologC.Fields{
				"method":  reflectionHelpers.GetThisFuncName(),
				"message": fmt.Sprintf("JM-SSH: could not accept channel: %v", err),
			})
			break
		}

		go func() {
			defer channel.Close()
			for {
				data := make([]byte, 1000000)
				len, err := channel.Read(data)
				if err != nil {
					if strings.Compare(err.Error(), "EOF") == 0 {
						golog.Instance().LogInfoWithFields(gologC.Fields{
							"method":  reflectionHelpers.GetThisFuncName(),
							"message": fmt.Sprintf("JM-SSH: CONNECTION-CLOSED: %v", err),
						})
						break
					} else {
						golog.Instance().LogErrorWithFields(gologC.Fields{
							"method":  reflectionHelpers.GetThisFuncName(),
							"message": fmt.Sprintf("JM-SSH: DATA-ERROR: %v", err),
						})
						break
					}
				}

				data = data[0:len]
				golog.Instance().LogDebugWithFields(gologC.Fields{
					"method":  reflectionHelpers.GetThisFuncName(),
					"message": fmt.Sprintf("JM-SSH: SSH-DATA: %s", string(data)),
				})

				apiEndpoing := APIEndpoint{}
				err = json.Unmarshal(data, &apiEndpoing)
				if err != nil {
					golog.Instance().LogErrorWithFields(gologC.Fields{
						"method":  reflectionHelpers.GetThisFuncName(),
						"message": fmt.Sprintf("JM-SSH: Missing ROUTE: %s", err.Error()),
					})
				}

				golog.Instance().LogDebugWithFields(gologC.Fields{
					"method":  reflectionHelpers.GetThisFuncName(),
					"message": fmt.Sprintf("JM-SSH: SSH-DATA: %s", string(data)),
				})

				// Make `http.Request`
				request, err := http.NewRequest("POST", apiEndpoing.Value, bytes.NewBuffer(data))
				if err != nil {
					golog.Instance().LogErrorWithFields(gologC.Fields{
						"method":  reflectionHelpers.GetThisFuncName(),
						"message": fmt.Sprintf("JM-SSH: SSH-DATA-ERROR: %s", err.Error()),
					})
					break
				}

				route := thisRef.router.Get(apiEndpoing.Value)
				if route == nil {
					golog.Instance().LogErrorWithFields(gologC.Fields{
						"method":  reflectionHelpers.GetThisFuncName(),
						"message": fmt.Sprintf("JM-SSH: Missing ROUTE: %s", apiEndpoing.Value),
					})
					break
				}

				route.GetHandler().ServeHTTP(&customResponseWriter{sshChannel: channel}, request)
				break
			}
		}()
	}
}
