package resource

type ResourceType int8

const (
	DynamoDBResourceType ResourceType = 1
	KinesisResourceType  ResourceType = 2
	XRayResourceType     ResourceType = 3
)

// 在使用中的资源
var resourceInUse []ResourceType

func GetResourceInUse() []ResourceType {
	return resourceInUse
}

func RegisterResource(resourceType ResourceType) {
	resourceInUse = append(resourceInUse, resourceType)
}
