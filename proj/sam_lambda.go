package proj

import (
	"fmt"

	"github.com/haozzzzzzzz/go-lambda/resource/iam"
	"github.com/sirupsen/logrus"
)

func (m *SAMTemplateConfig) BuildLambdaFunction() (err error) {
	stage := m.State
	projConfig := m.ProjectYamlFile
	awsConfig := m.AWSYamlFile
	projectPath := m.ProjectPath
	templateFile := m.SAMTemplateYamlFile

	lambdaFunctionName := m.LambdaFunctionName

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

		// 更改角色名称
		role.Properties.RoleName = fmt.Sprintf("%sRole", lambdaFunctionName)
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
		//deploymentType = "Canary10Percent10Minutes" // 10分钟完成转移
		deploymentType = "AllAtOnce" // 立即转移
	}

	// lambda函数
	resourceLambdaFunction := &SAMResource{
		Type: "AWS::Serverless::Function",
		Properties: map[string]interface{}{
			"Handler":          lambdaFunctionName,
			"FunctionName":     lambdaFunctionName,
			"Runtime":          "go1.x",
			"CodeUri":          fmt.Sprintf("./%s.zip", lambdaFunctionName),
			"Description":      projConfig.Description,
			"Role":             funcRole,
			"AutoPublishAlias": stage,
			"Timeout":          30,
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

	return
}
