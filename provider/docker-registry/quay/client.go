package quay

import (
	"fmt"
	"net/http"
	"strings"
)

type RegistryClient struct {
}

func NewQuayClient() *RegistryClient {
	return &RegistryClient{}
}

func (c RegistryClient) Login() error {
	return nil
}

func (c RegistryClient) DoesImageExist(repo, tag string) (bool, error) {
	params := strings.Split(repo, "/")
	url := fmt.Sprintf("https://quay.io/v2/%s/%s/manifests/%s", params[1], params[2], tag)
	response, err := http.Get(url)
	if err != nil {
		return false, err
	}
	if response.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}
