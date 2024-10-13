package simulation

import (
	"fmt"

	"github.com/seaneyre/nbasim/internal/retrieve"
	"encoding/json"
)

type Simulation struct {
	nba_game_id int
	time_factor float64
}

func New(nbaGameID int) *Simulation {
	return &Simulation{
		nba_game_id: nbaGameID,
		time_factor: 1.0,
	}
}

func (s *Simulation) Run() error {
	fmt.Println("Running simulation")
	fmt.Printf("NBA Game ID: %d\n", s.nba_game_id)
	fmt.Printf("Time factor: %f\n", s.time_factor)
	
	url := "https://cdn.nba.com/static/json/liveData/playbyplay/playbyplay_0022000180.json"
	fmt.Println("Fetching events from:", url)

	j, err := retrieve.FetchResponseFromURL(url)
	if err != nil {
		return err
	}

	var events retrieve.PlayByPlayResponse
	if err := json.Unmarshal(j, &events); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	str, _ := json.MarshalIndent(events.Game.Actions[0], "", "  ")
	fmt.Println(string(str))

	return nil
}
