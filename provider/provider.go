package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dockerRegistry "github.com/tunein/terraform-provider-docker/provider/docker-registry"
	"github.com/tunein/terraform-provider-docker/provider/helper"
)

type DockerProvider struct {
	dockerClient   *helper.DockerClient
	registryClient *dockerRegistry.RegistryHttpClient
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

	awsClient := helper.NewAwsClient()

	// AUTH for docker client
	authStr, err := awsClient.GetDockerAuthStrFromEcr()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	dockerClient, err := helper.NewDockerClient(authStr)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	provider.dockerClient = dockerClient

	// AUTH for registry
	authStr, err = awsClient.GetAuthStrFromEcr()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	username := d.Get("docker_hub_username").(string)
	password := d.Get("docker_hub_password").(string)

	registryClient := dockerRegistry.NewRegistryHttpClient(authStr, username, password)
	provider.registryClient = registryClient

	return provider, diags
}
