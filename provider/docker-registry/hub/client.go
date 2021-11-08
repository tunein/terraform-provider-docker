package hub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type DockerHubLoginResponse struct {
	Token string `json:"token"`
}

type RegistryClient struct {
	client     *http.Client
	v1endpoint string
	v2endpoint string
	username   string
	password   string
	token      string
}

func NewDockerHubClient(username, password string) *RegistryClient {
	client := &http.Client{}
	return &RegistryClient{
		client:     client,
		v1endpoint: "https://index.docker.io/v1",
		v2endpoint: "https://hub.docker.com/v2",
		username:   username,
		password:   password,
	}
}

func (c RegistryClient) Login() error {
	params := url.Values{}
	params.Set("username", c.username)
	params.Set("password", c.password)
	request, err := http.NewRequest("POST", c.v2endpoint+"/users/login", strings.NewReader(params.Encode()))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var result DockerHubLoginResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		return fmt.Errorf("can not unmarshal JSON")
	}
	c.token = result.Token
	return nil
}

func (c RegistryClient) DoesImageExist(repo, tag string) (bool, error) {
	repo = strings.ReplaceAll(repo, "docker.io/", "")
	request, err := http.NewRequest("GET", fmt.Sprintf("https://index.docker.io/v1/repositories/%s/tags/%s", repo, tag), nil)
	if err != nil {
		return false, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}
