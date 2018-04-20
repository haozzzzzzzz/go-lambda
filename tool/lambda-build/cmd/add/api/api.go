package api

import "github.com/spf13/cobra"

func CommandAddApi() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "api",
		Short: "add api",
		Run: func(cmd *cobra.Command, args []string) {
			// ...
		},
	}

	return cmd
}
