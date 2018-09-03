package proj

import (
	"fmt"
	"io/ioutil"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/utils/str"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type SAMTemplateConfig struct {
	State               string
	ProjectPath         string
	ProjectYamlFile     *ProjectYamlFile
	AWSYamlFile         *AWSYamlFile
	SAMTemplateYamlFile *SAMTemplateYamlFile
	LambdaFunctionName  string
}

func NewSAMTempalteConfig(stage string, projConfig *ProjectYamlFile, awsConfig *AWSYamlFile) (config *SAMTemplateConfig, err error) {
	config = &SAMTemplateConfig{
		State:               stage,
		ProjectPath:         projConfig.ProjectPath,
		ProjectYamlFile:     projConfig,
		AWSYamlFile:         awsConfig,
		SAMTemplateYamlFile: NewSAMTemplateYamlFile(),
		LambdaFunctionName:  fmt.Sprintf("%s%s", str.CamelString(stage), projConfig.Name),
	}

	err = config.Build()
	if nil != err {
		logrus.Errorf("build sam template failed. %s.", err)
		return
	}
	return
}

type SAMTemplateYamlFile struct {
	AWSTemplateFormatVersion string                 `yaml:"AWSTemplateFormatVersion"`
	Transform                string                 `yaml:"Transform"`
	Description              string                 `yaml:"Description"`
	Resources                map[string]interface{} `yaml:"Resources"`
	Outputs                  map[string]interface{} `yaml:"Outputs"`
}

func NewSAMTemplateYamlFile() (templateFile *SAMTemplateYamlFile) {
	templateFile = &SAMTemplateYamlFile{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Description:              "AWS Serverless Lambda Function",
		Resources:                make(map[string]interface{}),
		Outputs:                  make(map[string]interface{}),
	}
	return
}

type SAMResource struct {
	Type       string                 `yaml:"Type"`
	Properties map[string]interface{} `yaml:"Properties"`
}

func (m *SAMTemplateConfig) Build() (err error) {
	templateFile := m.SAMTemplateYamlFile
	projConfig := m.ProjectYamlFile

	// 对象初始化
	templateFile.Description = projConfig.Description

	// 创建lambda函数相关
	err = m.BuildLambdaFunction()
	if nil != err {
		logrus.Errorf("build lambda function failed. %s.", err)
		return
	}

	// api gateway event
	switch projConfig.EventSourceType {
	case ApiGatewayProxyEvent:
		err = m.BuildApiGatewayProxyEvent()
		if nil != err {
			logrus.Errorf("build api gateway proxy event failed. %s.", err)
			return
		}
	}

	return
}

func (m *SAMTemplateConfig) Save() (err error) {
	stage := m.State
	projectPath := m.ProjectPath

	samYamlFilePath := fmt.Sprintf("%s/deploy/%s/template.yaml", projectPath, stage)
	byteYaml, err := yaml.Marshal(m.SAMTemplateYamlFile)
	if nil != err {
		logrus.Errorf("marshal sam yaml file failed. \n%s.", err)
		return
	}

	err = ioutil.WriteFile(samYamlFilePath, byteYaml, project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write sam yaml file failed. \n%s.", err)
		return
	}

	return
}
