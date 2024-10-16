package simulation

import (
	"net/url"
	"os"
	"sort"
	"time"

	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/seaneyre/nbasim/internal/retrieve"
)

type Simulation struct {
	nba_game_id               string
	time_factor               float64
	real_start_time           time.Time
	simulated_game_clock_time int
	host_url                  string
}

func New(nbaGameID string, timeFactor float64, realStartTime time.Time, hostURL string) *Simulation {
	return &Simulation{
		nba_game_id:               nbaGameID,
		time_factor:               timeFactor,
		real_start_time:           realStartTime,
		host_url:                  hostURL,
		simulated_game_clock_time: 0,
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
}

func (s *Simulation) Run() error {
	log.Print("Running Simulation")
	log.Printf("NBA Game ID: %s", s.nba_game_id)
	log.Printf("Time factor: %f", s.time_factor)
	log.Printf("Real Start Time: %s", s.real_start_time.Format(time.RFC3339))

	serverURL := url.URL{Scheme: "ws", Host: s.host_url, Path: "/ws/game/" + s.nba_game_id, RawQuery: "type=simulator"}
	log.Info().Str("server_url", serverURL.String()).Msg("Dialing to Websocket with URL")
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal().Msgf("Error connecting to server: %v", err)
	}
	defer conn.Close()
	log.Info().Msgf("Connected to server for game ID: %s", s.nba_game_id)

	log.Info().Msg("Getting Play-By-Play data from NBA API")
	resp, err := retrieve.GetPlayByPlayResponse(s.nba_game_id)
	if err != nil {
		return err
	}

	log.Info().Msg("Preparing list of events from Play-By-Play response")
	events := PrepareEvents(resp)
	log.Info().Msgf("%d events prepared", len(events))

	log.Info().Msgf("now=%s", time.Now().Format(time.RFC3339))
	log.Info().Msgf("real_start_time=%s", s.real_start_time.Format(time.RFC3339))
	game_should_have_started := time.Now().Compare(s.real_start_time)
	log.Info().Msgf("game_should_have_started=%d", game_should_have_started)

	log.Info().Msg("Waiting for game start...")
	for time.Now().Compare(s.real_start_time) == -1 {
		time_to_go := s.real_start_time.Sub(time.Now())
		log.Info().Msgf("%f seconds until game starts", time_to_go.Seconds())
		sleep_duration := time_to_go.Seconds()
		if sleep_duration > 0 {
			log.Info().Msgf("Sleeping for %f seconds", sleep_duration)
			time.Sleep(time.Duration(sleep_duration * float64(time.Second)))
		}
	}

	log.Info().Msg("Starting game")
	for i, event := range events {
		log.Info().Msgf("Processing event %d/%d", i, len(events))
		log.Info().Msgf("The simulated clock time is currently %d. The event will happen at %d", s.simulated_game_clock_time, event.GameClockTime)
		for s.simulated_game_clock_time != event.GameClockTime {
			game_seconds_until_event := event.GameClockTime - s.simulated_game_clock_time
			log.Info().Msgf("%d seconds until event happens", game_seconds_until_event)
			time.Sleep(time.Duration(int(time.Second) * (game_seconds_until_event / int(s.time_factor))))
			s.simulated_game_clock_time += game_seconds_until_event
			log.Printf("Sleep finished s.simulated_game_clock_time=%d", s.simulated_game_clock_time)
		}
		log.Info().Msg("*Event Happens*")
		SendEvent(event, conn)
	}
	log.Info().Msg("Game finished")
	return nil
}

func SendEvent(event Event, conn *websocket.Conn) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal event")
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, eventJSON)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message through WebSocket")
	}
}

func PrepareEvents(resp retrieve.PlayByPlayResponse) []Event {
	var events []Event
	for _, action := range resp.Game.Actions {
		game_clock_time, err := GetGameClockTime(action.Clock, action.Period)
		if err != nil {
			log.Printf("Error getting Game Clock Time from clock=%s; period=%d", action.Clock, action.Period)
		}
		event := Event{
			GameClockTime: game_clock_time,
			ActionType:    action.ActionType,
			Action:        action,
		}
		events = append(events, event)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].GameClockTime < events[j].GameClockTime
	})
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
