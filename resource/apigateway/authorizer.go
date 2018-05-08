package apigateway

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/haozzzzzzzz/go-lambda/resource/iam"
)

func GetAllowAuthorizerResponse(principalId string, arn string) (response *events.APIGatewayCustomAuthorizerResponse) {
	return GetAuthorizerResponse(principalId, iam.EffectAllow, arn)
}

func GetDenyAuthorizerResponse(principalId string, arn string) (response *events.APIGatewayCustomAuthorizerResponse) {
	return GetAuthorizerResponse(principalId, iam.EffectDeny, arn)
}

func GetAuthorizerResponse(principalId string, effect string, arn string) (response *events.APIGatewayCustomAuthorizerResponse) {
	response = &events.APIGatewayCustomAuthorizerResponse{}
	response.PrincipalID = principalId
	policy := &response.PolicyDocument
	policy.Version = "2012-10-17"
	statement := events.IAMPolicyStatement{}
	statement.Effect = effect
	statement.Resource = append(statement.Resource, arn)
	statement.Action = append(statement.Action, "execute-api:Invoke")
	policy.Statement = append(policy.Statement, statement)
	return
}
