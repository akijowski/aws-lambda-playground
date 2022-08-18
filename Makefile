.PHONY: build ls doctor synth test

PROJECT_NAME=aws-lambda-playground
AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --profile adam | jq -r .Account)
REGION=us-east-2

synth: build
	npx aws-cdk synth -q

synth-sam: build
	npx aws-cdk synth --no-staging > template.cdk.yaml

build:
	./bin/build-lambdas ./functions

bootstrap:
	npx aws-cdk bootstrap --profile adam aws://$(AWS_ACCOUNT_ID)/$(REGION)

ls:
	npx aws-cdk ls --long

doctor:
	npx aws-cdk doctor

diff:
	npx aws-cdk diff --profile adam

build-stack:
	go build -o=$(PROJECT_NAME) -tags=stack ./infra

test-stack:
	go test ./infra -tags=stack

invoke-helloworld: synth-sam
	sam local invoke --no-event \
	HelloWorldFunction
