package _func

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func generateHandlerTemplate(lambdaFunc *LambdaFunction) (err error) {
	handlerDir := fmt.Sprintf("%s/handler", lambdaFunc.ProjectPath)
	err = os.MkdirAll(handlerDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project handler directory failed. \n%s.", err)
		return
	}

	basicExecutionFileName := fmt.Sprintf("%s/handler_basicexecution.go", handlerDir)
	err = ioutil.WriteFile(basicExecutionFileName, []byte(basicExecutionFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write handler/handler_basicexecution.go failed. \n%s.", err)
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
