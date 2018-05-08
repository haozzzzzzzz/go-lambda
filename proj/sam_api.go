package proj

import (
	"fmt"

	"github.com/pkg/errors"
)

func (m *SAMTemplateConfig) BuildApiGatewayProxyEvent() (err error) {
	lambdaFunctionName := m.LambdaFunctionName
	stage := m.State
	templateFile := m.SAMTemplateYamlFile
	resourceLambdaFunction, ok := templateFile.Resources[lambdaFunctionName].(*SAMResource)
	if !ok {
		err = errors.New("lambda function type should be *SAMResource")
		return
	}

	apiResourceName := "ApiGatewayApi"
	apiResource := &SAMResource{
		Type: "AWS::Serverless::Api",
		Properties: map[string]interface{}{
			"Name":      lambdaFunctionName, // 显示在ApiGateway控制台的资源名称
			"StageName": stage,
			"DefinitionBody": map[string]interface{}{
				"swagger": "2.0",
				"info": map[string]interface{}{
					"version": "1.0",
					"title":   lambdaFunctionName,
				},
				"basePath": fmt.Sprintf("/%s", stage),
				"schemes":  []string{"https"},
				"paths": map[string]interface{}{
					"/{proxy+}": map[string]interface{}{
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
					},
				},
			},
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

	// 输出
	templateFile.Outputs["ApiUrl"] = map[string]interface{}{
		"Description": fmt.Sprintf("%s(%s) Api URL", lambdaFunctionName, stage),
		"Value": map[string]interface{}{
			"Fn::Sub": fmt.Sprintf("https://${%s}.execute-api.${AWS::Region}.amazonaws.com/%s", apiResourceName, stage),
		},
	}

	return
}
