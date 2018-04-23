package remote

import (
	"github.com/spf13/cobra"
)

func CommandRemote() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "remote",
		Short:   "remote command",
		Example: "remote lamdba",
	}

	cmd.AddCommand(CommandRemoteFunction())

	return cmd
}
