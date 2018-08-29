package cmd

import (
	"github.com/haozzzzzzzz/go-lambda/tools/lamb/cmd/add"
	"github.com/haozzzzzzzz/go-lambda/tools/lamb/cmd/compile"
	"github.com/haozzzzzzzz/go-lambda/tools/lamb/cmd/config"
)

func init() {
	rootCmd.AddCommand(add.CommandAdd())
	rootCmd.AddCommand(compile.CommandCompile())
	rootCmd.AddCommand(config.CommandConfig())

}
