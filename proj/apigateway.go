package proj

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
)

const (
	AuthorizerFileName = "authorizer.yaml"
)

type AuthorizerConfig struct {
	ProjectPath             string
	ConfigFilePath          string
	ProjectYamlFile         *ProjectYamlFile
	AuthorizerYamlFile      *AuthorizerYamlFile
	AuthorizerYamlFileExsit bool
	Mode                    os.FileMode
}

func NewAuthorizerConfigFromProjPath(projPath string) (config *AuthorizerConfig, err error) {
	projYamlFile, err := LoadProjectYamlFile(projPath)
	if nil != err {
		logrus.Errorf("load project yaml file failed. %s.", err)
		return
	}

	config, err = NewAuthorizerConfig(projYamlFile)
	if nil != err {
		logrus.Errorf("new authorizer config from project yaml file failed. %s.", err)
		return
	}

	return
}

func NewAuthorizerConfig(projYamlFile *ProjectYamlFile) (config *AuthorizerConfig, err error) {
	if projYamlFile.EventSourceType != ApiGatewayProxyEvent {
		logrus.Info("project's event trigger should be ApiGatewayProxyEvent")
		return
	}

	projPath := projYamlFile.ProjectPath

	config = &AuthorizerConfig{
		ProjectPath:             projPath,
		ProjectYamlFile:         projYamlFile,
		ConfigFilePath:          fmt.Sprintf("%s/.proj/%s", projPath, AuthorizerFileName),
		AuthorizerYamlFile:      NewAuthorizerYamlFile(),
		AuthorizerYamlFileExsit: false,
		Mode: projYamlFile.Mode,
	}

	config.Load()

	return
}

func (m *AuthorizerConfig) Load() (err error) {
	if file.PathExists(m.ConfigFilePath) {
		m.AuthorizerYamlFileExsit = true
		err = yaml.ReadYamlFromFile(m.ConfigFilePath, m.AuthorizerYamlFile)
		if nil != err {
			logrus.Errorf("read %q failed. %s.", m.ConfigFilePath, err)
			return
		}
	}

	return
}

func (m *AuthorizerConfig) Save() (err error) {
	err = yaml.WriteYamlToFile(m.ConfigFilePath, m.AuthorizerYamlFile, m.Mode)
	if nil != err {
		logrus.Errorf("write authorizer %q failed. %s.", m.ConfigFilePath, err)
		return
	}
	return
}

type Authorizer struct {
	Name            string   `yaml:"name"`             // 名字
	FunctionHandler string   `yaml:"function_handler"` // 处理函数
	Headers         []string `yaml:"headers"`          // header校验字段名
	Queries         []string `yaml:"queries"`          // query参数名
}

func (m *Authorizer) StdinBuild() (err error) {
	inputReader := bufio.NewReader(os.Stdin)
	var input string

	fmt.Print("Input authorizer name: ")
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read authorizer name failed. %s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		m.Name = input
	}

	fmt.Print("Input authorizer function handler: ")
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read authorizer function handler failed. %s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		m.FunctionHandler = input
	}

	fmt.Print("Input authorizer headers(split by ,): ")
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read authorizer headers failed. %s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		m.Headers = strings.Split(input, ",")
	}

	fmt.Print("Input authorizer queries(split by ,): ")
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read authorizer queries failed. %s.", err)
		return
	}

	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		m.Queries = strings.Split(input, ",")
	}

	return
}

type AuthorizerYamlFile struct {
	Authorizers []*Authorizer `yaml:"authorizers"`
}

func NewAuthorizerYamlFile() *AuthorizerYamlFile {
	return &AuthorizerYamlFile{}
}

func (m *AuthorizerYamlFile) AddAuthorizer(authorizer *Authorizer) {
	m.Authorizers = append(m.Authorizers, authorizer)
}
