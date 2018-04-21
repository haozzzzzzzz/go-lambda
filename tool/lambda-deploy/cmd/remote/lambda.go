package remote

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandRemoteLambda() *cobra.Command {
	var projectPath string
	var cmd = &cobra.Command{
		Use:     "lambda",
		Short:   "remote lambda command",
		Example: "remote lambda -p ./",
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

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectPath, "path", "p", "./", "project path")

	return cmd
}
