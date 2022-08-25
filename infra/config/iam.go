package config

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func AppConfigDataPolicy() awsiam.PolicyStatement {
	return awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Sid:    jsii.String("readAppConfigData"),
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"appconfig:StartConfigurationSession",
			"appconfig:GetLatestConfiguration",
		),
		Resources: jsii.Strings(
			"*",
		),
	})
}
