package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func generateConstantTemplate(lambdaFunc *LambdaFunction) (err error) {
	constantDir := fmt.Sprintf("%s/constant", lambdaFunc.ProjectPath)
	err = os.MkdirAll(constantDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project constent directory failed. \n%s.", err)
		return
	}

	lambdaFileName := fmt.Sprintf("%s/lambda.go", constantDir)
	newLambdaFileText := fmt.Sprintf(lambdaFileText, lambdaFunc.Name)
	err = ioutil.WriteFile(lambdaFileName, []byte(newLambdaFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write constant/lambda.go failed. \n%s.", err)
		return
	}

	return
}

var lambdaFileText = `package constant

const (
	LambdaFunctionName = "%s"
)
`
