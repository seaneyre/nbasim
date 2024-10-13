package retrieve

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
}

func FetchResponseFromURL(url string) ([]byte, error) {
	log.Print("Getting response from URL:", url)

	resp, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return body, nil
}

func GetPlayByPlayResponse(nba_game_id string) (PlayByPlayResponse, error) {
	url := fmt.Sprintf("https://cdn.nba.com/static/json/liveData/playbyplay/playbyplay_%s.json", nba_game_id)

	j, err := FetchResponseFromURL(url)
	if err != nil {
		return PlayByPlayResponse{}, err
	}

	var events PlayByPlayResponse
	if err := json.Unmarshal(j, &events); err != nil {
		log.Err(err).Msg("Can not unmarshal JSON")
	}
	return events, nil
}
