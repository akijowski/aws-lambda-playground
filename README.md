# AWS Lambda Playground

This is intended to be a sample proving ground for different aspects of AWS Lambda.

1. :white_check_mark: AWS CDK + AWS SAM
1. :white_check_mark: AWS AppConfig Lambda Extension
1. Caching Lambda Extension
1. :white_check_mark: Building Go Lambda in a container with AWS CDK
1. :white_check_mark: Building a Go Lambda to use ARM Lambda instances
1. Blue/Green deployments with CodeDeploy and Lambda

## AWS CDK + AWS SAM

With CDK, it appears that SAM takes on a smaller role in development.  The CDK remains the tool that manages creating the template and deploying the CloudFormation stack.  SAM is really only needed when wanting to run local instances of a Lambda for testing.

Further information of how SAM works with CDK can be found [in the SAM documentation](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-cdk.html)

### CDK SAM Constructs

The CDK library contains only L1 constructs for AWS SAM.  Researching Github Issues and documentation, the sense is that any benefit to using the "Serverless CloudFormation Transform" and Serverless Resources (`AWS::Serverless::Function`, etc) could be replicated within CDK.  Therefore, using these SAM L1 constructs is unnecessary and the L2 constructs are a more effective method to defining resources.

Github Issue thread talking about how CDK replaces or can't replicate some of the features of SAM [#716](https://github.com/aws/aws-cdk/issues/716)

### Make Commands

To run a local Lambda container with AWS SAM:

```bash
make invoke-helloworld-local
```

To invoke the Lambda in AWS:

```bash
make invoke-helloworld
```

## AWS AppConfig Lambda Extension

[AWS AppConfig provides an offical Lambda Extension](https://docs.aws.amazon.com/appconfig/latest/userguide/appconfig-integration-lambda-extensions.html).  It handles invoking and caching configuration data for your Lambda.  This can have many benefits:

- Provides a central location to manage configuration instead of environment variables within Lambda
- Configuration changes do not require a deploy to the Lambda
- The extension, through caching, can save requests to AppConfig
- AppConfig also supports feature flags, which open up new opportunities for deployment and application behavior

It is not all :sunny: sunshine though.  There are some issues to take in to consideration.

### Many AppConfig Resources

There are a few AppConfig resources that need to be made and coordinated.

- **Application**: This is needed as a virtual resource to namespace the remaining resources
- **Environment**: Another resource to namespace configuration
- **Configuration Profile**: A resource that references a configuration data source
- **Hosted Configuration**: If you want AppConfig to host the config data, you will need one of these
- **Deployment Strategy**: You must determine how AppConfig deploys configuration data to targets
- **Deployment**: This is needed to have AppConfig deploy changes

### CDK Constructs

Currently there are no L2 CDK constructs that simplify building AppConfig resources.  This means you will need to create the above resources with low-level constructs.

### Extension Runtime Support

Looking at the [supported runtimes](https://docs.aws.amazon.com/appconfig/latest/userguide/appconfig-integration-lambda-extensions.html#appconfig-integration-lambda-extensions-runtimes), the Lambda extension does not support Go.  This means the Lambda will need to be running a Provided runtime.

### Make Commands

To run a local Lambda container with AWS SAM:

```bash
make invoke-app-config-local
```

To invoke the Lambda in AWS:

```bash
make invoke-app-config
```

## Caching Lambda Extension

TODO: https://github.com/aws-samples/aws-lambda-extensions/tree/main/cache-extension-demo

## Building Go Lambda in a container with AWS CDK

To build a Go Lambda in a container, it requires providing bundling options when declaring Lambda asset code in the `awslambda.Function` constructor.
To accomplish this I modified the options my `NewLambdaFunction` uses to include a boolean field.  When set it will set the `localBundler` to `nil`, which effectively forces a container build.
The alternative is for the `localBundler` to return a `false` value, which tells the CDK that local bundling failed and will fall back to a container build.

## Building a Go Lambda to use ARM Lambda instances

To build a Go Lambda for ARM64 architectures, two things must be configured:

1. The compiled Go code must be compatible with Linux/ARM64
1. The runtime must be Provided.AL2

The first requirement is accomplished by calling `go build` with `GOOS=linux` and `GOARCH=arm64`.
The second is changing the Lambda runtime to `provided`.

In this case, the function code is built in the Go runtime bundling Docker image, which provides all the tooling needed to build a function binary.
The function is then configured to load the binary in a Provided.AL2 Lambda runtime environment.

### Make Commands

To run a local Lambda container in AWS SAM:

```bash
make invoke-arm-local
```

To invoke the Lambda in AWS:

```bash
make invoke-arm
```
