package local

import (
	"github.com/spf13/cobra"
)

func CommandLocal() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "local",
		Short:   "local command",
		Example: "local func",
	}

	return cmd
}
