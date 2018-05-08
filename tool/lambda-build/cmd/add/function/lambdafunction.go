package function

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// 添加lambda函数命令
func CommandAddLambdaFunction() *cobra.Command {
	var handler LambdaFunction
	handler.Mode = os.ModePerm
	var eventType string
	var cmd = &cobra.Command{
		Use:   "func",
		Short: "add lambda function",
		Run: func(cmd *cobra.Command, args []string) {
			handler.EventSourceType = proj.NewLambdaFunctionEventSourceType(eventType)
			err := handler.Run()
			if nil != err {
				logrus.Errorf("run add lambda function command failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&handler.Name, "name", "n", "", "set lambda project name")
	flags.StringVarP(&handler.Description, "description", "d", "AWS Serverless Function", "lambda function description")
	flags.StringVarP(&handler.Path, "path", "p", "./", "set lambda project path")
	flags.StringVarP(&eventType, "event", "e", proj.BasicExecutionEvent.String(), "set lambda function event source type")

	return cmd
}

// 添加Lambda函数命令处理器
type LambdaFunction struct {
	Name            string      `json:"name" validate:"required"`
	Description     string      `yaml:"description"`
	Path            string      `json:"path" validate:"required"`
	Mode            os.FileMode `json:"mode" validate:"required"`
	ProjectPath     string      `json:"project_path"`
	EventSourceType proj.LambdaFunctionEventSourceType
}

func (m *LambdaFunction) Run() (err error) {
	/**
	这并不是Go的Bug，包括Linux系统调用都是这样的，创建目录除了给定的权限还要加上系统的Umask，Go也是如实遵循这种约定。
	Umask是权限的补码,用于设置创建文件和文件夹默认权限的,一般在 /etc/profile中或 $HOME/profile或 $HOME/.bash_profile中
	*/
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	dir, err := filepath.Abs(m.Path)
	if nil != err {
		logrus.Errorf("get absolute file path failed. \n%s.", err)
		return
	}
	m.ProjectPath = fmt.Sprintf("%s/%s", dir, m.Name)

	// 名字由字母和数字组成，由字母开始
	matched, err := regexp.MatchString("^[A-za-z][A-Za-z0-9]+$", m.Name)
	if nil != err {
		logrus.Errorf("run regexp match string failed. %s.", err)
		return
	}

	if !matched {
		err = errors.New("function name should be composed of letters and numbers, and it begins with letter")
		return
	}

	err = validator.New().Struct(m)
	if nil != err {
		logrus.Errorf("validate struct failed. \n%s.", err)
		return
	}

	if file.PathExists(m.ProjectPath) {
		err = errors.New("project directory has existed")
		if nil != err {
			return
		}
	}

	mode := m.Mode
	// project root
	err = os.MkdirAll(m.ProjectPath, mode)
	if nil != err {
		logrus.Errorf("make project directory failed. \n%s.", err)
		return
	}

	err = generateProjTemplate(m)
	if nil != err {
		logrus.Errorf("generate proj template failed. \n%s.", err)
		return
	}

	// config
	err = generateConfigTemplate(m)
	if nil != err {
		logrus.Errorf("generate config template failed. \n%s.", err)
		return
	}

	// constant
	err = generateConstantTemplate(m)
	if nil != err {
		logrus.Errorf("generate constant template failed. \n%s.", err)
		return
	}

	// handler
	err = generateHandlerTemplate(m)
	if nil != err {
		logrus.Errorf("generate handler template failed. \n%s.", err)
		return
	}

	switch m.EventSourceType {
	case proj.BasicExecutionEvent:
	case proj.CustomEvent:
	case proj.ApiGatewayProxyEvent:
		err = generateApiTemplate(m)
		if nil != err {
			logrus.Errorf("generate api template failed. \n%s.", err)
			return
		}

	case proj.ApiGatewayAuthorizerEvent:
		err = generateApiGatewayAuthorizer(m)
		if nil != err {
			logrus.Errorf("generate api gateway authorizer failed. %s.", err)
			return
		}
	}

	// create detector mai file
	err = generateDetectorMainTemplate(m)
	if nil != err {
		logrus.Errorf("generate project detector template failed. \n%s.", err)
		return
	}

	//create main file
	err = generateMainTemplate(m)
	if nil != err {
		logrus.Warnf("generate main template failed. \n%s.", err)
		return
	}

	// do go imports
	goimports.DoGoImports([]string{m.ProjectPath}, true)
	return
}
