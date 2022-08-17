# AWS Lambda Playground

This is intended to be a sample proving ground for different aspects of AWS Lambda.

1. AWS CDK + AWS SAM
1. AWS AppConfig Lambda Extension
1. Caching Lambda Extension
1. Building Go Lambda in a container
1. Building a Go Lambda to use ARM Lambda instances

## AWS CDK + AWS SAM

With CDK, it appears that SAM takes on a smaller role in development.  The CDK remains the tool that manages creating the template and deploying the CloudFormation stack.  SAM is really only needed when wanting to run local instances of a Lambda for testing.

> *SAM may also be necessary for building Lambda binaries, this is undetermined at this time*

Further information of how SAM works with CDK can be found [in the SAM documentation](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-cdk.html)

Github Issue thread talking about how CDK replaces or can't replicate some of the features of SAM [#716](https://github.com/aws/aws-cdk/issues/716)
