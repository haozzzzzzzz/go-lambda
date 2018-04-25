package dynamodb

import (
	"github.com/haozzzzzzzz/go-lambda/resource"
)

func init() {
	resource.RegisterResource(resource.DynamoDBResourceType)
}
