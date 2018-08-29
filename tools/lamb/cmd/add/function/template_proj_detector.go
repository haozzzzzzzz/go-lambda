package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func generateDetectorMainTemplate(lambdaFunc *LambdaFunction) (err error) {
	detectorDir := fmt.Sprintf("%s/.proj/detector", lambdaFunc.ProjectPath)
	err = os.MkdirAll(detectorDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make dir %q failed. \n%s.", detectorDir, err)
		return
	}

	detectorFilePath := fmt.Sprintf("%s/main.go", detectorDir)
	err = ioutil.WriteFile(detectorFilePath, []byte(detectorMainGoFileText), project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write %q failed. \n%s.", detectorFilePath, err)
		return
	}

	return
}

var detectorMainGoFileText = `package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/haozzzzzzzz/go-lambda/resource/iam"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error
	filePath := flag.String("path", "", "role.yaml save file path")
	flag.Parse()

	if *filePath == "" {
		logrus.Errorf("wrong role save target file path")
		return
	}

	handler.GetMainHandler()

	roleFilePath, err := filepath.Abs(*filePath)
	if nil != err {
		logrus.Errorf("get absolute file path failed. \n%s.", err)
		return
	}

	logrus.Info("detecting project")
	defer func() {
		logrus.Info("detecting project finish")
	}()

	role := iam.NewExecutionRole(fmt.Sprintf("%sRole", constant.LambdaFunctionName))
	err = role.WriteTo(roleFilePath, os.ModePerm)
	if nil != err {
		logrus.Errorf("write role to yaml failed. \n%s.", err)
		return
	}
}
`
