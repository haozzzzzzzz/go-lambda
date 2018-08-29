package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func generateApiLocalTemplate(lambdaFunc *LambdaFunction) (err error) {
	projPath := lambdaFunc.ProjectPath

	// local
	localDir := fmt.Sprintf("%s/local/api", projPath)
	err = os.MkdirAll(localDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make directory %q failed. %s.", localDir, err)
		return
	}

	localMainFilePath := fmt.Sprintf("%s/main.go", localDir)
	err = ioutil.WriteFile(localMainFilePath, []byte(localMainFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write %q failed. %s.", localMainFilePath, err)
		return
	}

	return
}

var localMainFileText = `package main

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

func main()  {
	engine := ginbuilder.GetEngine()
	api.BindRouters(engine)
	engine.Run(":8100")
}
`
