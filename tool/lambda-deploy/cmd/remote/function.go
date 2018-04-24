package remote

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-lambda/proj"
	time2 "github.com/haozzzzzzzz/go-rapid-development/utils/time"
	"github.com/serenize/snaker.git"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandRemoteFunction() *cobra.Command {
	var remoteLambdaFunction RemoteLambdaFunction
	var cmd = &cobra.Command{
		Use:     "func",
		Short:   "deploy remote lambda function command",
		Example: "remote func -p ./",
		Run: func(cmd *cobra.Command, args []string) {
			if remoteLambdaFunction.ProjectPath == "" {
				logrus.Errorf("need specify project path")
				return
			}

			var err error
			remoteLambdaFunction.ProjectPath, err = filepath.Abs(remoteLambdaFunction.ProjectPath)
			if nil != err {
				logrus.Errorf("get absolute project path failed. \n%s.", err)
				return
			}

			remoteLambdaFunction.Mode = os.ModePerm
			err = remoteLambdaFunction.Run()
			if nil != err {
				logrus.Errorf("run remote lambda function failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&remoteLambdaFunction.ProjectPath, "path", "p", "./", "project path")
	flags.StringVarP(&remoteLambdaFunction.Stage, "stage", "s", "test", "deploy stage name. \"test\" or \"prod\"")
	return cmd
}

type RemoteLambdaFunction struct {
	ProjectPath string      `yaml:"project_path" validate:"required"`
	Stage       string      `yaml:"stage" validate:"required"`
	Mode        os.FileMode `yaml:"mode"`
}

func (m *RemoteLambdaFunction) Run() (err error) {

	err = validator.New().Struct(m)
	if nil != err {
		logrus.Errorf("validate RemoteLambdaFunction failed. \n%s.", err)
		return
	}

	projectPath := m.ProjectPath
	// read aws config
	awsYamlFile, _, err := proj.CheckAWSYamlFile(projectPath, os.ModePerm, false)
	if nil != err {
		logrus.Errorf("check aws yaml file failed. \n%s.", err)
		return
	}

	// project yaml config
	projectConfig, err := proj.LoadProjectYamlFile(projectPath)
	if nil != err {
		logrus.Errorf("load project yaml file failed. \n%s.", err)
		return
	}

	// 创建一个s3的bucket用于存放代码包
	packageBucket := fmt.Sprintf("lambda-%s-%s", snaker.CamelToSnake(projectConfig.Name), m.Stage)
	packageBucket = strings.Replace(packageBucket, "_", "-", -1)
	logrus.Infof("checking s3 bucket %q", packageBucket)

	_, err = awsYamlFile.RunAWSCliCommand(
		"aws",
		"s3", "mb", fmt.Sprintf("s3://%s", packageBucket),
		"--region", awsYamlFile.Region,
	)
	if nil != err {
		logrus.Errorf("create aws s3 bucket %s failed. \n%s.", packageBucket, err)
		return
	}

	// 打包
	logrus.Info("packaging function")
	stageDeployPath := fmt.Sprintf("%s/deploy/%s", projectPath, m.Stage)
	dayStartTime, err := time2.DayStartTime(time.Now())
	if nil != err {
		logrus.Errorf("get day start time failed. \n%s.", err)
		return
	}

	strDayStartTime := time2.DateStringFormat(dayStartTime)
	_, err = awsYamlFile.RunAWSCliCommand(
		"aws",
		"cloudformation", "package",
		"--template-file", fmt.Sprintf("%s/template.yaml", stageDeployPath),
		"--output-template-file", fmt.Sprintf("%s/serverless-output.yaml", stageDeployPath),
		"--s3-bucket", packageBucket,
		"--s3-prefix", strDayStartTime,
	)
	if nil != err {
		logrus.Errorf("cloudformation package lambda function failed. \n%s.", err)
		return
	}

	// 发布
	logrus.Info("deploying package")
	stackName := packageBucket
	_, err = awsYamlFile.RunAWSCliCommand(
		"aws",
		"cloudformation", "deploy",
		"--template-file", fmt.Sprintf("%s/serverless-output.yaml", stageDeployPath),
		"--stack-name", stackName,
		"--capabilities", "CAPABILITY_IAM",
	)
	if nil != err {
		logrus.Errorf("cloudformation deploy lambda function failed. \n%s.", err)
		return
	}

	// 成功
	logrus.Info("deploy lambda function successfully")

	return
}
