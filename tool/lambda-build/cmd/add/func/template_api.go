package _func

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api"
	"github.com/sirupsen/logrus"
)

func generateApiTemplate(lambdaFunc *LambdaFunction) (err error) {
	// api
	apiDir := fmt.Sprintf("%s/api", lambdaFunc.ProjectPath)
	err = os.MkdirAll(apiDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project api directory failed. \n%s.", err)
		return
	}

	routersFileName := fmt.Sprintf("%s/routers.go", apiDir)
	err = ioutil.WriteFile(routersFileName, []byte(routersFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write api/routers.go failed. \n%s.", err)
		return
	}

	// 建立一个api示例
	metricDir := fmt.Sprintf("%s/metric", apiDir)
	err = os.MkdirAll(metricDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make api example \"metric\" directory failed. \n%s.", err)
		return
	}

	err = api.CreateApiSource(&api.ApiItem{
		HttpMethod:        "GET",
		RelativePath:      "/metric",
		ApiHandlerFunc:    "MetricHandlerFunc",
		ApiHandlerPackage: "metric",
		SourceFile:        fmt.Sprintf("%s/api_metric.go", metricDir),
	})
	if nil != err {
		logrus.Errorf("create api metric/api_info.go file failed. \n%s.", err)
		return
	}

	return
}

var routersFileText = `package api

import (
	"github.com/gin-gonic/gin"
)

func BindRouters(engine *gin.Engine) (err error) {
	return
}
`
