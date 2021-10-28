package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDownstreamImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDownstreamImageCreate,
		ReadContext:   resourceDownstreamImageRead,
		UpdateContext: resourceDownstreamImageUpdate,
		DeleteContext: resourceDownstreamImageDelete,
		Schema: map[string]*schema.Schema{
			"upstream_repo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"downstream_repo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDownstreamImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*DockerProvider)

	var upstreamRepo = d.Get("upstream_repo").(string)
	var downstreamRepo = d.Get("downstream_repo").(string)
	var tag = d.Get("tag").(string)

	err := provider.dockerClient.ImagePull(upstreamRepo, tag)
	if err != nil {
		return diag.FromErr(err)
	}

	err = provider.dockerClient.ImageTag(upstreamRepo, downstreamRepo, tag)
	if err != nil {
		return diag.FromErr(err)
	}

	err = provider.dockerClient.ImagePush(downstreamRepo, tag)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("1")

	return nil
}

func resourceDownstreamImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceDownstreamImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceDownstreamImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
