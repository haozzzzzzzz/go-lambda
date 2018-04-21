package function

import (
	"fmt"
	"os"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/sirupsen/logrus"
)

func generateProjTemplate(lambdaFunc *LambdaFunction) (err error) {
	projDir := fmt.Sprintf("%s/.proj/", lambdaFunc.ProjectPath)
	err = os.MkdirAll(projDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project proj direcotry failed. \n%s.", err)
		return
	}

	// create project yaml
	projYamlConfig := proj.ProjectYamlFile{
		Name:        lambdaFunc.Name,
		ProjectPath: lambdaFunc.ProjectPath,
		Mode:        lambdaFunc.Mode,
	}
	err = projYamlConfig.Save()
	if nil != err {
		logrus.Errorf("save project config failed. \n%s.", err)
		return
	}

	return
}
