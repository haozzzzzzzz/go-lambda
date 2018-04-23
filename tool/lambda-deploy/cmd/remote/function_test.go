package remote

import (
	"os"
	"testing"
)

func TestRemoteLambdaFunction_Run(t *testing.T) {
	remoteLambda := RemoteLambdaFunction{
		ProjectPath: "/Users/hao/Documents/Projects/Github/go_lambda_learning/src/ExampleLambdaBasic",
		Mode:        os.ModePerm,
	}

	err := remoteLambda.Run()
	if nil != err {
		t.Error(err)
		return
	}
}
