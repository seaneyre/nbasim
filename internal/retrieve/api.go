package retrieve

import (
	"fmt"
	"net/http"
	"io"
	"time"

	"encoding/json"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func FetchResponseFromURL(url string) ([]byte, error) {
	fmt.Println("Fetching JSON from URL:", url)

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
	fmt.Println("Fetching events from:", url)

	j, err := FetchResponseFromURL(url)
	if err != nil {
		return PlayByPlayResponse{}, err
	}

	var events PlayByPlayResponse
	if err := json.Unmarshal(j, &events); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	return events, nil
}
