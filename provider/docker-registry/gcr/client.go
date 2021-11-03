package gcr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type RegistryClient struct {
	client *http.Client
}

func NewGCRClient() *RegistryClient {
	return &RegistryClient{}
}

func (c RegistryClient) Login() error {
	return nil
}

func (c RegistryClient) IfImageExist(repo, tag string) error {
	params := strings.Split(repo, "/")
	url := fmt.Sprintf("https://k8s.gcr.io/v2/%s/%s/manifests/%s", params[1], params[2], tag)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Image was not found: " + repo + ":" + tag)
	}
	return nil
}
