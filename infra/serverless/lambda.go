package serverless

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

var (
	GO_Runtime       = awslambda.Runtime_GO_1_X()
	PROVIDED_Runtime = awslambda.Runtime_PROVIDED_AL2()
)

// LambdaOpts is a collection of options for configuring an AWS Lambda Function.
//
// CodeURI is the relative path from the project root to the source directory.
// Handler is the name of the executable that will run in the Lambda environment.
type LambdaOpts struct {
	FunctionName        string
	FunctionDescription string
	CodeURI             string
	Handler             string
	Runtime             awslambda.Runtime
}

// LocalBundler satisfies the [awscdk.ILocalBundling] interface by implementing TryBundle.
// It will run go build so that the CDK can bundle the function assets correctly.
//
// The implementation of the interface must be a struct as it must be a marshalable object in Go,
// therefore a functional approach is not possible.
type LocalBundler struct {
	CodeUri string
	Handler string
}

// TryBundle is used to build the Lambda function in the local environment.  It returns `false` if unable to build the function,
// and `true` if it was successful.  It will place the build output in the provided `outputDir`.
//
// For more information see: https://aws.amazon.com/blogs/devops/building-apps-with-aws-cdk/
func (lb *LocalBundler) TryBundle(outputDir *string, options *awscdk.BundlingOptions) *bool {
	checkCmd := exec.Command("go", "version")
	buildCmd := exec.Command("go", "build", "-o", fmt.Sprintf("%s/%s", *outputDir, lb.Handler), lb.CodeUri)
	buildCmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")

	if err := checkCmd.Run(); err != nil {
		fmt.Printf("check error: %s\n", err)
		return jsii.Bool(false)
	}

	outB, err := buildCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("build output: %s\n", string(outB))
		fmt.Printf("build error: %q\n", err)
		return jsii.Bool(false)
	}
	return jsii.Bool(true)
}

// NewLambdaFunction builds a Go Lambda with sensible defaults.
func NewLambdaFunction(scope constructs.Construct, id *string, opts *LambdaOpts) awslambda.Function {

	l := awslambda.NewFunction(scope, id, &awslambda.FunctionProps{
		Runtime:      opts.Runtime,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(5)),
		FunctionName: jsii.String(opts.FunctionName),
		Description:  jsii.String(opts.FunctionDescription),
		Handler:      jsii.String(opts.Handler),
		Code: awslambda.NewAssetCode(jsii.String(mustCwd()), &awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				// See: https://github.com/aws/aws-cdk/issues/20907
				Image:   opts.Runtime.BundlingImage(),
				Command: jsii.Strings("bash", "-c", fmt.Sprintf("GOOS=linux go build -o /asset-output/%s %s", opts.Handler, opts.CodeURI)),
				Local:   &LocalBundler{CodeUri: opts.CodeURI, Handler: opts.Handler},
				User:    jsii.String("root"),
			},
		}),
		LogRetention: awslogs.RetentionDays_FIVE_DAYS,
		Tracing:      awslambda.Tracing_ACTIVE,
	})

	awscdk.NewCfnOutput(scope, jsii.String(fmt.Sprintf("%sARN", *id)), &awscdk.CfnOutputProps{
		Description: jsii.String(fmt.Sprintf("Lambda Function ARN for %s", opts.FunctionName)),
		Value:       l.FunctionArn(),
	})

	return l
}
