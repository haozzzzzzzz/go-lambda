package cmd

import (
	_init "github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/init"
)

func init() {
	rootCmd.AddCommand(_init.CommandInit())
}
