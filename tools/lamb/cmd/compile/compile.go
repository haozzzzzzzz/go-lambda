package compile

import (
	"github.com/haozzzzzzzz/go-lambda/tools/lamb/cmd/compile/function"
	"github.com/spf13/cobra"
)

func CommandCompile() *cobra.Command {
	// compileCmd represents the compile command
	var compileCmd = &cobra.Command{
		Use:   "compile",
		Short: "compile component",
		Example: `compile api --path ./
compile function --path ./`,
	}

	compileCmd.AddCommand(
		function.CommandCompileFunction(),
	)

	return compileCmd
}
