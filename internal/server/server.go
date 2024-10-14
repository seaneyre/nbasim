package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
}

type Server struct {
	host string
	port int
}

func New(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) Run() error {
	log.Debug().Msg("Running Server")

	r := mux.NewRouter()
	path := "/ws/game/"
	r.HandleFunc(path, handleWebSocket)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Info().Str("address", addr).Msg("Starting server")
	return http.ListenAndServe(addr, r)
}



func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Received WebSocket connection request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error upgrading to WebSocket")
		return
	}
	defer conn.Close()

	log.Debug().Msg("WebSocket connection established")

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("WebSocket read error")
			} else {
				log.Debug().Err(err).Msg("WebSocket connection closed")
			}
			break
		}

		log.Info().Str("message", string(data)).Msg("Message Received")

		err = conn.WriteMessage(messageType, data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message through WebSocket")
			break
		}

		log.Debug().Str("message", string(data)).Msg("Message sent back to client")
	}
}
