package _func

// 创建main文件
func AddMainFile(fileName string) (err error) {
	return
}

var mainFileText = `
import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.MainHandler)
}
`
