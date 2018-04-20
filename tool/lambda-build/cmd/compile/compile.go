package compile

import (
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/compile/api"
	"github.com/spf13/cobra"
)

func CommandCompile() *cobra.Command {
	// compileCmd represents the compile command
	var compileCmd = &cobra.Command{
		Use:     "compile",
		Short:   "compile component",
		Example: `compile api --path ./`,
	}

	compileCmd.AddCommand(
		api.CommandCompileApiFunction(),
	)

	return compileCmd
}
