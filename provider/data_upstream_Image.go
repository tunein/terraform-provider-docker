package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Description: "Verifies the image presence in the repository",
		ReadContext: dataSourceImageRead,
		Schema: map[string]*schema.Schema{
			"repo": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "full repository name",
			},
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "latest",
				Description: "image tag",
			},
		},
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*DockerProvider)
	var diags diag.Diagnostics

	var image = &Image{}

	repo, ok := d.GetOk("repo")
	if ok {
		image.Repo = repo.(string)
	}

	tag, ok := d.GetOk("tag")
	if ok {
		image.Tag = tag.(string)
	}

	err := provider.registryClient.IfImageExist(image.Repo, image.Tag)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	err = d.Set("repo", image.Repo)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("tag", image.Tag)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
