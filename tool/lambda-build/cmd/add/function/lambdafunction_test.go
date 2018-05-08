package function

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func TestLambdaFunction_Run(t *testing.T) {
	lambdaFunc := LambdaFunction{
		Name: "LambdaHandler",
		Path: "/Users/hao/Documents/Projects/Github/go_lambda_learning/src/github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/add/function",
		Mode: os.ModePerm,
	}
	err := lambdaFunc.Run()
	if err != nil {
		t.Error(err)
	}
}

func TestNameLimit(t *testing.T) {
	matched, err := regexp.MatchString("^[A-za-z][A-Za-z0-9]+$", "1Helsss111loWorld")
	if nil != err {
		return
	}
	fmt.Println(matched)
}
