package _func

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func generateProjTemplate(lambdaFunc *LambdaFunction) (err error) {
	projDir := fmt.Sprintf("%s/.proj/", lambdaFunc.ProjectPath)
	err = os.MkdirAll(projDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project proj direcotry failed. \n%s.", err)
		return
	}

	// create project yaml
	projYamlFileName := fmt.Sprintf("%s/proj.yaml", projDir)
	projYamlConfig := proj.ProjectYamlFile{
		Name:        lambdaFunc.Name,
		ProjectPath: lambdaFunc.ProjectPath,
	}
	byteProjYamlFile, err := yaml.Marshal(projYamlConfig)
	if nil != err {
		logrus.Errorf("marshal proj yaml config file failed. \n%s.", err)
		return
	}
	err = ioutil.WriteFile(projYamlFileName, byteProjYamlFile, lambdaFunc.Mode)
	if nil != err {
		logrus.Warnf("write .proj/proj.yaml failed. \n%s.", err)
		return
	}

	return
}
