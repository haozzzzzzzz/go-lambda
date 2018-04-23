package remote

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/serenize/snaker.git"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandRemoteFunction() *cobra.Command {
	var projectPath string
	var cmd = &cobra.Command{
		Use:     "func",
		Short:   "deploy remote lambda function command",
		Example: "remote func -p ./",
		Run: func(cmd *cobra.Command, args []string) {
			if projectPath == "" {
				logrus.Errorf("need specify project path")
				return
			}

			var err error
			projectPath, err = filepath.Abs(projectPath)
			if nil != err {
				logrus.Errorf("get absolute project path failed. \n%s.", err)
				return
			}

			remoteLambdaFunc := &RemoteLambdaFunction{
				ProjectPath: projectPath,
				Mode:        os.ModePerm,
			}
			err = remoteLambdaFunc.Run()
			if nil != err {
				logrus.Errorf("run remote lambda function failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectPath, "path", "p", "./", "project path")

	return cmd
}

type RemoteLambdaFunction struct {
	ProjectPath string      `yaml:"project_path"`
	Mode        os.FileMode `yaml:"mode"`
}

func (m *RemoteLambdaFunction) Run() (err error) {
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
	packageBucket := fmt.Sprintf("lambda-%s", snaker.CamelToSnake(projectConfig.Name))
	logrus.Infof("checking s3 bucket %q", packageBucket)

	packageBucket = strings.Replace(packageBucket, "_", "-", -1)
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
	_, err = awsYamlFile.RunAWSCliCommand(
		"aws",
		"cloudformation", "package",
		"--template-file", fmt.Sprintf("%s/deploy/template.yaml", projectPath),
		"--output-template-file", fmt.Sprintf("%s/deploy/serverless-output.yaml", projectPath),
		"--s3-bucket", packageBucket,
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
		"--template-file", fmt.Sprintf("%s/deploy/serverless-output.yaml", projectPath),
		"--stack-name", stackName,
	)
	if nil != err {
		logrus.Errorf("cloudformation deploy lambda function failed. \n%s.", err)
		return
	}

	// 成功
	logrus.Info("deploy lambda function successfully")

	return
}
