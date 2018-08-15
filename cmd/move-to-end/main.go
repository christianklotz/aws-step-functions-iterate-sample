package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(list []interface{}) ([]interface{}, error) {
	if len(list) == 0 {
		return list, nil
	}
	return append(list[1:], list[0]), nil
}

func main() {
	lambda.Start(handler)
}
