package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/add/func/LambdaHandler/handler"
)

func main() {
	lambda.Start(handler.ApiGatewayHandler)
}
