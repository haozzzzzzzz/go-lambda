package proj

// project yaml 文件
type ProjectYamlFile struct {
	Name        string `json:"name" yaml:"name"`
	ProjectPath string `json:"project_path" yaml:"project_path"`
}
