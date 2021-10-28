package ecr

import (
	"fmt"
	"net/http"
	"strings"
)

type RegistryClient struct {
	client  *http.Client
	authStr string
}

func NewEcrClient(authStr string) *RegistryClient {
	client := &http.Client{}
	return &RegistryClient{
		client:  client,
		authStr: authStr,
	}
}

func (c RegistryClient) Login() error {
	return nil
}

func (c RegistryClient) DoesImageExist(repo, tag string) (bool, error) {
	params := strings.SplitN(repo, "/", 2)
	if len(params) != 2 {
		return false, fmt.Errorf("image should have valid ECR definition")
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v2/%s/manifests/%s", params[0], params[1], tag), nil)
	if err != nil {
		return false, err
	}
	request.Header.Set("Authorization", "Basic "+c.authStr)
	response, err := c.client.Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}
