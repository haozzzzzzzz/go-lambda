package function

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func generateMainTemplate(lambdaFunc *LambdaFunction) (err error) {
	projectPath := lambdaFunc.ProjectPath
	mode := lambdaFunc.Mode

	var handlerName string
	switch lambdaFunc.EventSourceType {
	case BasicExecutionEvent:
		handlerName = "BasicExecutionEventHandler"
	case CustomEvent:
		handlerName = "CustomEventHandler"
	case ApiGatewayEvent:
		handlerName = "ApiGatewayEventHandler"
	}

	newMainFileText := fmt.Sprintf(mainFileText, handlerName)
	mainGoFileName := fmt.Sprintf("%s/main.go", projectPath)
	err = ioutil.WriteFile(mainGoFileName, []byte(newMainFileText), mode)
	if nil != err {
		logrus.Errorf("write %q failed. \n%s.", mainGoFileName, err)
		return
	}
	return
}

var mainFileText = `
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

var mainHandler = handler.%s

func main() {
	lambda.Start(mainHandler)
}
`
