package docker_client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"log"
	"strings"
)

type SdkClient struct {
	context context.Context
	client  *client.Client
	authStr string
}

func NewClient(authStr string) (dockerClient *SdkClient, err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	dockerClient = &SdkClient{
		client:  cli,
		context: ctx,
		authStr: authStr,
	}
	return
}

func (c *SdkClient) ImagePull(repo, tag string) (err error) {
	options := types.ImagePullOptions{}

	//if registry is AWS ECR - add authentication string to options
	if strings.Contains(repo, ".dkr.ecr.") {
		options.RegistryAuth = c.authStr
	}

	reader, err := c.client.ImagePull(c.context, fmt.Sprintf("%s:%s", repo, tag), options)
	if err != nil {
		return
	}
	defer reader.Close()
	if logs, err := io.ReadAll(reader); err == nil {
		log.Print(string(logs))
	}
	return
}

func (c SdkClient) ImageTag(sourceRepo, targetRepo, tag string) (err error) {
	err = c.client.ImageTag(c.context, fmt.Sprintf("%s:%s", sourceRepo, tag), fmt.Sprintf("%s:%s", targetRepo, tag))
	return
}

func (c SdkClient) ImagePush(repo, tag string) (err error) {
	options := types.ImagePushOptions{
		RegistryAuth: c.authStr,
	}
	reader, err := c.client.ImagePush(c.context, fmt.Sprintf("%s:%s", repo, tag), options)
	if err != nil {
		return
	}
	type ErrorMessage struct {
		Error string
	}
	var errorMessage ErrorMessage
	buffIOReader := bufio.NewReader(reader)
	for {
		streamBytes, err := buffIOReader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err := json.Unmarshal(streamBytes, &errorMessage); err != nil {
			return err
		}
		if errorMessage.Error != "" {
			return fmt.Errorf("error pushing image: %s", errorMessage.Error)
		}
	}
	return
}
