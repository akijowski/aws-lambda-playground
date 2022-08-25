//go:build stack

package main

import (
	"fmt"

	"github.com/akijowski/aws-lambda-playground/infra/config"
	"github.com/akijowski/aws-lambda-playground/infra/serverless"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
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
		Runtime:             serverless.GO_Runtime,
	})

	// Using a provided runtime requires the handler to be named bootstrap
	configDemo := serverless.NewLambdaFunction(stack, jsii.String("AppConfigDemoFunction"), &serverless.LambdaOpts{
		FunctionName:        fmt.Sprintf("%s-app-config-demo", *stack.StackName()),
		FunctionDescription: "Lambda to test the App Config Layer or Extension",
		CodeURI:             "./functions/appConfigDemo",
		Handler:             "bootstrap",
		Runtime:             serverless.PROVIDED_Runtime,
	})

	appConfig := config.NewAppConfig(stack, "AppConfigDemo").
		WithApplication(&config.AppConfigApplicationInput{
			Name:        "app-config-demo",
			Description: "AppConfig Demo project",
		}).
		WithEnvironment(&config.AppConfigEnvironmentInput{
			Name:        "stage",
			Description: "AppConfig Demo environment",
		}).
		WithHostedFreeformProfile(&config.AppConfigProfileInput{
			Name:        "main",
			Description: "AppConfig Demo profile",
		}).
		WithSimpleDeployStrategy("allAtOnce").
		WithHostedConfiguration(&config.AppConfigConfigurationInput{
			Description: "AppConfig Demo configuration",
			ContentType: "application/json",
			ContentPath: "reference/config/config.json",
		}).
		Deploy("a deployment for AppConfig Demo")

	appConfigExtension := awslambda.LayerVersion_FromLayerVersionArn(stack, jsii.String("AppConfigDemoLayer"), jsii.String("arn:aws:lambda:us-east-2:728743619870:layer:AWS-AppConfig-Extension:50"))
	configDemo.AddLayers(appConfigExtension)
	configDemo.AddToRolePolicy(config.AppConfigDataPolicy())
	configDemo.AddEnvironment(jsii.String("APP_CONFIG_PATH"), appConfig.ConfigurationPath(), nil)
	configDemo.AddEnvironment(jsii.String("AWS_APPCONFIG_EXTENSION_LOG_LEVEL"), jsii.String("debug"), nil)

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
