package function

import (
	"testing"
)

func TestCompileFunction_Run(t *testing.T) {
	compileFunc := CompileFunction{
		ProjectPath: "/Users/hao/Documents/Projects/Github/go_lambda_learning/src/ExampleLambdaApi",
		Stage:       "test",
	}
	err := compileFunc.Run()
	if nil != err {
		t.Error(err)
		return
	}
}
