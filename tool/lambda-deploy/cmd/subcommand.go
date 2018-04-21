package cmd

import (
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-deploy/cmd/remote"
)

func init() {
	rootCmd.AddCommand(remote.CommandRomote())
}
