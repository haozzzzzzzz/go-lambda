package proj

import (
	"io/ioutil"
	"os"

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

func (m *AWSYamlFile) Save(filePath string) (err error) {
	byteYaml, err := yaml.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal aws yaml failed. \n%s.", err)
		return
	}

	err = ioutil.WriteFile(filePath, byteYaml, m.Mode)
	if nil != err {
		logrus.Errorf("write aws yaml file failed. \n%s.", err)
		return
	}

	return
}

func LoadAWSYamlFile(filePath string) (yamlFile *AWSYamlFile, err error) {
	byteYaml, err := ioutil.ReadFile(filePath)
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
