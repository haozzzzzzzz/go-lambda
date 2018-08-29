package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func generateProjTemplate(lambdaFunc *LambdaFunction) (err error) {
	innerProjDir := fmt.Sprintf("%s/.proj/", lambdaFunc.ProjectPath)
	err = os.MkdirAll(innerProjDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project proj folder failed. \n%s.", err)
		return
	}

	// create secret folder
	secretDir := fmt.Sprintf("%s/secret", innerProjDir)
	err = os.MkdirAll(secretDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project secret folder failed. \n%s.", err)
		return
	}

	// create .gitignore
	err = ioutil.WriteFile(fmt.Sprintf("%s/.gitignore", innerProjDir), []byte("secret"), project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("add project .gitignore file failed. \n%s.", err)
		return
	}

	// create project yaml
	projYamlConfig := proj.ProjectYamlFile{
		Name:            lambdaFunc.Name,
		Description:     lambdaFunc.Description,
		ProjectPath:     lambdaFunc.ProjectPath,
		EventSourceType: lambdaFunc.EventSourceType,
	}
	err = projYamlConfig.Save()
	if nil != err {
		logrus.Errorf("save project config failed. \n%s.", err)
		return
	}

	// check aws file
	_, _, err = proj.CheckAWSYamlFile(lambdaFunc.ProjectPath, true)
	if nil != err {
		logrus.Errorf("check aws yaml file failed. \n%s.", err)
		return
	}

	return
}
