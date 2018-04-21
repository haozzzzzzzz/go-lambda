package function

import (
	"fmt"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandCompileFunction() *cobra.Command {
	var compileFunction CompileFunction
	var cmd = &cobra.Command{
		Use:   "func",
		Short: "compile func",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			compileFunction.ProjectPath, err = filepath.Abs(compileFunction.ProjectPath)
			if nil != err {
				logrus.Errorf("get absolute project path failed. \n%s.", err)
				return
			}

			err = compileFunction.Run()
			if nil != err {
				logrus.Errorf("run compile func failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&compileFunction.ProjectPath, "path", "p", "./", "project path")

	return cmd
}

// compile function command
type CompileFunction struct {
	ProjectPath     string `json:"project_path" yaml:"project_path" validate:"required"`
	ProjectYamlFile *proj.ProjectYamlFile
}

func (m *CompileFunction) Run() (err error) {
	err = validator.New().Struct(m)
	if nil != err {
		logrus.Errorf("validate CompileFunction object failed. \n%s.", err)
		return
	}

	m.ProjectYamlFile, err = proj.LoadProjectYamlFile(m.ProjectPath)
	if nil != err {
		logrus.Errorf("load project yaml file failed. \n%s.", err)
		return
	}

	// go build
	err = m.runGoBuild()
	if nil != err {
		logrus.Errorf("run go build failed. \n%s.", err)
		return
	}
	return
}

func (m *CompileFunction) runGoBuild() (err error) {
	projConfig := m.ProjectYamlFile
	binTarget := fmt.Sprintf("%s/bin/%s", projConfig.ProjectPath, projConfig.Name)
	mainFile := fmt.Sprintf("%s/main.go", projConfig.ProjectPath)

	// go build
	logrus.Info("go building binary")
	exit, output, err := cmd.RunCommand("go", "build", "-v", "-o", binTarget, mainFile)
	if nil != err || exit != 0 {
		logrus.Errorf("run go build command failed. \n%s.", err)
		return
	}

	strOutput := output.String()
	if strOutput != "" {
		logrus.Info(strOutput)
	}

	// zip
	logrus.Info("zip building zip file")
	zipTarget := fmt.Sprintf("%s/bin/%s.zip", projConfig.ProjectPath, projConfig.Name)
	exit, output, err = cmd.RunCommand("zip", zipTarget, binTarget)
	if nil != err || exit != 0 {
		logrus.Errorf("run zip command failed. \n%s.", err)
		return
	}

	strOutput = output.String()
	if strOutput != "" {
		logrus.Info(strOutput)
	}

	return
}
