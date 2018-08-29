package cmd

import (
	"github.com/haozzzzzzzz/go-lambda/tools/lamd/cmd/remote"
)

func init() {
	rootCmd.AddCommand(remote.CommandRemote())
}
