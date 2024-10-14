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
	connections map[string]map[*websocket.Conn]bool
	mutex       sync.Mutex
}

func New(host string, port int) *Server {
	return &Server{
		host:        host,
		port:        port,
		connections: make(map[string]map[*websocket.Conn]bool),
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
	r.HandleFunc("/ws/game/{game_id}", s.handleWebSocket)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Info().Str("address", addr).Msg("Starting server")
	return http.ListenAndServe(addr, r)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["game_id"]
	log.Debug().Str("game_id", gameID).Msg("Received WebSocket connection request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error upgrading to WebSocket")
		return
	}
	defer conn.Close()
	log.Debug().Msg("WebSocket connection established")

	log.Debug().Msg("Adding connection to server connection list")
	s.mutex.Lock()
	if s.connections[gameID] == nil {
        s.connections[gameID] = make(map[*websocket.Conn]bool)
    }
	s.connections[gameID][conn] = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.connections[gameID], conn)
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
		s.Broadcast(gameID, messageType, message)
		log.Debug().Msg("Message broadcast complete")
	}
}

func (s *Server) Broadcast(gameId string, messageType int, message []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn := range s.connections[gameId] {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			log.Error().Err(err).Msg("Failed to broadcast message to a client")
			delete(s.connections[gameId], conn)
			conn.Close()
		}
	}
}
