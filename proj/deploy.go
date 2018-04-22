package proj

import (
	"os"
	"time"
)

// 发布
type DeployYamlFile struct {
	Version    os.FileMode `yaml:"version"`
	UpdateTime time.Time   `yaml:"update_time"`
	Mode       os.FileMode `yaml:"mode"`
}

func (m *DeployYamlFile) Save(projectPath string) (err error) {
	return
}

func LoadDeployYamlFile(projectPath string) (err error) {
	return
}
