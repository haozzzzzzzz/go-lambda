package cmd

import (
	"github.com/haozzzzzzzz/go-lambda/tools/lambda-build/cmd/add"
	"github.com/haozzzzzzzz/go-lambda/tools/lambda-build/cmd/compile"
	"github.com/haozzzzzzzz/go-lambda/tools/lambda-build/cmd/config"
)

func init() {
	rootCmd.AddCommand(add.CommandAdd())
	rootCmd.AddCommand(compile.CommandCompile())
	rootCmd.AddCommand(config.CommandConfig())

}
