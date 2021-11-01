package ecr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type EcrClient struct {
	client  *http.Client
	authStr string
}

func NewEcrClient(authStr string) *EcrClient {
	client := &http.Client{}
	return &EcrClient{
		client:  client,
		authStr: authStr,
	}
}

func (c EcrClient) Login() error {
	return nil
}

func (c EcrClient) IfImageExist(repo, tag string) error {
	params := strings.SplitN(repo, "/", 2)
	request, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v2/%s/manifests/%s", params[0], params[1], tag), nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", "Basic "+c.authStr)
	response, err := c.client.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		return errors.New("Image was not found: " + repo + ":" + tag)
	}
	return nil
}
