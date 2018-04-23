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

	// main.go
	newMainFileText := fmt.Sprintf(mainFileText, handlerName)
	mainGoFileName := fmt.Sprintf("%s/main.go", projectPath)
	err = ioutil.WriteFile(mainGoFileName, []byte(newMainFileText), mode)
	if nil != err {
		logrus.Errorf("write %q failed. \n%s.", mainGoFileName, err)
		return
	}

	// deploy.sh
	deployShellFilePath := fmt.Sprintf("%s/deploy.sh", projectPath)
	err = ioutil.WriteFile(deployShellFilePath, []byte(deployShellFileText), mode)
	if nil != err {
		logrus.Errorf("write %q failed. \n%s.", deployShellFilePath, err)
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

var deployShellFileText = `#!/usr/bin/env bash
echo building...
lamb compile func

echo deploying...
lamd remote func`
