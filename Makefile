.PHONY: build ls doctor synth test

PROJECT_NAME=aws-lambda-playground
AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --profile adam | jq -r .Account)
REGION=us-east-2

synth:
	npx aws-cdk synth -q

synth-sam:
	npx aws-cdk synth --no-staging

clean:
	rm -rf ./cdk.out/ ./build/ template.cdk.yaml

deploy: synth
	npx aws-cdk deploy --profile adam --outputs-file ./cdk-outputs.json

destroy: synth
	npx aws-cdk destroy --profile adam

bootstrap:
	npx aws-cdk bootstrap --profile adam aws://$(AWS_ACCOUNT_ID)/$(REGION)

ls:
	npx aws-cdk ls --long

cdk-help:
	npx aws-cdk -h

doctor:
	npx aws-cdk doctor

diff:
	npx aws-cdk diff --profile adam

test-stack:
	go test ./infra -tags=stack

invoke-helloworld-local: synth-sam
	sam local invoke --no-event \
	HelloWorldFunction

invoke-app-config-local: synth-sam
	sam local invoke --no-event \
	AppConfigDemoFunction \
	--env-vars local/app-config-lambda-env.json \
	--profile adam \
	--region us-east-2

invoke-arm-local: synth-sam
	sam local invoke --no-event \
	ArmLambdaDemo \

invoke-helloworld:
	aws lambda invoke \
	--function-name $$(cat ./cdk-outputs.json| jq --raw-output '."aws-lambda-playground"."HelloWorldFunctionARN"') \
	--payload '{}' \
	--profile adam \
	--no-cli-pager \
	response.json

invoke-app-config:
	aws lambda invoke \
	--function-name $$(cat ./cdk-outputs.json| jq --raw-output '."aws-lambda-playground"."AppConfigDemoFunctionARN"') \
	--payload '{}' \
	--profile adam \
	--no-cli-pager \
	response.json

invoke-arm:
	aws lambda invoke \
	--function-name $$(cat ./cdk-outputs.json| jq --raw-output '."aws-lambda-playground"."ArmLambdaDemoARN"') \
	--payload '{}' \
	--profile adam \
	--no-cli-pager \
	response.json
