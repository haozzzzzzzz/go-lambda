package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api"
	"github.com/sirupsen/logrus"
)

func generateApiTemplate(lambdaFunc *LambdaFunction) (err error) {
	// add handler
	apigatewayFileName := fmt.Sprintf("%s/handler/handler_apigateway.go", lambdaFunc.ProjectPath)
	err = ioutil.WriteFile(apigatewayFileName, []byte(apigatewayFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write handler/apigateway.go failed. \n%s.", err)
		return
	}

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

	// 建立本地api测试main文件
	err = generateApiLocalTemplate(lambdaFunc)
	if nil != err {
		logrus.Errorf("generate api local template failed. %s.", err)
		return
	}

	return
}

var apigatewayFileText = `package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy.git/gin"
	"github.com/haozzzzzzzz/go-lambda/resource/xray"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/sirupsen/logrus"
)

// gin lambda adapter
var ginLambda *ginadapter.GinLambda

func NewGinLambda() (err error) {
	logrus.Infof("Lambda function %s initializing...", constant.LambdaFunctionName)
	ginEngine := ginbuilder.GetEngine()

	// bind xray
	ginEngine.Use(xray.XRayGinMiddleware(constant.LambdaFunctionName))

	// bind routers
	err = api.BindRouters(ginEngine)
	if nil != err {
		logrus.Errorf("set http router failed. \n%s.", err)
		return
	}

	ginLambda = ginadapter.New(ginEngine)
	return
}

func init() {
	err := NewGinLambda()
	if nil != err {
		logrus.Errorf("new gin lambda failed. %s", err)
		return
	}
}

func ApiGatewayProxyEventHandler(ctx context.Context, request *events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	if nil == ginLambda {
		err = NewGinLambda()
		if nil != err {
			logrus.Errorf("new gin lambda failed. \n%s.", err)
			return
		}
	}

	return ginLambda.Proxy(*request)
}
`

var routersFileText = `package api

import (
	"github.com/gin-gonic/gin"
)

// 注意：BindRouters函数体内不能自定义添加任何声明，由lambda-build compile api命令生成api绑定声明
func BindRouters(engine *gin.Engine) (err error) {
	return
}
`
