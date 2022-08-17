package serverless

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// LambdaOpts is a collection of options for configuring an AWS Serverless Lambda Function.
//
// CodeURI needs to be relative to the calling function in order to be resolved correctly in the template.
type LambdaOpts struct {
	FunctionName        string
	FunctionDescription string
	CodeURI             string
	Handler             string
	CreateLogGroup      bool
	LogRetentionDays    float64
}

func NewLambdaFunction(scope constructs.Construct, id *string, opts *LambdaOpts) awssam.CfnFunction {
	l := awssam.NewCfnFunction(scope, id, &awssam.CfnFunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X().Name(),
		FunctionName: jsii.String(opts.FunctionName),
		Description:  jsii.String(opts.FunctionDescription),
		CodeUri:      awslambda.AssetCode_FromAsset(jsii.String(opts.CodeURI), nil).Path(),
		Handler:      jsii.String(opts.Handler),
	})

	if opts.CreateLogGroup {
		awslogs.NewCfnLogGroup(scope, jsii.String(fmt.Sprintf("%sLogGroup", *id)), &awslogs.CfnLogGroupProps{
			LogGroupName:    jsii.String(fmt.Sprintf("aws/lambda/%s", opts.FunctionName)),
			RetentionInDays: jsii.Number(opts.LogRetentionDays),
		})
	}
	return l
}
