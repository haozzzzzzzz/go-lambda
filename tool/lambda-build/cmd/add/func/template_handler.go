package _func

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func generateHandlerTemplate(lambdaFunc *LambdaFunction) (err error) {
	handlerDir := fmt.Sprintf("%s/handler", lambdaFunc.ProjectPath)
	err = os.MkdirAll(handlerDir, lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("make project handler directory failed. \n%s.", err)
		return
	}

	apigatewayFileName := fmt.Sprintf("%s/apigateway.go", handlerDir)
	err = ioutil.WriteFile(apigatewayFileName, []byte(apigatewayFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write handler/apigateway.go failed. \n%s.", err)
		return
	}

	return
}

var apigatewayFileText = `package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy.git/gin"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/sirupsen/logrus"
)

// gin lambda adapter
var ginLambda *ginadapter.GinLambda

func NewGinLambda() (err error) {
	logrus.Infof("Lambda function %s initializing...", constant.LambdaFunctionName)
	ginEngine := ginbuilder.GetEngine()
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

func ApiGatewayHandler(ctx context.Context, request *events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
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
