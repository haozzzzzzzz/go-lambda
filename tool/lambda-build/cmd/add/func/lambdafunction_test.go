package _func

import (
	"os"
	"testing"
)

func TestLambdaFunction_Run(t *testing.T) {
	lambdaFunc := LambdaFunction{
		Name: "LambdaHandler",
		Path: "/Users/hao/Documents/Projects/Github/go_lambda_learning/src/github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/add/func",
		Mode: os.ModePerm,
	}
	err := lambdaFunc.Run()
	if err != nil {
		t.Error(err)
	}
}
