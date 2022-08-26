package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

var logger *log.Logger

func handler(ctx context.Context) error {
	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		logger.Printf("%+v\n", lc)
	}
	logger.Println("This is a performant ARM64 Lambda!")
	return nil
}

func main() {
	logger = log.Default()
	logger.SetPrefix("arm_lambda ")
	logger.SetFlags(log.Lshortfile)
	lambda.Start(handler)
}
