package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tunein/terraform-provider-docker/provider/helper"
)

type DockerProvider struct {
	dockerHubUsername string
	dockerHubPassword string
	dockerClient      *helper.DockerClient
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"docker_hub_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"docker_hub_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"docker_downstream_image": ResourceDownstreamImage(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"docker_upstream_image": DataUpstreamImage(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	provider := &DockerProvider{}

	username := d.Get("docker_hub_username").(string)
	password := d.Get("docker_hub_password").(string)

	if (username != "") && (password != "") {
		provider.dockerHubUsername = username
		provider.dockerHubPassword = password
	}

	awsClient := helper.NewAwsClient()
	authStr, err := awsClient.GetDockerAuthStrFromEcr()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	dockerClient, err := helper.NewDockerClient(authStr)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	provider.dockerClient = dockerClient

	return provider, diags
}
