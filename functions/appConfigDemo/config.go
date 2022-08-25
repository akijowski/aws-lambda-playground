package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getConfig(ctx context.Context, path string) (*SampleConfig, error) {
	var config *SampleConfig
	b, err := doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func doGet(ctx context.Context, path string) ([]byte, error) {
	host := "http://localhost:2772"
	url := fmt.Sprintf("%s/%s", host, path)
	logger.Printf("Making config request to: %s\n", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	logger.Printf("status code: %q\n", resp.Status)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
