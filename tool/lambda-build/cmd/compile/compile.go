package compile

import (
	"github.com/spf13/cobra"
)

func CommandCompile() *cobra.Command {
	// compileCmd represents the compile command
	var compileCmd = &cobra.Command{
		Use:   "compile",
		Short: "compile lambda project",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return compileCmd
}
