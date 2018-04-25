package proj

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-lambda/resource/iam"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type SAMTemplateYamlFile struct {
	AWSTemplateFormatVersion string                 `yaml:"AWSTemplateFormatVersion"`
	Transform                string                 `yaml:"Transform"`
	Description              string                 `yaml:"Description"`
	Resources                map[string]interface{} `yaml:"Resources"`
	Outputs                  map[string]interface{} `yaml:"Outputs"`
}

type SAMResource struct {
	Type       string                 `yaml:"Type"`
	Properties map[string]interface{} `yaml:"Properties"`
}

func NewSAMTemplateYamlFileByExistConfig(stage string, projConfig *ProjectYamlFile, awsConfig *AWSYamlFile) (templateFile *SAMTemplateYamlFile, err error) {
	projectPath := projConfig.ProjectPath
	templateFile = &SAMTemplateYamlFile{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Description:              projConfig.Description,
		Resources:                make(map[string]interface{}),
		Outputs:                  make(map[string]interface{}),
	}

	lambdaFunctionName := projConfig.Name

	// 角色
	var funcRole interface{}
	if awsConfig.Role == "" {
		roleYamlFilePath := fmt.Sprintf("%s/.proj/role.yaml", projectPath)
		role, errLoad := iam.LoadRoleFromFile(roleYamlFilePath)
		if nil != errLoad {
			err = errLoad
			logrus.Errorf("load role.yaml from file failed. \n%s.", err)
			return
		}

		roleName := role.Properties.RoleName
		templateFile.Resources[roleName] = role

		funcRole = map[string]interface{}{
			"Fn::GetAtt": []string{
				roleName, "Arn",
			},
		}

	} else {
		funcRole = fmt.Sprintf("arn:aws:iam::%s:role/%s", awsConfig.AccountId, awsConfig.Role)

	}

	// 发布流量转移类型
	var deploymentType string
	switch stage {
	case TestStage.String():
		deploymentType = "AllAtOnce" // 立刻转移
	case ProdStage.String():
		deploymentType = "Canary10Percent10Minutes" // 10分钟完成转移
	}

	// lambda函数
	resourceLambdaFunction := SAMResource{
		Type: "AWS::Serverless::Function",
		Properties: map[string]interface{}{
			"Handler":          projConfig.Name,
			"FunctionName":     projConfig.Name,
			"Runtime":          "go1.x",
			"CodeUri":          fmt.Sprintf("./%s.zip", projConfig.Name),
			"Description":      projConfig.Description,
			"Role":             funcRole,
			"AutoPublishAlias": stage,
			"DeploymentPreference": map[string]interface{}{
				"Type": deploymentType,
				//"Alarms": []interface{}{ // A list of alarms that you want to monitor
				//	map[string]interface{}{
				//		"Ref": "AliasErrorMetricGreaterThanZeroAlarm",
				//	},
				//	map[string]interface{}{
				//		"Ref": "LatestVersionErrorMetricGreaterThanZeroAlarm",
				//	},
				//},
				//"Hooks": map[string]interface{}{ //Validation Lambda functions that are run before & after traffic shifting
				//	"PreTraffic": map[string]interface{}{
				//		"Ref": lambdaFunctionName,
				//	},
				//	"PostTraffic": map[string]interface{}{
				//		"Ref": lambdaFunctionName,
				//	},
				//},
			},
		},
	}

	templateFile.Resources[lambdaFunctionName] = resourceLambdaFunction

	// api gateway event
	switch projConfig.EventSourceType {
	case ApiGatewayEvent:
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
	}

	return
}

func (m *SAMTemplateYamlFile) Save(stage string, projectPath string, mode os.FileMode) (err error) {
	samYamlFilePath := fmt.Sprintf("%s/deploy/%s/template.yaml", projectPath, stage)
	byteYaml, err := yaml.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal sam yaml file failed. \n%s.", err)
		return
	}

	err = ioutil.WriteFile(samYamlFilePath, byteYaml, mode)
	if nil != err {
		logrus.Errorf("write sam yaml file failed. \n%s.", err)
		return
	}

	return
}
