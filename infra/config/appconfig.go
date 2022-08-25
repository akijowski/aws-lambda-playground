package config

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappconfig"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AppConfig struct {
	Scope              constructs.Construct
	LogicalIDPrefix    string
	Application        awsappconfig.CfnApplication
	Environment        awsappconfig.CfnEnvironment
	Configuration      awsappconfig.CfnHostedConfigurationVersion
	Profile            awsappconfig.CfnConfigurationProfile
	DeploymentStrategy awsappconfig.CfnDeploymentStrategy
	Deployment         awsappconfig.CfnDeployment
}

type AppConfigApplicationInput struct {
	Name        string
	Description string
}

type AppConfigEnvironmentInput struct {
	Name        string
	Description string
}

type AppConfigProfileInput struct {
	Name        string
	Description string
}

type AppConfigConfigurationInput struct {
	Description string
	ContentType string
	ContentPath string
}

func NewAppConfig(stack constructs.Construct, prefix string) *AppConfig {
	return &AppConfig{Scope: stack, LogicalIDPrefix: prefix}
}

func (ac *AppConfig) WithApplication(input *AppConfigApplicationInput) *AppConfig {
	app := awsappconfig.NewCfnApplication(ac.Scope, ac.logicalID("App"), &awsappconfig.CfnApplicationProps{
		Name:        jsii.String(input.Name),
		Description: jsii.String(input.Description),
	})
	ac.Application = app
	return ac
}

func (ac *AppConfig) WithEnvironment(input *AppConfigEnvironmentInput) *AppConfig {
	ac.Environment = awsappconfig.NewCfnEnvironment(ac.Scope, ac.logicalID("Env"), &awsappconfig.CfnEnvironmentProps{
		ApplicationId: ac.Application.Ref(),
		Name:          jsii.String(input.Name),
		Description:   jsii.String(input.Description),
	})
	return ac
}

func (ac *AppConfig) WithHostedFreeformProfile(input *AppConfigProfileInput) *AppConfig {
	ac.Profile = awsappconfig.NewCfnConfigurationProfile(ac.Scope, ac.logicalID("Profile"), &awsappconfig.CfnConfigurationProfileProps{
		ApplicationId: ac.Application.Ref(),
		Name:          jsii.String(input.Name),
		Description:   jsii.String(input.Description),
		Type:          jsii.String("AWS.Freeform"),
		LocationUri:   jsii.String("hosted"),
	})
	return ac
}

func (ac *AppConfig) WithSimpleDeployStrategy(name string) *AppConfig {
	ac.DeploymentStrategy = awsappconfig.NewCfnDeploymentStrategy(ac.Scope, ac.logicalID("Strategy"), &awsappconfig.CfnDeploymentStrategyProps{
		Name:                        jsii.String(name),
		DeploymentDurationInMinutes: jsii.Number(1),
		GrowthFactor:                jsii.Number(50),
		Description:                 jsii.String("Sample deploy strategy"),
		FinalBakeTimeInMinutes:      jsii.Number(1),
		GrowthType:                  jsii.String("LINEAR"),
		ReplicateTo:                 jsii.String("NONE"),
	})
	return ac
}

func (ac *AppConfig) WithHostedConfiguration(input *AppConfigConfigurationInput) *AppConfig {
	c := loadFromFile(input.ContentPath)
	ac.Configuration = awsappconfig.NewCfnHostedConfigurationVersion(ac.Scope, ac.logicalID("Config"), &awsappconfig.CfnHostedConfigurationVersionProps{
		ApplicationId:          ac.Application.Ref(),
		ConfigurationProfileId: ac.Profile.Ref(),
		Description:            jsii.String(input.Description),
		ContentType:            jsii.String(input.ContentType),
		Content:                c,
	})
	return ac
}

func (ac *AppConfig) Deploy(description string) *AppConfig {
	deploy := awsappconfig.NewCfnDeployment(ac.Scope, ac.logicalID("Deployment"), &awsappconfig.CfnDeploymentProps{
		Description:            jsii.String(description),
		ApplicationId:          ac.Application.Ref(),
		EnvironmentId:          ac.Environment.Ref(),
		DeploymentStrategyId:   ac.DeploymentStrategy.Ref(),
		ConfigurationProfileId: ac.Profile.Ref(),
		ConfigurationVersion:   ac.Configuration.Ref(),
	})

	ac.Deployment = deploy

	awscdk.NewCfnOutput(ac.Scope, jsii.String("AppConfigID"), &awscdk.CfnOutputProps{
		Description: jsii.String("ID for the AppConfig Application"),
		Value:       ac.Application.Ref(),
	})

	return ac
}

func (ac *AppConfig) ConfigurationPath() *string {
	path := fmt.Sprintf("applications/%s/environments/%s/configurations/%s", *ac.Application.Name(), *ac.Environment.Name(), *ac.Profile.Name())
	return jsii.String(path)
}

func (ac *AppConfig) logicalID(name string) *string {
	return jsii.String(fmt.Sprintf("%s%s", ac.LogicalIDPrefix, name))
}
