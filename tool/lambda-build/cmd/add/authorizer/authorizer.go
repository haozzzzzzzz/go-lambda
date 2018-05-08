package authorizer

import (
	"path/filepath"

	"github.com/haozzzzzzzz/go-lambda/proj"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandAddAuthorizerFunction() *cobra.Command {
	var projectPath string
	cmd := &cobra.Command{
		Use:   "authorizer",
		Short: "add apig authorizer",
		Run: func(cmd *cobra.Command, args []string) {
			projectPath, err := filepath.Abs(projectPath)
			if nil != err {
				logrus.Errorf("get absolute file path failed. %s.", err)
				return
			}

			authorizerConfig, err := proj.NewAuthorizerConfigFromProjPath(projectPath)
			if nil != err {
				logrus.Errorf("new authorizer config failed. %s.", err)
				return
			}

			authorizer := &proj.Authorizer{}
			err = authorizer.StdinBuild()
			if nil != err {
				logrus.Errorf("build authorizer failed. %s.", err)
				return
			}

			authorizerConfig.AuthorizerYamlFile.AddAuthorizer(authorizer)
			err = authorizerConfig.Save()
			if nil != err {
				logrus.Errorf("save authorizer config failed. %s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectPath, "path", "p", "./", "source file path")

	return cmd
}
