package retrieve

import (
	"fmt"
	"net/http"
	"io"
	"time"
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