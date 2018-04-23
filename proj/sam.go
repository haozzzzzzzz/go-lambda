package proj

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type SAMTemplateYamlFile struct {
	AWSTemplateFormatVersion string                 `yaml:"AWSTemplateFormatVersion"`
	Transform                string                 `yaml:"Transform"`
	Description              string                 `yaml:"Description"`
	Resources                map[string]interface{} `yaml:"Resources"`
}

type SAMResource struct {
	Type       string                 `yaml:"Type"`
	Properties map[string]interface{} `yaml:"Properties"`
}

func NewSAMTemplateYamlFile(projectPath string) (templateFile *SAMTemplateYamlFile, err error) {
	projectConfig, err := LoadProjectYamlFile(projectPath)
	if nil != err {
		logrus.Errorf("load project yaml config file failed. \n%s.", err)
		return
	}

	awsConfig, err := LoadAWSYamlFile(projectPath)
	if nil != err {
		logrus.Errorf("load aws yaml config file failed. \n%s.", err)
		return
	}

	templateFile = NewSAMTemplateYamlFileByExistConfig(projectConfig, awsConfig)
	return
}

func NewSAMTemplateYamlFileByExistConfig(projConfig *ProjectYamlFile, awsConfig *AWSYamlFile) (templateFile *SAMTemplateYamlFile) {
	templateFile = &SAMTemplateYamlFile{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Description:              projConfig.Description,
		Resources:                make(map[string]interface{}),
	}

	resourceLambdaFunction := SAMResource{
		Type: "AWS::Serverless::Function",
		Properties: map[string]interface{}{
			"Handler":          projConfig.Name,
			"FunctionName":     projConfig.Name,
			"Runtime":          "go1.x",
			"CodeUri":          fmt.Sprintf("./%s.zip", projConfig.Name),
			"Description":      projConfig.Description,
			"Role":             fmt.Sprintf("arn:aws:iam::%s:role/%s", awsConfig.AccountId, awsConfig.Role),
			"AutoPublishAlias": "live",
			"DeploymentPreference": map[string]interface{}{
				"Type": "Canary10Percent10Minutes",
				"Alarms": []string{ // A list of alarms that you want to monitor
					"!Ref AliasErrorMetricGreaterThanZeroAlarm",
					"!Ref LatestVersionErrorMetricGreaterThanZeroAlarm",
				},
				"Hooks": map[string]interface{}{ //Validation Lambda functions that are run before & after traffic shifting
					"PreTraffic":  "!Ref PreTrafficLambdaFunction",
					"PostTraffic": "!Ref PostTrafficLambdaFunction",
				},
			},
		},
	}

	switch projConfig.EventSourceType {
	case ApiGatewayEvent:
		lambdaFunctionEvents := make(map[string]interface{})
		resourceLambdaFunction.Properties["Events"] = lambdaFunctionEvents
		lambdaFunctionEvents[projConfig.Name] = map[string]interface{}{
			"Type": "Api",
			"Properties": map[string]interface{}{
				"Path":   "/{proxy+}",
				"Method": "any",
			},
		}

	}

	templateFile.Resources[projConfig.Name] = resourceLambdaFunction

	return
}

func (m *SAMTemplateYamlFile) Save(projectPath string, mode os.FileMode) (err error) {
	samYamlFilePath := fmt.Sprintf("%s/deploy/template.yaml", projectPath)
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
