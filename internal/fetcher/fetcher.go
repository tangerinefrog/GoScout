package fetcher

import (
	"context"
	"io"
	"net/http"
)

func Fetch(url string) (body []byte, err error) {
	client := http.Client{}

	req, err := http.NewRequestWithContext(context.TODO(), "GET", url, nil)
	configureMockHeaders(req)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
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
