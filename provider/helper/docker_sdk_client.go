package helper

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"strings"
)

type DockerClient struct {
	context context.Context
	client  *client.Client
	authStr string
}

func NewDockerClient(authStr string) (dockerClient *DockerClient, err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	dockerClient = &DockerClient{
		client:  cli,
		context: ctx,
		authStr: authStr,
	}
	return
}

func (d *DockerClient) ImageSearch(repo, tag string) (err error) {
	options := types.ImageSearchOptions{
		RegistryAuth: d.authStr,
		Limit:        10,
	}
	result, err := d.client.ImageSearch(d.context, repo, options)
	if err != nil {
		return
	}
	fmt.Println(result)
	return
}

func (d *DockerClient) ImagePull(repo, tag string) (err error) {
	options := types.ImagePullOptions{}

	//if registry is AWS ECR - add authentication string to options
	if strings.Contains(repo, ".dkr.ecr.") {
		options.RegistryAuth = d.authStr
	}

	reader, err := d.client.ImagePull(d.context, fmt.Sprintf("%s:%s", repo, tag), options)
	if err != nil {
		return
	}
	defer reader.Close()
	return
}

func (d DockerClient) ImageTag(sourceRepo, targetRepo, tag string) (err error) {
	err = d.client.ImageTag(d.context, fmt.Sprintf("%s:%s", sourceRepo, tag), fmt.Sprintf("%s:%s", targetRepo, tag))
	return
}

func (d DockerClient) ImagePush(repo, tag string) (err error) {
	options := types.ImagePushOptions{
		RegistryAuth: d.authStr,
	}
	reader, err := d.client.ImagePush(d.context, fmt.Sprintf("%s:%s", repo, tag), options)
	if err != nil {
		return
	}
	defer reader.Close()
	var out string
	if b, err := io.ReadAll(reader); err == nil {
		out = string(b)
	}
	fmt.Println(out)
	return
}
