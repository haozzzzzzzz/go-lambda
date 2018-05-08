package proj

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (m *SAMTemplateConfig) BuildApiGatewayProxyEvent() (err error) {
	lambdaFunctionName := m.LambdaFunctionName
	stage := m.State
	templateFile := m.SAMTemplateYamlFile

	// lambda function
	resourceLambdaFunction, ok := templateFile.Resources[lambdaFunctionName].(*SAMResource)
	if !ok {
		err = errors.New("lambda function type should be *SAMResource")
		return
	}

	proxyPath := map[string]interface{}{
		"x-amazon-apigateway-any-method": map[string]interface{}{
			"produces": []string{
				"application/json",
			},
			"x-amazon-apigateway-integration": map[string]interface{}{
				"type":                "aws_proxy",
				"httpMethod":          "POST",
				"passthroughBehavior": "when_no_match",
				"uri": map[string]interface{}{
					"Fn::Sub": fmt.Sprintf("arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${%s.Arn}/invocations", lambdaFunctionName),
				},
			},
		},
	}

	apiResourceName := "ApiGatewayApi"
	apiDefinitionBody := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"version": "1.0",
			"title":   lambdaFunctionName,
		},
		"basePath": fmt.Sprintf("/%s", stage),
		"schemes":  []string{"https"},
		"paths": map[string]interface{}{
			"/{proxy+}": proxyPath,
		},
	}

	apiResource := &SAMResource{
		Type: "AWS::Serverless::Api",
		Properties: map[string]interface{}{
			"Name":           lambdaFunctionName, // 显示在ApiGateway控制台的资源名称
			"StageName":      stage,
			"DefinitionBody": apiDefinitionBody,
		},
	}

	// lambda function api event
	lambdaFunctionEvents := make(map[string]interface{})
	resourceLambdaFunction.Properties["Events"] = lambdaFunctionEvents
	apiEventName := "ApiEvent"
	lambdaFunctionEvents[apiEventName] = map[string]interface{}{
		"Type": "Api",
		"Properties": map[string]interface{}{
			"RestApiId": map[string]interface{}{
				"Ref": apiResourceName,
			},
			"Path":   "/{proxy+}",
			"Method": "ANY",
		},
	}

	templateFile.Resources[apiResourceName] = apiResource

	// permission
	apiAccessPermissionName := "ApiAccessPermission"
	apiPermissionResource := &SAMResource{
		Type: "AWS::Lambda::Permission",
		Properties: map[string]interface{}{
			"Action": "lambda:InvokeFunction",
			"FunctionName": map[string]interface{}{
				"Ref": lambdaFunctionName,
			},
			"Principal": "apigateway.amazonaws.com",
			"SourceArn": map[string]interface{}{
				"Fn::Sub": fmt.Sprintf("arn:aws:execute-api:${AWS::Region}:${AWS::AccountId}:${%s}/%s/*/*", apiResourceName, stage),
			},
		},
	}

	templateFile.Resources[apiAccessPermissionName] = apiPermissionResource

	// authorizer
	authorizerConfig, err := NewAuthorizerConfig(m.ProjectYamlFile)
	if nil != err {
		logrus.Errorf("new authorizer config failed. %s.", err)
		return
	}

	if authorizerConfig.AuthorizerYamlFileExsit {
		sercurity := make([]interface{}, 0)
		sercurityDefinitions := make(map[string]interface{})
		for _, authorizer := range authorizerConfig.AuthorizerYamlFile.Authorizers {
			sercurity = append(sercurity, map[string][]string{
				authorizer.Name: make([]string, 0),
			})

			var identitySources []string
			for _, header := range authorizer.Headers {
				identitySources = append(identitySources, fmt.Sprintf("method.request.header.%s", header))
			}

			for _, query := range authorizer.Queries {
				identitySources = append(identitySources, fmt.Sprintf("method.request.querystring.%s", query))
			}

			identitySource := strings.Join(identitySources, ",")
			sercurityDefinitions[authorizer.Name] = map[string]interface{}{
				"type": "apiKey",
				"name": "Unused",
				"in":   "header",
				"x-amazon-apigateway-authtype": "custom",
				"x-amazon-apigateway-authorizer": map[string]interface{}{
					"authorizerUri": map[string]interface{}{
						"Fn::Sub": fmt.Sprintf("arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:%s/invocations", authorizer.FunctionHandler),
					},
					"authorizerResultTtlInSeconds": 300,
					"identitySource":               identitySource,
					"type":                         "request",
				},
			}

			// authorizer access permission
			authorizerAccessPermissionName := fmt.Sprintf("Authorizer%sAccessPermission", authorizer.Name)
			accessPermissionResource := &SAMResource{
				Type: "AWS::Lambda::Permission",
				Properties: map[string]interface{}{
					"Action":       "lambda:InvokeFunction",
					"FunctionName": authorizer.FunctionHandler,
					"Principal":    "apigateway.amazonaws.com",
					"SourceArn": map[string]interface{}{
						"Fn::Sub": fmt.Sprintf("arn:aws:execute-api:${AWS::Region}:${AWS::AccountId}:${%s}/authorizers/*", apiResourceName),
					},
				},
			}

			templateFile.Resources[authorizerAccessPermissionName] = accessPermissionResource
		}

		proxyPath["x-amazon-apigateway-any-method"].(map[string]interface{})["security"] = sercurity
		apiDefinitionBody["securityDefinitions"] = sercurityDefinitions

	}

	// 输出
	templateFile.Outputs["ApiUrl"] = map[string]interface{}{
		"Description": fmt.Sprintf("%s(%s) Api URL", lambdaFunctionName, stage),
		"Value": map[string]interface{}{
			"Fn::Sub": fmt.Sprintf("https://${%s}.execute-api.${AWS::Region}.amazonaws.com/%s", apiResourceName, stage),
		},
	}

	return
}
