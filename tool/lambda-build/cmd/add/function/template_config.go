package function

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func generateConfigTemplate(lambdaFunc *LambdaFunction) (err error) {
	configDir := fmt.Sprintf("%s/config", lambdaFunc.ProjectPath)
	err = os.MkdirAll(configDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project config directory failed. \n%s.", err)
		return
	}

	return
}
