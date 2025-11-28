package fetcher

import (
	"context"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

var ErrorUnsuccessfulStatusCode error = errors.New("unsuccessful status code")

func Fetch(ctx context.Context, url string) (body []byte, err error) {
	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	configureMockHeaders(req)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, ErrorUnsuccessfulStatusCode
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FetchWithRetry(ctx context.Context, url string, retryCount int) (body []byte, err error) {
	for i := range retryCount {
		//delay formula: 2^i seconds
		delay := time.Duration(math.Pow(2, float64(i))) * time.Second

		body, err = Fetch(ctx, url)
		if err != nil {
			if errors.Is(err, ErrorUnsuccessfulStatusCode) {
				log.Printf("got unsuccessful status code from '%s'", url)
				time.Sleep(delay)
				continue
			}
			return
		}
		return
	}

	return
}

func configureMockHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0") // Add a new header
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("DNT", "1")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
}
