package proj

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type AWSYamlFile struct {
	AccountId    string `yaml:"account_id"`
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
	Region       string `yaml:"region"`
	OutputFormat string `yaml:"output_format"`
	Role         string `yaml:"role"`
	CodeS3Bucket string `yaml:"code_s3_bucket"`
}

func (m *AWSYamlFile) Save(projectPath string) (err error) {
	awsYamlFileDir := fmt.Sprintf("%s/.proj/secret", projectPath)
	if !file.PathExists(awsYamlFileDir) {
		err = os.MkdirAll(awsYamlFileDir, project.ProjectDirMode)
		if nil != err {
			logrus.Errorf("make project secret directory failed. \n%s.", err)
			return
		}
	}

	awsYamlFilePath := fmt.Sprintf("%s/aws.yaml", awsYamlFileDir)
	byteYaml, err := yaml.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal aws yaml failed. \n%s.", err)
		return
	}

	err = ioutil.WriteFile(awsYamlFilePath, byteYaml, project.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write aws yaml file failed. \n%s.", err)
		return
	}

	return
}

func LoadAWSYamlFile(projectPath string) (yamlFile *AWSYamlFile, err error) {
	awsYamlFilePath := fmt.Sprintf("%s/.proj/secret/aws.yaml", projectPath)
	byteYaml, err := ioutil.ReadFile(awsYamlFilePath)
	if nil != err {
		logrus.Errorf("read aws yaml file failed. \n%s.", err)
		return
	}

	yamlFile = &AWSYamlFile{}
	err = yaml.Unmarshal(byteYaml, yamlFile)
	if nil != err {
		logrus.Errorf("unmarshal aws yaml file failed. \n%s.", err)
		return
	}

	return
}

func CheckAWSYamlFile(projectPath string, overwrite bool) (awsYamlFile *AWSYamlFile, exist bool, err error) {
	awsYamlFilePath := fmt.Sprintf("%s/.proj/secret/aws.yaml", projectPath)
	if file.PathExists(awsYamlFilePath) {
		exist = true
		awsYamlFile, err = LoadAWSYamlFile(projectPath)
		if nil != err {
			logrus.Errorf("load asw yaml file failed. \n%s.", err)
			return
		}

		// 如果不重新配置，则返回
		if !overwrite {
			return
		}

	} else {
		awsYamlFile = &AWSYamlFile{}
	}

	inputReader := bufio.NewReader(os.Stdin)
	var input string

	fmt.Print(fmt.Sprintf("Input AWS account id (%s):", awsYamlFile.AccountId))
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS account id failed. \n%s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		awsYamlFile.AccountId = input
	}

	fmt.Print(fmt.Sprintf("Input AWS access key (%s):", awsYamlFile.AccessKey))
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS access key failed. \n%s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		awsYamlFile.AccessKey = input
	}

	fmt.Print(fmt.Sprintf("Input AWS secret key (%s):", awsYamlFile.SecretKey))
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS secret key failed. \n%s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		awsYamlFile.SecretKey = input
	}

	fmt.Print(fmt.Sprintf("Input AWS region (%s):", awsYamlFile.Region))
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS region failed. \n%s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		awsYamlFile.Region = input
	}

	fmt.Print(fmt.Sprintf("Input AWS output format (%s):", awsYamlFile.OutputFormat))
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS output format failed. \n%s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		awsYamlFile.OutputFormat = input
	}

	fmt.Print(fmt.Sprintf("Input Lambda Execution Role(%s):", awsYamlFile.Role))
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read lambda execution role failed. \n%s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		awsYamlFile.Role = input
	}

	// code deploy s3 bucket
	fmt.Print(fmt.Sprintf("Input Code Deploy S3 Bucket(%s):", awsYamlFile.CodeS3Bucket))
	input, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read lambda execution role failed. \n%s.", err)
		return
	}
	input = strings.Replace(input, "\n", "", -1)
	if input != "" {
		awsYamlFile.CodeS3Bucket = input
	}

	err = awsYamlFile.Save(projectPath)
	if nil != err {
		logrus.Errorf("save AWS yaml file failed. \n%s.", err)
		return
	}

	return
}

// 运行AWS Cli命令时需要设置环境变量
func (m *AWSYamlFile) RunAWSCliCommand(command string, args ...string) (exit int, err error) {
	os.Setenv("AWS_ACCESS_KEY_ID", m.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", m.SecretKey)
	os.Setenv("AWS_DEFAULT_REGION", m.Region)
	os.Setenv("AWS_DEFAULT_OUTPUT", m.OutputFormat)
	return cmd.RunCommand("./", command, args...)
}
