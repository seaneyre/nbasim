package server

import (
	"encoding/json"
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

type ConnectionType int

const (
	SimulatorConnection ConnectionType = iota
	ListenerConnection
)

type Connection struct {
	conn     *websocket.Conn
	connType ConnectionType
}

type Server struct {
	host       string
	port       int
	listeners  map[string]map[*Connection]bool
	simulators map[string]map[*Connection]bool
	mutex      sync.Mutex
}

func New(host string, port int) *Server {
	return &Server{
		host:       host,
		port:       port,
		listeners:  make(map[string]map[*Connection]bool),
		simulators: make(map[string]map[*Connection]bool),
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
	r.HandleFunc("/api/status", s.handleAPIStatus).Methods("GET")

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Info().Str("address", addr).Msg("Starting server")
	return http.ListenAndServe(addr, r)
}

func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	simulationGameIDs := make([]string, 0, len(s.simulators))
	for gameID := range s.simulators {
		simulationGameIDs = append(simulationGameIDs, gameID)
	}
	status := map[string]interface{}{
		"server_status":          "running",
		"simulation_count":       len(s.simulators),
		"listener_count":         len(s.listeners),
		"total_connection_count": len(s.simulators) + len(s.listeners),
		"simulation_game_ids":    simulationGameIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["game_id"]
	connType := r.URL.Query().Get("type")

	log.Debug().Str("game_id", gameID).Str("type", connType).Msg("Received WebSocket connection request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error upgrading to WebSocket")
		return
	}
	defer conn.Close()
	log.Debug().Msg("WebSocket connection established")

	log.Debug().Msg("Parsing connection type")
	var connectionType ConnectionType
	switch connType {
	case "simulator":
		connectionType = SimulatorConnection
		log.Debug().Msg("Creating Simulator Connection")

		connection := &Connection{
			conn:     conn,
			connType: connectionType,
		}

		log.Debug().Msg("Adding connection to server simulation list")
		s.mutex.Lock()
		if s.simulators[gameID] == nil {
			s.simulators[gameID] = make(map[*Connection]bool)
		}
		s.simulators[gameID][connection] = true
		s.mutex.Unlock()

		defer func() {
			s.mutex.Lock()
			delete(s.simulators[gameID], connection)
			s.mutex.Unlock()
			conn.Close()
		}()
	case "listener":
		connectionType = ListenerConnection
		log.Debug().Msg("Creating Listener Connection")

		connection := &Connection{
			conn:     conn,
			connType: connectionType,
		}

		log.Debug().Msg("Adding connection to server listener list")
		s.mutex.Lock()
		if s.listeners[gameID] == nil {
			s.listeners[gameID] = make(map[*Connection]bool)
		}
		s.listeners[gameID][connection] = true
		s.mutex.Unlock()

		defer func() {
			s.mutex.Lock()
			delete(s.listeners[gameID], connection)
			s.mutex.Unlock()
			conn.Close()
		}()
	case "":
		connectionType = ListenerConnection
		log.Debug().Msg("Creating Default (Listener) Connection")

		connection := &Connection{
			conn:     conn,
			connType: connectionType,
		}

		log.Debug().Msg("Adding connection to server listener list")
		s.mutex.Lock()
		if s.listeners[gameID] == nil {
			s.listeners[gameID] = make(map[*Connection]bool)
		}
		s.listeners[gameID][connection] = true
		s.mutex.Unlock()

		defer func() {
			s.mutex.Lock()
			delete(s.listeners[gameID], connection)
			s.mutex.Unlock()
			conn.Close()
		}()
	default:
		log.Error().Str("type", connType).Msg("Invalid connection type")
		conn.Close()
		return
	}

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
		s.BroadcastToListeners(gameID, messageType, message)
		log.Debug().Msg("Message broadcast complete")
	}
}

func (s *Server) BroadcastToListeners(gameId string, messageType int, message []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn := range s.listeners[gameId] {
		err := conn.conn.WriteMessage(messageType, message)
		if err != nil {
			log.Error().Err(err).Msg("Failed to broadcast message to a client")
			delete(s.listeners[gameId], conn)
			conn.conn.Close()
		}
	}
}
