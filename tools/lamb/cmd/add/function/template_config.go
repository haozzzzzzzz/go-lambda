package function

import (
	"fmt"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func generateConfigTemplate(lambdaFunc *LambdaFunction) (err error) {
	projectPath := lambdaFunc.ProjectPath

	// dev配置
	devConfigDir := fmt.Sprintf("%s/stage/dev/config", projectPath)
	err = os.MkdirAll(devConfigDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project dev config directory failed. \n%s.", err)
		return
	}

	// test配置
	testConfigDir := fmt.Sprintf("%s/stage/test/config", projectPath)
	err = os.MkdirAll(testConfigDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project test config directory failed. \n%s.", err)
		return
	}

	// pre配置
	preConfigDir := fmt.Sprintf("%s/stage/pre/config", projectPath)
	err = os.MkdirAll(preConfigDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project pre config directory failed. %s.", err)
		return
	}

	// prod配置
	configDir := fmt.Sprintf("%s/stage/prod/config", projectPath)
	err = os.MkdirAll(configDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make project prod config directory failed. \n%s.", err)
		return
	}

	return
}
