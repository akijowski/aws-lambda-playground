package serverless

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// LambdaOpts is a collection of options for configuring an AWS Serverless Lambda Function.
type LambdaOpts struct {
	FunctionName        string
	FunctionDescription string
	CodeURI             string
	Handler             string
}

// func NewSAMLambdaFunction(scope constructs.Construct, id *string, opts *LambdaOpts) awssam.CfnFunction {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fullCodeURI := filepath.Join(cwd, opts.CodeURI)
// 	l := awssam.NewCfnFunction(scope, id, &awssam.CfnFunctionProps{
// 		Runtime:      awslambda.Runtime_GO_1_X().Name(),
// 		Timeout:      jsii.Number(5),
// 		FunctionName: jsii.String(opts.FunctionName),
// 		Description:  jsii.String(opts.FunctionDescription),
// 		CodeUri:      awslambda.NewAssetCode(jsii.String(fullCodeURI), nil).Path(),
// 		Handler:      jsii.String(opts.Handler),
// 	})

// 	if opts.CreateLogGroup {
// 		awslogs.NewCfnLogGroup(scope, jsii.String(fmt.Sprintf("%sLogGroup", *id)), &awslogs.CfnLogGroupProps{
// 			LogGroupName:    jsii.String(fmt.Sprintf("aws/lambda/%s", opts.FunctionName)),
// 			RetentionInDays: jsii.Number(opts.LogRetentionDays),
// 		})
// 	}
// 	return l
// }

func NewLambdaFunction(scope constructs.Construct, id *string, opts *LambdaOpts) awslambda.Function {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fullCodeURI := filepath.Join(cwd, opts.CodeURI)
	return awslambda.NewFunction(scope, id, &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(5)),
		FunctionName: jsii.String(opts.FunctionName),
		Description:  jsii.String(opts.FunctionDescription),
		Code:         awslambda.NewAssetCode(jsii.String(fullCodeURI), nil),
		Handler:      jsii.String(opts.Handler),
		LogRetention: awslogs.RetentionDays_FIVE_DAYS,
		Tracing:      awslambda.Tracing_ACTIVE,
	})
}
