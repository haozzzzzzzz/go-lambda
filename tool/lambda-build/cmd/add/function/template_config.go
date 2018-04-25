package function

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func generateConfigTemplate(lambdaFunc *LambdaFunction) (err error) {
	projectPath := lambdaFunc.ProjectPath
	mode := lambdaFunc.Mode

	// dev配置
	devConfigDir := fmt.Sprintf("%s/config_dev", projectPath)
	err = os.MkdirAll(devConfigDir, mode)
	if nil != err {
		logrus.Errorf("make project dev config directory failed. \n%s.", err)
		return
	}

	// test配置
	testConfigDir := fmt.Sprintf("%s/config_test", projectPath)
	err = os.MkdirAll(testConfigDir, mode)
	if nil != err {
		logrus.Errorf("make project test config directory failed. \n%s.", err)
		return
	}

	// prod配置
	configDir := fmt.Sprintf("%s/config_prod", projectPath)
	err = os.MkdirAll(configDir, mode)
	if nil != err {
		logrus.Errorf("make project prod config directory failed. \n%s.", err)
		return
	}

	return
}
