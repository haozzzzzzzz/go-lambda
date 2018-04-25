package function

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/pkg/errors"
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
	flags.StringVarP(&compileFunction.Stage, "stage", "s", "test", "stage name, test or prod")
	return cmd
}

// compile function command
type CompileFunction struct {
	ProjectPath     string `json:"project_path" yaml:"project_path" validate:"required"`
	Stage           string `yaml:"stage" validate:"required"`
	ProjectYamlFile *proj.ProjectYamlFile
	AWSYamlFile     *proj.AWSYamlFile
	SAMYamlFile     *proj.SAMTemplateYamlFile
}

func (m *CompileFunction) Run() (err error) {
	err = validator.New().Struct(m)
	if nil != err {
		logrus.Errorf("validate CompileFunction object failed. \n%s.", err)
		return
	}

	// read project yaml file
	m.ProjectYamlFile, err = proj.LoadProjectYamlFile(m.ProjectPath)
	if nil != err {
		logrus.Errorf("load project yaml file failed. \n%s.", err)
		return
	}

	// read AWS yaml file
	m.AWSYamlFile, _, err = proj.CheckAWSYamlFile(m.ProjectPath, m.ProjectYamlFile.Mode, false)
	if nil != err {
		logrus.Errorf("check aws yaml file failed. \n%s.", err)
		return
	}

	// go build
	err = m.runGoBuild()
	if nil != err {
		logrus.Errorf("run go build failed. \n%s.", err)
		return
	}

	// save sam template file
	m.SAMYamlFile, err = proj.NewSAMTemplateYamlFileByExistConfig(m.Stage, m.ProjectYamlFile, m.AWSYamlFile)
	if nil != err {
		logrus.Errorf("new sam template yaml obj failed. \n%s.", err)
		return
	}
	err = m.SAMYamlFile.Save(m.Stage, m.ProjectPath, m.ProjectYamlFile.Mode)
	if nil != err {
		logrus.Errorf("save sam template failed. \n%s.", err)
		return
	}

	// zip package
	err = m.zipPackage()
	if nil != err {
		logrus.Errorf("zip package failed. \n%s.", err)
		return
	}

	return
}

func (m *CompileFunction) runGoBuild() (err error) {
	projConfig := m.ProjectYamlFile
	projPath := projConfig.ProjectPath
	deployTarget := fmt.Sprintf("%s/deploy/%s/%s", projPath, m.Stage, projConfig.Name)
	mainFile := fmt.Sprintf("%s/main.go", projPath)

	// go run detector
	detectorMainFile := fmt.Sprintf("%s/.proj/detector/main.go", projPath)
	detectorMain := fmt.Sprintf("%s/.proj/detector/main", projPath)
	exit, err := cmd.RunCommand("go", "build", "-v", "-o", detectorMain, detectorMainFile)
	if nil != err || exit != 0 {
		logrus.Errorf("build detector failed. \n%s.", err)
		return
	}

	roleYamlFilePath := fmt.Sprintf("%s/.proj/role.yaml", projPath)
	exit, err = cmd.RunCommand(detectorMain, "--path", roleYamlFilePath)
	if nil != err {
		logrus.Errorf("generate role.yaml failed. \n%s.", err)
		return
	}

	// go build
	logrus.Info("go building binary")
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")
	exit, err = cmd.RunCommand("go", "build", "-v", "-o", deployTarget, mainFile)
	if nil != err || exit != 0 {
		logrus.Errorf("run go build command failed. \n%s.", err)
		return
	}

	return
}

func (m *CompileFunction) zipPackage() (err error) {
	projConfig := m.ProjectYamlFile
	projectPath := projConfig.ProjectPath
	deployTarget := fmt.Sprintf("%s/deploy/%s/%s", projectPath, m.Stage, projConfig.Name)

	// 配置源
	var stageConfigDir string
	switch m.Stage {
	case proj.TestStage.String():
		stageConfigDir = fmt.Sprintf("%s/config_test", projectPath)
	case proj.ProdStage.String():
		stageConfigDir = fmt.Sprintf("%s/config_prod", projectPath)
	default:
		err = errors.New("not supported stage")
		return
	}

	// 设置config
	configDir := fmt.Sprintf("%s/config", projectPath)
	err = os.RemoveAll(configDir)
	if nil != err {
		logrus.Errorf("remove %q failed. \n%s.", configDir)
		return
	}

	err = file.Copy(stageConfigDir, configDir)
	if nil != err {
		logrus.Errorf("copy %q to %q failed. \n%s.", stageConfigDir, configDir, err)
		return
	}

	// zip
	logrus.Info("zip building zip file")
	zipTarget := fmt.Sprintf("%s/deploy/%s/%s.zip", projectPath, m.Stage, projConfig.Name)

	// 打包可执行文件和配置文件
	exit, err := cmd.RunCommand("zip", "-j", zipTarget, deployTarget, configDir)
	if nil != err || exit != 0 {
		logrus.Errorf("run zip command failed. \n%s.", err)
		return
	}

	return
}
