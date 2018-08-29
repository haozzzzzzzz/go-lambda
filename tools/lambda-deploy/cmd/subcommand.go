package cmd

import (
	"github.com/haozzzzzzzz/go-lambda/tools/lambda-deploy/cmd/remote"
)

func init() {
	rootCmd.AddCommand(remote.CommandRemote())
}
