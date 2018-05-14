package function

import (
	"fmt"
	"io/ioutil"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/sirupsen/logrus"
)

func generateMainTemplate(lambdaFunc *LambdaFunction) (err error) {
	projectPath := lambdaFunc.ProjectPath
	mode := lambdaFunc.Mode

	// main.go
	newMainFileText := mainFileText
	mainGoFileName := fmt.Sprintf("%s/main.go", projectPath)
	err = ioutil.WriteFile(mainGoFileName, []byte(newMainFileText), mode)
	if nil != err {
		logrus.Errorf("write %q failed. \n%s.", mainGoFileName, err)
		return
	}

	lambdaFuncName := lambdaFunc.Name

	// create .gitignore
	gitIgnoreFileText := `detector
main
%sprod
%spre
%stest
%sdev
`
	gitIgnoreFileText = fmt.Sprintf(gitIgnoreFileText, lambdaFuncName, lambdaFuncName, lambdaFuncName, lambdaFuncName)
	gitIgnoreFilePath := fmt.Sprintf("%s/.gitignore", projectPath)
	err = ioutil.WriteFile(gitIgnoreFilePath, []byte(gitIgnoreFileText), mode)
	if nil != err {
		logrus.Errorf("write %q failed. %s.", gitIgnoreFilePath, err)
		return
	}

	err = createDeployShellFile(lambdaFunc)
	if nil != err {
		logrus.Errorf("create deploy shell file failed. \n%s.", err)
		return
	}

	return
}

var mainFileText = `
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.GetMainHandler())
}
`

func createDeployShellFile(lambdaFunc *LambdaFunction) (err error) {
	projectPath := lambdaFunc.ProjectPath
	mode := lambdaFunc.Mode
	var deployShellFileText string
	var runShellFileText string

	switch lambdaFunc.EventSourceType {
	case proj.ApiGatewayProxyEvent:
		runShellFileText = `#!/usr/bin/env bash
echo generating api
lbuild compile api

echo running main
go build -o main main.go
./main`

		deployShellFileText = `#!/usr/bin/env bash

stage=$1

if [ -z ${stage} ]
then
    stage="dev"
fi

echo deploy target stage environment: ${stage}

echo generating api
lbuild compile api

echo building...
lbuild compile func -s ${stage}

echo deploying...
ldeploy remote func -s ${stage}
`

	default:
		runShellFileText = `#!/usr/bin/env bash
echo running main
go build -o main main.go
./main`

		deployShellFileText = `#!/usr/bin/env bash

stage=$1

if [ -z ${stage} ]
then
    stage="dev"
fi

echo deploy target stage environment: ${stage}

echo building...
lbuild compile func -s ${stage}

echo deploying...
ldeploy remote func -s ${stage}
`

	}

	// run.sh
	runShellFilePath := fmt.Sprintf("%s/run.sh", projectPath)
	err = ioutil.WriteFile(runShellFilePath, []byte(runShellFileText), mode)
	if nil != err {
		logrus.Errorf("write %q failed. \n%s.", runShellFilePath, err)
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
