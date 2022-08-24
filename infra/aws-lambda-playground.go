//go:build stack

package main

import (
	"fmt"

	"github.com/akijowski/aws-lambda-playground/infra/config"
	"github.com/akijowski/aws-lambda-playground/infra/serverless"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
)

type AwsLambdaPlaygroundStackProps struct {
	awscdk.StackProps
}

func NewAwsLambdaPlaygroundStack(scope constructs.Construct, id string, props *AwsLambdaPlaygroundStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	serverless.NewLambdaFunction(stack, jsii.String("HelloWorldFunction"), &serverless.LambdaOpts{
		FunctionName:        fmt.Sprintf("%s-hello-world", *stack.StackName()),
		FunctionDescription: "Simple Hello World Lambda",
		CodeURI:             "./functions/helloWorld",
		Handler:             "helloWorld",
	})
	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("AwsLambdaPlaygroundQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	config.NewAppConfig(stack, "AppConfigDemo").
		WithApplication(&config.AppConfigApplicationInput{
			Name:        "app-config-demo",
			Description: "AppConfig Demo project",
		}).
		WithEnvironment(&config.AppConfigEnvironmentInput{
			Name:        "app-config-demo",
			Description: "AppConfig Demo environment",
		}).
		WithHostedFreeformProfile(&config.AppConfigProfileInput{
			Name:        "app-config-demo",
			Description: "AppConfig Demo profile",
		}).
		WithSimpleDeployStrategy("app-config-demo").
		WithHostedConfiguration(&config.AppConfigConfigurationInput{
			Description: "AppConfig Demo configuration",
			ContentType: "application/json",
		}).
		Deploy("a deployment for AppConfig Demo")

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	stackName := app.Node().TryGetContext(jsii.String("stack-name")).(string)

	NewAwsLambdaPlaygroundStack(app, stackName, &AwsLambdaPlaygroundStackProps{
		awscdk.StackProps{
			Description: jsii.String("Lambda Samples"),
			Env:         env(),
			Tags: &map[string]*string{
				"Project": jsii.String("aws-lambda-playground"),
			},
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
