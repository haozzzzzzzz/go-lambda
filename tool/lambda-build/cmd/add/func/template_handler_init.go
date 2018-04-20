package _func

func AddHandlerInitFile(fileName string) (err error) {
	return
}

var initFileText = `
// init program running
func init () {
	logrus.Infof("Lambda function %p initializing ...")
}

func MainHandler(
	ctx context.Context,
	request *events.APIGatewayProxyRequest,
) (
	response events.APIGatewayProxyResponse,
	err error
){
	
}
`
