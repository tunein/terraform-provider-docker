package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

// AUTH str for ECR
func (c *AwsClient) GetAuthStrFromEcr() (string, error) {
	tokenOutput, err := c.ecr.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}
	authStr := aws.StringValue(tokenOutput.AuthorizationData[0].AuthorizationToken)
	return authStr, nil
}

// AUTH str for Docker
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
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid auhtorization token. Report an issue to this docker provider")
	}

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
