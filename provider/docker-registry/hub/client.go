package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type DockerHubLoginResponce struct {
	Token string `json:"token"`
}

type DockerHubClient struct {
	client     *http.Client
	v1endpoint string
	v2endpoint string
	username   string
	password   string
	token      string
}

func NewDockerHubClient(username, password string) *DockerHubClient {
	client := &http.Client{}
	return &DockerHubClient{
		client:     client,
		v1endpoint: "https://index.docker.io/v1",
		v2endpoint: "https://hub.docker.com",
		username:   username,
		password:   password,
	}
}

func (c DockerHubClient) Login() error {
	params := url.Values{}
	params.Set("username", c.username)
	params.Set("password", c.password)
	request, err := http.NewRequest("POST", c.v2endpoint+"/v2/users/login", strings.NewReader(params.Encode()))
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
	var result DockerHubLoginResponce
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	c.token = result.Token
	return nil
}

func (c DockerHubClient) IfImageExist(repo, tag string) error {
	params := strings.Split(repo, "/")
	request, err := http.NewRequest("GET", fmt.Sprintf("https://index.docker.io/v1/repositories/%s/%s/tags/%s", params[1], params[2], tag), nil)
	if err != nil {
		panic(err)
	}
	response, err := c.client.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		return errors.New("Image was not found: " + repo + ":" + tag)
	}
	return nil
}
