package cmd

import (
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/add"
	"github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/compile"
	_init "github.com/haozzzzzzzz/go-lambda/tool/lambda-build/cmd/init"
)

func init() {
	rootCmd.AddCommand(add.CommandAdd())
	rootCmd.AddCommand(compile.CommandCompile())
	rootCmd.AddCommand(_init.CommandInit())

}
