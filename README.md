# AWS Step Functions Iteration Sample
A complete serverless stack demonstrating how to process an arbitrary number of
jobs using AWS Step Functions. This is the sample project to the article
[Processing an arbitrary number of jobs with AWS Step Functions](https://medium.com/@christianklotz/c185c2d2608).

## Prerequisites
In order to build, deploy and run this sample app you'll need an AWS account
and the following tools installed on your machine.

- [Go](https://golang.org/) to compile the Lambda functions
- [AWS CLI](https://aws.amazon.com/cli/) for deployment
- [AWS SAM CLI](https://github.com/awslabs/aws-sam-cli) for deployment

Make sure to clone the repository into your `$GOPATH`.

## Build
Get all dependencies and build the Lambda functions for the Amazon Linux.

    go get ./cmd/...

    GOOS=linux go build -o ./cmd/move-to-end/move-to-end ./cmd/move-to-end
    GOOS=linux go build -o ./cmd/process-execute/process-execute ./cmd/process-execute

## Deploy

    sam package --template-file ./template.yaml \
      --s3-bucket <S3_BUCKET> \ 
      --output-template-file ./template.packaged.yaml

    aws cloudformation deploy --template-file "$(PWD)/template.packaged.yaml" \
      --stack-name <STACK_NAME> \
      --capabilities CAPABILITY_NAMED_IAM


## Run
The workflow is started by `process-execute`, passing the input data to the
state machine.

```json
{
  "guid": "execution-id",
  "jobs": [{
    "guid": "execution-id-0",
    "input": "First job",
    "done": false
  }, {
    "guid": "execution-id-1",
    "input": "Second job",
    "done": false
  }]
}
```

Alternatively, start an execution using the AWS CLI and the correct state 
machine arn.

    aws stepfunctions start-execution --state-machine-arn <STATE_MACHINE_ARN> \
      --input "{\"jobs\": [{\"input\": \"First job\", \"done\": false}, {\"input\": \"Second job\", \"done\": false}]}"
