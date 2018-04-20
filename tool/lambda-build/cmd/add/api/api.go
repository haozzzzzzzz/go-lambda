package api

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api"
	"github.com/serenize/snaker.git"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CommandAddApiFunction() *cobra.Command {
	var apiItem api.ApiItem
	var cmd = &cobra.Command{
		Use:   "api",
		Short: "add gin api",
		Run: func(cmd *cobra.Command, args []string) {
			path, err := filepath.Abs(apiItem.SourceFile)
			if nil != err {
				logrus.Errorf("get absolute file path failed. \n%s.", err)
				return
			}

			curDir := filepath.Base(path)
			if curDir == "api" && file.PathExists(fmt.Sprintf("%s/routers.go", path)) {
				logrus.Errorf("api file can not has same directory as routers.go")
				return
			}

			apiItem.ApiHandlerPackage = snaker.CamelToSnake(curDir)
			apiItem.SourceFile = fmt.Sprintf("%s/api_%s.go", path, strings.ToLower(apiItem.ApiHandlerFunc))
			err = validator.New().Struct(apiItem)
			if nil != err {
				logrus.Errorf("invalid api item. \n%s.", err)
				return
			}

			err = api.CreateApiSource(&apiItem)
			if nil != err {
				logrus.Errorf("create api source failed. \n%s.", err)
				return
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&apiItem.ApiHandlerFunc, "name", "n", "VarFuncName", "api handler function name")
	flags.StringVarP(&apiItem.HttpMethod, "method", "m", "GET", "GET POST PUT PATCH HEAD OPTIONS DELETE CONNECT TRACE")
	flags.StringVarP(&apiItem.RelativePath, "uri", "u", "/", "api uri")
	flags.StringVarP(&apiItem.SourceFile, "path", "p", "./", "source file path")
	return cmd
}
