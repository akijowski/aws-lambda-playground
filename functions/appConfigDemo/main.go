package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

var (
	logger      *log.Logger
	isColdStart bool
)

type SampleConfig struct {
	Enabled      bool   `json:"enabled"`
	StringValue  string `json:"string_value"`
	NumericValue int    `json:"numeric_value"`
}

type Response struct {
	Config    *SampleConfig `json:"config"`
	ColdStart bool          `json:"cold_start"`
	InvokedAt time.Time     `json:"invoked_at"`
}

func handler(ctx context.Context) (*Response, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		logger.Printf("%+v\n", lc)
	}
	configPath := os.Getenv("APP_CONFIG_PATH")
	if configPath == "" {
		return nil, errors.New("APP_CONFIG_PATH is not set")
	}
	config, err := getConfig(ctx, configPath)
	if err != nil {
		return nil, err
	}
	resp := &Response{ColdStart: isColdStart, Config: config, InvokedAt: time.Now()}
	isColdStart = false
	return resp, nil
}

func main() {
	logger = log.Default()
	logger.SetPrefix("app_config_demo ")
	logger.SetFlags(log.Lshortfile)
	isColdStart = true
	lambda.Start(handler)
}
