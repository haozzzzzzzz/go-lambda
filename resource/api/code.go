package api

import (
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

var (
	// DynamoDB错误
	CodeErrorDynamoDB = &ginbuilder.ReturnCode{
		Code:    2000,
		Message: "dynamodb error",
	}

	// 获取DynamoDB记录失败
	CodeErrorDynamoDBGetItemFailed = &ginbuilder.ReturnCode{
		Code:    2001,
		Message: "dynamodb get item failed",
	}

	// 添加DynamoDB记录失败
	CodeErrorDynamoDBPutItemFailed = &ginbuilder.ReturnCode{
		Code:    2002,
		Message: "dynamodb put item failed",
	}
)
