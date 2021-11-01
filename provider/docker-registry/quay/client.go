package quay

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type QuayClient struct {
}

func NewQuayClient() *QuayClient {
	return &QuayClient{}
}

func (c QuayClient) Login() error {
	return nil
}

func (c QuayClient) IfImageExist(repo, tag string) error {
	params := strings.Split(repo, "/")
	url := fmt.Sprintf("https://quay.io/v2/%s/%s/manifests/%s", params[1], params[2], tag)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Image was not found: " + repo + ":" + tag)
	}
	return nil
}