package _func

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// 添加lambda函数命令
func CommandAddLambdaFunction() *cobra.Command {
	var handler LambdaFunction
	var cmd = &cobra.Command{
		Use:   "func",
		Short: "add lambda function",
		Run: func(cmd *cobra.Command, args []string) {
			err := handler.Run()
			if nil != err {
				logrus.Errorf("run add lambda function command failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&handler.Name, "name", "n", "", "set lambda project name")
	flags.StringVarP(&handler.Path, "path", "p", "", "set lambda project path")

	return cmd
}

// 添加Lambda函数命令处理器
type LambdaFunction struct {
	Name string `json:"name" validate:"required"`
	Path string `json:"path" validate:"required"`
}

func (m *LambdaFunction) Run() (err error) {
	err = validator.New().Struct(m)
	if nil != err {
		logrus.Errorf("validate struct failed. \n%s.", err)
		return
	}

	// make project directory
	dir, err := filepath.Abs(m.Path)
	if nil != err {
		logrus.Errorf("get absolute file path failed. \n%s.", err)
		return
	}

	projectDirectory := fmt.Sprintf("%s/%s", dir, m.Name)
	if file.PathExists(projectDirectory) {
		err = errors.New("project directory has existed")
		if nil != err {
			return
		}
	}

	// project root
	err = os.MkdirAll(projectDirectory, os.ModePerm)
	if nil != err {
		logrus.Errorf("make project directory failed. \n%s.", err)
		return
	}

	// project handlers
	err = os.MkdirAll(fmt.Sprintf("%s/handler", projectDirectory), os.ModePerm)
	if nil != err {
		logrus.Errorf("make project handler directory failed. \n%s.", err)
		return
	}

	//create main file

	return
}
