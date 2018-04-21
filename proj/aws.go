package proj

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type AWSYamlFile struct {
	AccessKey    string      `yaml:"access_key"`
	SecretKey    string      `yaml:"secret_key"`
	Region       string      `yaml:"region"`
	OutputFormat string      `yaml:"output_format"`
	Mode         os.FileMode `yaml:"mode"`
}

func (m *AWSYamlFile) Save(projectPath string) (err error) {
	awsYamlFilePath := fmt.Sprintf("%s/.proj/secret/aws.yaml", projectPath)
	byteYaml, err := yaml.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal aws yaml failed. \n%s.", err)
		return
	}

	err = ioutil.WriteFile(awsYamlFilePath, byteYaml, m.Mode)
	if nil != err {
		logrus.Errorf("write aws yaml file failed. \n%s.", err)
		return
	}

	return
}

func LoadAWSYamlFile(projectPath string) (yamlFile *AWSYamlFile, err error) {
	awsYamlFilePath := fmt.Sprintf("%s/secret/aws.yaml")
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

func (m *AWSYamlFile) CheckAWSYamlFile(projectPath string, overWrite bool) (exist bool, err error) {
	awsYamlFilePath := fmt.Sprintf("%s/secret/aws.yaml")
	if file.PathExists(awsYamlFilePath) && !overWrite {
		exist = true
		return
	}

	awsYamlFile := AWSYamlFile{
		Mode: m.Mode,
	}
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Print("Input AWS access key:")
	awsYamlFile.AccessKey, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS access key failed. \n%s.", err)
		return
	}

	fmt.Print("Input AWS secret key:")
	awsYamlFile.SecretKey, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS secret key failed. \n%s.", err)
		return
	}

	fmt.Print("Input AWS region(e.g. us-east-1):")
	awsYamlFile.Region, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS region failed. \n%s.", err)
		return
	}

	fmt.Print("Input AWS output format(e.g. json):")
	awsYamlFile.OutputFormat, err = inputReader.ReadString('\n')
	if nil != err {
		logrus.Errorf("read AWS output format failed. \n%s.", err)
		return
	}

	err = awsYamlFile.Save(projectPath)
	if nil != err {
		logrus.Errorf("save AWS yaml file failed. \n%s.", err)
		return
	}

	return
}
