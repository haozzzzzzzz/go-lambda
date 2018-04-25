package iam

import (
	"fmt"
	"testing"

	_ "github.com/haozzzzzzzz/go-lambda/resource/dynamodb"

	"gopkg.in/yaml.v2"
)

func TestGetExecutionRole(t *testing.T) {
	role := NewExecutionRole("ExampleApiExecutionRole")
	byteRole, err := yaml.Marshal(role)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Print(string(byteRole))
}
