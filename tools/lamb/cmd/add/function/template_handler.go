package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func generateHandlerTemplate(lambdaFunc *LambdaFunction) (err error) {
	handlerDir := fmt.Sprintf("%s/handler", lambdaFunc.ProjectPath)
	err = os.MkdirAll(handlerDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project handler directory failed. \n%s.", err)
		return
	}

	basicExecutionFileName := fmt.Sprintf("%s/handler_basicexecution.go", handlerDir)
	err = ioutil.WriteFile(basicExecutionFileName, []byte(basicExecutionFileText), project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("write handler/handler_basicexecution.go failed. \n%s.", err)
		return
	}

	// init file
	strEventType := lambdaFunc.EventSourceType.String()
	newInitFileText := fmt.Sprintf(initFileText, strEventType)
	initFilePath := fmt.Sprintf("%s/init.go", handlerDir)
	err = ioutil.WriteFile(initFilePath, []byte(newInitFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write %q failed. \n%s.", initFilePath, err)
		return
	}

	return
}

var basicExecutionFileText = `package handler
import (
	"context"
	"fmt"
)

func BasicExecutionEventHandler(ctx context.Context, event interface{})(string, error) {
	return fmt.Sprintf("Hello, world."), nil
}
`

var initFileText = `package handler

var mainHandler = %sHandler

func GetMainHandler() interface{} {
	return mainHandler
}
`
