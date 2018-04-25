package iam

import (
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-lambda/resource"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type AssumeRolePolicyDocumentStatement struct {
	Effect    string `yaml:"Effect"`
	Principal struct {
		Service string `yaml:"Service"`
	} `yaml:"Principal"`
	Action []string `yaml:"Action"`
}

type AssumeRolePolicyDocument struct {
	Version   string                               `yaml:"Version"`
	Statement []*AssumeRolePolicyDocumentStatement `yaml:"Statement"`
}

type PolicyDocumentStatement struct {
	Effect   string   `yaml:"Effect"`
	Action   []string `yaml:"Action"`
	Resource string   `yaml:"Resource"`
}

type Policy struct {
	PolicyName     string `yaml:"PolicyName"`
	PolicyDocument struct {
		Version   string                     `yaml:"Version"`
		Statement []*PolicyDocumentStatement `yaml:"Statement"`
	} `yaml:"PolicyDocument"`
}

type Role struct {
	Type       string `yaml:"Type"`
	Properties struct {
		AssumeRolePolicyDocument AssumeRolePolicyDocument `yaml:"AssumeRolePolicyDocument"`
		Path                     string                   `yaml:"Path"`
		Policies                 []*Policy                `yaml:"Policies"`
		RoleName                 string                   `yaml:"RoleName"`
	} `yaml:"Properties"`
}

func (m *Role) Bytes() (byteRole []byte) {
	byteRole, err := yaml.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal iam role failed. \n%s.", err)
		return
	}
	return
}

func (m *Role) String() (strRole string) {
	strRole = string(m.Bytes())
	return
}

func (m *Role) WriteTo(filePath string, perm os.FileMode) (err error) {
	err = ioutil.WriteFile(filePath, m.Bytes(), perm)
	if nil != err {
		logrus.Errorf("write role to %q failed. \n%s.", filePath, err)
		return
	}
	return
}

func (m *Role) AddAssumeRolePolicyDocumentStatement(statement *AssumeRolePolicyDocumentStatement) {
	statements := m.Properties.AssumeRolePolicyDocument.Statement
	m.Properties.AssumeRolePolicyDocument.Statement = append(statements, statement)
}

func (m *Role) AddPolicy(policy *Policy) {
	m.Properties.Policies = append(m.Properties.Policies, policy)
}

func (m *Role) AddDynamoDbPolicy() {
	policy := &Policy{
		PolicyName: "DynamoDBFullAccess",
	}
	policy.PolicyDocument.Version = "2012-10-17"
	statement := &policy.PolicyDocument.Statement
	*statement = append(*statement, &PolicyDocumentStatement{
		Effect:   "Allow",
		Resource: "*",
		Action: []string{
			"dynamodb:*",
		},
	})

	m.AddPolicy(policy)
}

func NewExecutionRole(roleName string) (role *Role) {
	role = &Role{}
	role.Properties.RoleName = roleName
	role.Type = "AWS::IAM::Role"
	properties := &role.Properties
	properties.Path = "/"
	properties.AssumeRolePolicyDocument.Version = "2012-10-17"

	// 默认添加lambda信任关系
	var lambdaRolePolicyStatement AssumeRolePolicyDocumentStatement
	lambdaRolePolicyStatement.Effect = "Allow"
	lambdaRolePolicyStatement.Principal.Service = "lambda.amazonaws.com"
	lambdaRolePolicyStatement.Action = append(lambdaRolePolicyStatement.Action, "sts:AssumeRole")
	role.AddAssumeRolePolicyDocumentStatement(&lambdaRolePolicyStatement)

	// 默认添加CloudWatch policy
	var policy Policy
	policy.PolicyName = "CloudWatchLogs"
	policy.PolicyDocument.Version = "2012-10-17"
	policyStatement := policy.PolicyDocument.Statement
	policy.PolicyDocument.Statement = append(policyStatement, &PolicyDocumentStatement{
		Effect: "Allow",
		Action: []string{
			"logs:CreateLogStream",
			"logs:CreateLogGroup",
			"logs:PutLogEvents",
		},
		Resource: "arn:aws:logs:*:*:*",
	})

	role.AddPolicy(&policy)

	resourceInUse := resource.GetResourceInUse()
	for _, resourceType := range resourceInUse {
		switch resourceType {
		case resource.DynamoDBResourceType:
			role.AddDynamoDbPolicy()
		}
	}

	return
}
