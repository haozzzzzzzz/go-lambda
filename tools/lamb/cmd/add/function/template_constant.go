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
		logrus.Errorf("make project constant directory failed. \n%s.", err)
		return
	}

	constantFileName := fmt.Sprintf("%s/constant.go", constantDir)
	newConstantFileText := fmt.Sprintf(constantFileText, lambdaFunc.Name)
	err = ioutil.WriteFile(constantFileName, []byte(newConstantFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write config/constant.go failed. \n%s.", err)
		return
	}

	return
}

var constantFileText = `package constant

const (
	LambdaFunctionName = "%s"
)
`
