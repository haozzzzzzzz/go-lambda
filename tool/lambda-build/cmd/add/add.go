package add

import (
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/add/api"
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/add/authorizer"
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/add/function"
	"github.com/spf13/cobra"
)

func CommandAdd() *cobra.Command {
	// addCmd represents the add command
	var addCmd = &cobra.Command{
		Use:     "add",
		Short:   "add component",
		Example: `add func --name LambdaHandler --path ./`,
	}

	addCmd.AddCommand(
		function.CommandAddLambdaFunction(),
		api.CommandAddApiFunction(),
		authorizer.CommandAddAuthorizerFunction(),
	)
	return addCmd
}
