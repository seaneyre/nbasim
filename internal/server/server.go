package server

import (
	"fmt"
	"net/http"
	"os"
	"sync"
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
	host        string
	port        int
	connections map[*websocket.Conn]bool
	mutex       sync.Mutex
}

func New(host string, port int) *Server {
	return &Server{
		host:        host,
		port:        port,
		connections: make(map[*websocket.Conn]bool),
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
	r.HandleFunc(path, s.handleWebSocket)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Info().Str("address", addr).Msg("Starting server")
	return http.ListenAndServe(addr, r)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Received WebSocket connection request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error upgrading to WebSocket")
		return
	}
	defer conn.Close()
	log.Debug().Msg("WebSocket connection established")

	log.Debug().Msg("Adding connection to server connection list")
	s.mutex.Lock()
	s.connections[conn] = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.connections, conn)
		s.mutex.Unlock()
		conn.Close()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("WebSocket read error")
			} else {
				log.Debug().Err(err).Msg("WebSocket connection closed")
			}
			break
		}
		log.Info().Msg("Message Received")
		log.Info().Msg("Broadcasting message")
		s.Broadcast(messageType, message)
		log.Debug().Msg("Message broadcast complete")
	}
}

func (s *Server) Broadcast(messageType int, message []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn := range s.connections {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			log.Error().Err(err).Msg("Failed to broadcast message to a client")
			delete(s.connections, conn)
			conn.Close()
		}
	}
}
