package helper

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"strings"
)

type AwsClient struct {
	client client.ConfigProvider
	ecr    *ecr.ECR
}

func NewAwsClient() (awsClient *AwsClient) {
	awsClient = &AwsClient{}

	sessOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	awsConfig := &aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}

	sess := session.Must(session.NewSessionWithOptions(sessOptions))

	awsClient.ecr = ecr.New(sess, awsConfig)

	return awsClient
}

func (c *AwsClient) GetDockerAuthStrFromEcr() (string, error) {
	tokenOutput, err := c.ecr.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}

	decodedToken, err := base64.StdEncoding.DecodeString(aws.StringValue(tokenOutput.AuthorizationData[0].AuthorizationToken))
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(string(decodedToken), ":", 2)

	authConfig := types.AuthConfig{
		Username:      parts[0],
		Password:      parts[1],
		ServerAddress: aws.StringValue(tokenOutput.AuthorizationData[0].ProxyEndpoint),
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	return authStr, nil

}