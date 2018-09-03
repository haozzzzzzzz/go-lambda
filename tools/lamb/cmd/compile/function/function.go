package function

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/utils/str"
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
	flags.StringVarP(&compileFunction.Stage, "stage", "s", "dev", "stage name, [dev test pre prod]")
	return cmd
}

// compile function command
type CompileFunction struct {
	ProjectPath       string `json:"project_path" yaml:"project_path" validate:"required"`
	Stage             string `yaml:"stage" validate:"required"`
	ProjectYamlFile   *proj.ProjectYamlFile
	AWSYamlFile       *proj.AWSYamlFile
	SAMTemplateConfig *proj.SAMTemplateConfig
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
	m.AWSYamlFile, _, err = proj.CheckAWSYamlFile(m.ProjectPath, false)
	if nil != err {
		logrus.Errorf("check aws yaml file failed. \n%s.", err)
		return
	}

	// copy config
	err = m.copyConfig(fmt.Sprintf("%s/config", m.ProjectPath))
	if nil != err {
		logrus.Errorf("copy config folder failed. %s.", err)
		return
	}

	// go build
	err = m.runGoBuild()
	if nil != err {
		logrus.Errorf("run go build failed. \n%s.", err)
		return
	}

	// save sam template file
	m.SAMTemplateConfig, err = proj.NewSAMTempalteConfig(m.Stage, m.ProjectYamlFile, m.AWSYamlFile)
	if nil != err {
		logrus.Errorf("new sam template config failed. \n%s.", err)
		return
	}
	err = m.SAMTemplateConfig.Save()
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
	lambdaFuncName := fmt.Sprintf("%s%s", str.CamelString(m.Stage), projConfig.Name)
	deployTarget := fmt.Sprintf("%s/%s", projPath, lambdaFuncName)
	mainFile := fmt.Sprintf("%s/main.go", projPath)

	stageDeployFolder := fmt.Sprintf("%s/deploy/%s", projPath, m.Stage)
	err = os.MkdirAll(stageDeployFolder, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make stage deploy folder failed. %s.", err)
		return
	}

	// go run detector
	detectorMainFile := fmt.Sprintf("%s/.proj/detector/main.go", projPath)
	detectorMain := fmt.Sprintf("%s/detector", projPath)
	exit, err := cmd.RunCommand("./", "go", "build", "-v", "-o", detectorMain, detectorMainFile)
	if nil != err || exit != 0 {
		logrus.Errorf("build detector failed. \n%s.", err)
		return
	}

	roleYamlFilePath := fmt.Sprintf("%s/.proj/role.yaml", projPath)
	// detector only running on dev environment
	exit, err = cmd.RunCommand(fmt.Sprintf("%s/stage/dev/", projPath), detectorMain, "--path", roleYamlFilePath)
	if nil != err {
		logrus.Errorf("generate role.yaml failed. \n%s.", err)
		return
	}

	// go build
	logrus.Info("go building binary")
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")
	exit, err = cmd.RunCommand("./", "go", "build", "-v", "-o", deployTarget, mainFile)
	if nil != err || exit != 0 {
		logrus.Errorf("run go build command failed. \n%s.", err)
		return
	}

	return
}

func (m *CompileFunction) copyConfig(configDir string) (err error) {
	projectPath := m.ProjectPath

	// 配置源
	var stageConfigDir string
	switch m.Stage {
	case proj.DevStage.String():
		stageConfigDir = fmt.Sprintf("%s/stage/%s/config", projectPath, m.Stage)
	case proj.TestStage.String():
		stageConfigDir = fmt.Sprintf("%s/stage/%s/config", projectPath, m.Stage)
	case proj.PreStage.String():
		stageConfigDir = fmt.Sprintf("%s/stage/%s/config", projectPath, m.Stage)
	case proj.ProdStage.String():
		stageConfigDir = fmt.Sprintf("%s/stage/%s/config", projectPath, m.Stage)
	default:
		err = errors.New("not supported stage")
		return
	}

	// 设置config
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

	return
}

func (m *CompileFunction) zipPackage() (err error) {
	projConfig := m.ProjectYamlFile
	projectPath := projConfig.ProjectPath
	lambdaFuncName := fmt.Sprintf("%s%s", str.CamelString(m.Stage), projConfig.Name)

	// zip
	logrus.Info("zip building zip file")
	zipWorkPath := fmt.Sprintf("%s/deploy/%s", projectPath, m.Stage)
	zipFileName := fmt.Sprintf("%s.zip", lambdaFuncName)
	zipTarget := fmt.Sprintf("%s/%s", zipWorkPath, zipFileName)

	exit, err := cmd.RunCommand("./", "zip", "-r", zipTarget, lambdaFuncName, "config")
	if nil != err || exit != 0 {
		logrus.Errorf("run zip command failed. \n%s.", err)
		return
	}

	gitIgnoreFilePath := fmt.Sprintf("%s/.gitignore", zipWorkPath)
	if !file.PathExists(gitIgnoreFilePath) {
		gitIgnoreFileText := zipFileName
		err = ioutil.WriteFile(gitIgnoreFilePath, []byte(gitIgnoreFileText), project.ProjectFileMode)
		if nil != err {
			logrus.Errorf("write %q failed. %s.", gitIgnoreFilePath, err)
			return
		}
	}

	return
}
