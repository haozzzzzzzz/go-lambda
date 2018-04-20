package api

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandCompileApiFunction() *cobra.Command {
	var projectPath string
	var cmd = &cobra.Command{
		Use:   "api",
		Short: "compile api",
		Run: func(cmd *cobra.Command, args []string) {
			if projectPath == "" {
				logrus.Errorf("need path")
				return
			}

			err := api.NewApiParser(projectPath).MapApi()
			if nil != err {
				logrus.Errorf("mapping api failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectPath, "path", "p", "./", "project path")

	return cmd
}
