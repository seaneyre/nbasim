package simulation

import (
	_ "fmt"
	"os"
	"time"

	"encoding/json"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/seaneyre/nbasim/internal/retrieve"
)

type Simulation struct {
	nba_game_id string
	time_factor float64
	real_start_time time.Time
}

func New(nbaGameID string, realStartTime time.Time) *Simulation {
	return &Simulation{
		nba_game_id: nbaGameID,
		time_factor: 1.0,
		real_start_time: realStartTime,
	}
}

func (s *Simulation) Run() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})

	log.Print("hello world")

	log.Print("Running simulation")
	log.Printf("NBA Game ID: %s", s.nba_game_id)
	log.Printf("Time factor: %f", s.time_factor)
	log.Printf("Real Start Time: %s", s.real_start_time.Format(time.RFC3339))
	
	log.Info().Msg("Getting Play-By-Play data from NBA API")
	resp, err := retrieve.GetPlayByPlayResponse(s.nba_game_id)
	if err != nil {
		return err
	}

	str, _ := json.MarshalIndent(resp.Game.Actions[0], "", "  ")
	log.Print(string(str))

	log.Info().Msg("Preparing list of events from Play-By-Play response")
	events := PrepareEvents(resp)
	log.Info().Msgf("%d events prepared", len(events))

	return nil
}

func PrepareEvents(resp retrieve.PlayByPlayResponse) []Event {
	var events []Event
	for _, action := range resp.Game.Actions {
		game_clock_time, err := GetGameClockTime(action.Clock, action.Period)
		if err != nil {
			log.Printf("Error getting Game Clock Time from clock=%s; period=%d", action.Clock, action.Period)
		}
		// log.Printf("clock=%s; period=%d; game_clock_time=%d\n", action.Clock, action.Period, game_clock_time)

		event := Event{
			GameClockTime: game_clock_time,
			ActionType: action.ActionType,
			Action: action,
		}
		events = append(events, event)
	}
	return events
}

func GetGameClockTime(clock_string string, period int) (int, error) {
	// TODO write tests for this!
	minutes, err := strconv.Atoi(clock_string[2:4])
	seconds, err := strconv.Atoi(clock_string[5:7])
	if err != nil {
		return 0, err
	}
	game_clock_time := ((12 * (period - 1)) * 60) + ((12 - minutes) * 60) + (60 - seconds) - 60
	return game_clock_time, nil
}