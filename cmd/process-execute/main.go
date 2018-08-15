package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/sfn/sfniface"
	"github.com/gofrs/uuid"
)

// The payload to invoke the Lambda function with, in different scenarios this
// could be an S3, CloudWatch or different event.
type event struct {
	// Input is expected to be a comma-separated list of strings, resulting in
	// a corresponding number of jobs the workflow is executed with.
	// Example: foo,bar will create two jobs, one for "foo", another for "bar".
	Input string `json:"input"`
}

type job struct {
	GUID  string `json:"guid"`
	Input string `json:"input"`

	Done bool `json:"done"`
}

type lambdaFunction struct {
	runner       sfniface.SFNAPI
	stateMachine string
}

type eventHandler func(evt event) error

func (fn *lambdaFunction) handler() eventHandler {
	return func(evt event) error {
		if evt.Input == "" {
			return errors.New("no jobs to process")
		}

		id, err := uuid.NewV4()
		if err != nil {
			return err
		}

		var input = struct {
			GUID string `json:"guid"`
			Jobs []job  `json:"jobs"`
		}{
			GUID: id.String(),
		}

		vals := strings.Split(evt.Input, ",")
		for i, v := range vals {
			input.Jobs = append(input.Jobs, job{
				GUID:  fmt.Sprintf("%s-%d", input.GUID, i),
				Input: v,
				Done:  false,
			})
		}

		var data []byte
		if data, err = json.Marshal(input); err != nil {
			return err
		}

		params := &sfn.StartExecutionInput{
			StateMachineArn: aws.String(fn.stateMachine),
			Name:            aws.String(id.String()),
			Input:           aws.String(string(data)),
		}
		_, err = fn.runner.StartExecution(params)
		return err
	}
}

// Starts the function, waiting for invocation.
func (fn *lambdaFunction) Start() {
	lambda.Start(fn.handler())
}

func main() {
	var (
		arn = flag.String("state-machine-arn", os.Getenv("STATE_MACHINE_ARN"),
			"The arn string of the state machine that should be started (default $STATE_MACHINE_ARN).")
	)

	flag.Parse()
	if *arn == "" {
		fmt.Println("missing state machine arn")
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("failed to init session:", err)
	}

	fn := lambdaFunction{
		runner:       sfn.New(sess),
		stateMachine: *arn,
	}
	fn.Start()
}
