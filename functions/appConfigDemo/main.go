package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

var (
	logger      *log.Logger
	isColdStart bool
)

type SampleConfig struct {
	Enabled      bool   `json:"enabled"`
	StringName   string `json:"string_name"`
	NumericValue int    `json:"numeric_value"`
}

type Response struct {
	Config    *SampleConfig `json:"config"`
	ColdStart bool          `json:"cold_start"`
}

func handler(ctx context.Context) (*Response, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		logger.Printf("%+v\n", lc)
	}
	config, err := getConfig(ctx)
	if err != nil {
		return nil, err
	}
	resp := &Response{ColdStart: isColdStart, Config: config}
	isColdStart = false
	return resp, nil
}

func main() {
	logger = log.Default()
	logger.SetPrefix("hello_world ")
	logger.SetFlags(log.Lshortfile)
	isColdStart = true
	lambda.Start(handler)
}
