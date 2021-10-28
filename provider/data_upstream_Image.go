package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tunein/terraform-provider-docker/provider/helper"
	"strconv"
	"time"
)

type Image struct {
	Repo string
	Tag  string
	Sha  string
}

func DataUpstreamImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,
		Schema: map[string]*schema.Schema{
			"repo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "latest",
			},
		},
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*DockerProvider)
	var diags diag.Diagnostics

	var image = &Image{}

	awsClient := helper.NewAwsClient()
	authStr, err := awsClient.GetDockerAuthStrFromEcr()
	if err != nil {
		diag.FromErr(err)
	}
	registryHttpClient := helper.NewRegistryHttpClient(authStr, provider.dockerHubUsername, provider.dockerHubPassword)

	repo, ok := d.GetOk("repo")
	if ok {
		image.Repo = repo.(string)
	}

	tag, ok := d.GetOk("tag")
	if ok {
		image.Tag = tag.(string)
	}

	err = registryHttpClient.IfImageExist(image.Repo, image.Tag)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("repo", image.Repo)
	d.Set("tag", image.Tag)

	return diags
}
