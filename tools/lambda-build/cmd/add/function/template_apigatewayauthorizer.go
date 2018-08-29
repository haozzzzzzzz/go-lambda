package function

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func generateApiGatewayAuthorizer(lambdaFunc *LambdaFunction) (err error) {
	projectPath := lambdaFunc.ProjectPath
	// add handler
	authorizerFilePath := fmt.Sprintf("%s/handler/handler_apigatewayauthorizer.go", projectPath)
	err = ioutil.WriteFile(authorizerFilePath, []byte(authorizerFileText), lambdaFunc.Mode)
	if nil != err {
		logrus.Errorf("write %q file failed. %s", err)
		return
	}

	return
}

var authorizerFileText = `package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/haozzzzzzzz/go-lambda/resource/apigateway"
)

func ApiGatewayAuthorizerEventHandler(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (response *events.APIGatewayCustomAuthorizerResponse, err error) {
	response = apigateway.GetAllowAuthorizerResponse(constant.LambdaFunctionName, request.MethodArn)
	return
}
`
