package resource

type ResourceType map[string]interface{}

var Resources map[string]ResourceType

func AddResource(resourceName string, resource ResourceType) {
	Resources[resourceName] = resource
}
