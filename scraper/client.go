package scraper

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func Download(spec string) (io.ReadCloser, error) {
	ctx := context.Background()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.fishbase.se/summary/"+spec+".html",
		nil,
	)
	if err != nil {
		log.Printf("error creating request: %s", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error making http request: %s", err)
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("error: unexpected status code %d", resp.StatusCode)
		return nil, fmt.Errorf("error: unexpected status code: %d", resp.StatusCode)
	}
	return resp.Body, nil
}
