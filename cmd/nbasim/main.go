package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/seaneyre/nbasim/internal/simulation"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sim *simulation.Simulation
var simMutex sync.Mutex

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	nba_game_id := "0022000180" //TODO: make this a flag which defaults to a random game
	sim = simulation.New(nba_game_id, 4.00, time.Now().Add(time.Second*2))

	go func() {
		if err := sim.Run(); err != nil {
			log.Printf("Simulation error: %v", err)
		}
	}()

	r := mux.NewRouter()
	path := "/ws/game/" + nba_game_id
	r.HandleFunc(path, handleWebSocket)

	port := "8080"
	log.Printf("Starting server on port %s", port)
	return http.ListenAndServe(":"+port, r)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Received WebSocket connection request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	simMutex.Lock()
	sim.AddConnection(conn)
	simMutex.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
	}

	simMutex.Lock()
	sim.RemoveConnection(conn)
	simMutex.Unlock()
}
