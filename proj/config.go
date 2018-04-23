package proj

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Lambda函数事件源
type LambdaFunctionEventSourceType int8

const (
	BasicExecutionEvent LambdaFunctionEventSourceType = 0 // 基本执行
	CustomEvent         LambdaFunctionEventSourceType = 1 // 自定义事件
	ApiGatewayEvent     LambdaFunctionEventSourceType = 2 // API GATEWAY事件
)

func NewLambdaFunctionEventSourceType(strEvent string) LambdaFunctionEventSourceType {
	switch strEvent {
	case CustomEvent.String():
		return CustomEvent
	case ApiGatewayEvent.String():
		return ApiGatewayEvent
	case BasicExecutionEvent.String():
		return BasicExecutionEvent
	}
	return BasicExecutionEvent
}

func (m LambdaFunctionEventSourceType) String() string {
	switch m {
	case CustomEvent:
		return "CustomEvent"
	case ApiGatewayEvent:
		return "ApiGatewayEvent"
	case BasicExecutionEvent:
		fallthrough
	default:
		return "BasicExecutionEvent"
	}

	return ""
}

// project yaml 文件
type ProjectYamlFile struct {
	Name            string                        `json:"name" yaml:"name"`
	Description     string                        `yaml:"description"`
	ProjectPath     string                        `json:"project_path" yaml:"project_path"`
	EventSourceType LambdaFunctionEventSourceType `yaml:"event_source_type"`
	Mode            os.FileMode                   `json:"mode" yaml:"mode"`
}

func (m *ProjectYamlFile) Save() (err error) {
	projYamlFileName := fmt.Sprintf("%s/.proj/proj.yaml", m.ProjectPath)
	byteProjYamlFile, err := yaml.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal proj yaml config file failed. \n%s.", err)
		return
	}
	err = ioutil.WriteFile(projYamlFileName, byteProjYamlFile, m.Mode)
	if nil != err {
		logrus.Warnf("write .proj/proj.yaml failed. \n%s.", err)
		return
	}

	return
}

func LoadProjectYamlFile(projectPath string) (yamlFile *ProjectYamlFile, err error) {
	projYamlFileName := fmt.Sprintf("%s/.proj/proj.yaml", projectPath)
	byteProjYamlFile, err := ioutil.ReadFile(projYamlFileName)
	if nil != err {
		logrus.Errorf("read %q project yaml config file failed. \n%s.", err)
		return
	}

	yamlFile = &ProjectYamlFile{}
	err = yaml.Unmarshal(byteProjYamlFile, yamlFile)
	if nil != err {
		yamlFile = nil
		logrus.Errorf("unmarshal project yaml file failed. \n%s.", err)
		return
	}
	return
}
