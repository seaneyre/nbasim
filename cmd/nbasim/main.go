package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/seaneyre/nbasim/internal/server"
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
	fmt.Fprintln(os.Stderr, "  server\n  simulator")

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
	serverPort := serverCmd.Int("port", 8080, "Port number for the server")

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Debug().Msg("No cli args specified, printing usage and exit 1")
		flag.Usage()
		os.Exit(1)
	}

	subCmd := flag.Arg(0)
	subCmdArgs := flag.Args()[1:]

	switch subCmd {
	case "server":
		serverCmd.Parse(subCmdArgs)
		srv := server.New(*serverHost, *serverPort)
		srv.Run()
	case "simulate":

		gameID := "0022000180"
		sim := simulation.New(gameID, 4.00, time.Now().Add(time.Second*2))
		sim.Run()
	default:
		log.Debug().Msgf("Unknown subcommand provided: %s", subCmd)
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(subCmd, subCmdArgs)
}
