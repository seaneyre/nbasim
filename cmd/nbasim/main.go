package main

import (
	"fmt"
	"os"

	"github.com/seaneyre/nbasim/internal/simulation"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	sim := simulation.New()
	return sim.Run()
}
