# AWS Step Functions Iteration Sample
A complete serverless stack demonstrating how to process an arbitrary number of
jobs using AWS Step Functions. This is the sample project to the article
[Processing an arbitrary number of jobs with AWS Step Functions](https://medium.com/@christianklotz/c185c2d2608).

## Build

    GOOS=linux go build -o ./cmd/move-to-end/move-to-end ./cmd/move-to-end
    GOOS=linux go build -o ./cmd/move-to-end/process-execute ./cmd/process-execute

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
