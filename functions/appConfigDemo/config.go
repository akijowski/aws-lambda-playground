package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func getConfig(ctx context.Context) (*SampleConfig, error) {
	var config *SampleConfig
	b, err := doGet(ctx)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func doGet(ctx context.Context) ([]byte, error) {
	url := "foo"
	host := "http://localhost:2772"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", host, url), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b := make([]byte, 0)
	if _, err := resp.Body.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
