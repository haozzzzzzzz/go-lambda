package init

import (
	"os"
	"path/filepath"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandInit() *cobra.Command {
	var projectPath string
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "init lambda function project",
		Run: func(cmd *cobra.Command, args []string) {
			if projectPath == "" {
				logrus.Errorf("need path")
				return
			}

			var err error
			projectPath, err = filepath.Abs(projectPath)
			if nil != err {
				logrus.Errorf("get absolute project path failed. \ns%s.", err)
				return
			}

			awsYamlFile := proj.AWSYamlFile{
				Mode: os.ModePerm,
			}
			_, err = awsYamlFile.CheckAWSYamlFile(projectPath)
			if nil != err {
				logrus.Errorf("check yaml file failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectPath, "path", "p", "./", "project path")

	return cmd
}
