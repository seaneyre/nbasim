package main

import (
	"fmt"
	"os"
	"time"

	"github.com/seaneyre/nbasim/internal/simulation"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	nba_game_id := "0022000180" //TODO: make this a flag which defaults to a random game
	sim := simulation.New(nba_game_id, 4.00, time.Now().Add(time.Second * 2))
	return sim.Run()
}
