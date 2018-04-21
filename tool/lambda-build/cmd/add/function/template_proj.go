package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/sirupsen/logrus"
)

func generateProjTemplate(lambdaFunc *LambdaFunction) (err error) {
	projDir := fmt.Sprintf("%s/.proj/", lambdaFunc.ProjectPath)
	err = os.MkdirAll(projDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project proj folder failed. \n%s.", err)
		return
	}

	// create secret folder
	err = os.MkdirAll(fmt.Sprintf("%s/secret", projDir), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project secret folder failed. \n%s.", err)
		return
	}

	// create .gitignore
	err = ioutil.WriteFile(fmt.Sprintf("%s/.gitignore", projDir), []byte("secret"), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("add project .gitignore file failed. \n%s.", err)
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
