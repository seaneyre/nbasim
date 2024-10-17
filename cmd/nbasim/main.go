package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/seaneyre/nbasim/internal/simulation"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
}

func usage() {
	intro := `nbasim is a framework to simulate NBA games play-by-play and serve events to a Websocket API.

Usage:
  nbasim [flags] <command> [command flags]`
	fmt.Fprintln(os.Stderr, intro)

	fmt.Fprintln(os.Stderr, "\nCommands:")
	fmt.Fprintln(os.Stderr, "  simulate")

	fmt.Fprintln(os.Stderr, "\nFlags:")
	// Prints a help string for each flag

	flag.PrintDefaults()

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Run `nbasim <command> -h` to get help for a specific command\n\n")
}

func main() {
	log.Debug().Msg("Starting nbasim")

	flag.Usage = usage
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverHost := serverCmd.String("host", "localhost", "Host address for the server")
	serverPort := serverCmd.Int("port", 8000, "Port number for the server")

	simulateCmd := flag.NewFlagSet("simulate", flag.ExitOnError)
	simulateURL := simulateCmd.String("url", "localhost:8000", "Host URL for the server.")
	simulateGameID := simulateCmd.String("game-id", "0022000180", "NBA API Game ID to simulate")
	simulateTimeFactor := simulateCmd.Float64("time-factor", 4.00, "Time factor to run the simulation at (e.g. if 2 then the simulation will run at 2x the speed of real-time)")

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Debug().Msg("No cli args specified, printing usage and exit 1")
		flag.Usage()
		os.Exit(1)
	}

	subCmd := flag.Arg(0)
	subCmdArgs := flag.Args()[1:]

	switch subCmd {
	case "simulate":
		simulateCmd.Parse(subCmdArgs)
		sim := simulation.New(*simulateGameID, *simulateTimeFactor, time.Now().Add(time.Second*2), *simulateURL)
		sim.Run()
	default:
		log.Debug().Msgf("Unknown subcommand provided: %s", subCmd)
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(subCmd, subCmdArgs)
}
