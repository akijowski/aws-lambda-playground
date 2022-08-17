package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

var logger *log.Logger

func handler(ctx context.Context) (string, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		logger.Printf("%+v\n", lc)
	}
	return "Hello World!\n", nil
}

func main() {
	logger = log.Default()
	logger.SetPrefix("hello_world ")
	logger.SetFlags(log.Lshortfile)
	lambda.Start(handler)
}
