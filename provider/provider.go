package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dockerClient "github.com/tunein/terraform-provider-docker/provider/docker-client"
	dockerRegistry "github.com/tunein/terraform-provider-docker/provider/docker-registry"
	"github.com/tunein/terraform-provider-docker/provider/helper"
)

type DockerProvider struct {
	dockerClient   *dockerClient.SdkClient
	registryClient *dockerRegistry.HttpClient
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"docker_hub_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"docker_hub_password": {
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

	client, err := dockerClient.NewClient(authStr)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	provider.dockerClient = client

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
