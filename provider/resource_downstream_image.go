package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func ResourceDownstreamImage() *schema.Resource {
	return &schema.Resource{
		Description:   "Pulls image from a public registry and pushes it to the private ECR",
		CreateContext: resourceDownstreamImageCreate,
		ReadContext:   resourceDownstreamImageRead,
		UpdateContext: resourceDownstreamImageUpdate,
		DeleteContext: resourceDownstreamImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"upstream_repo": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "full public repository name",
			},
			"downstream_repo": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "full ECR repository name",
			},
			"tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "image tag",
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

	d.SetId(downstreamRepo + ":" + tag)

	return nil
}

func resourceDownstreamImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var upstreamRepo = d.Get("upstream_repo").(string)
	var downstreamRepo = d.Get("downstream_repo").(string)
	var tag = d.Get("tag").(string)
	provider := m.(*DockerProvider)

	if d.State().ID != "" {
		image := strings.SplitN(d.State().ID, ":", 2)
		if len(image) != 2 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Wrong import ID",
				Detail:   "Import image should have tag after \":\"",
			})
			return diags
		}
		repo := strings.SplitN(image[0], "/", 2)
		if len(repo) != 2 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Wrong import ID",
				Detail:   "Import image should have valid ECR definition",
			})
			return diags
		}
		upstreamRepo = repo[1]
		downstreamRepo = image[0]
		tag = image[1]
	}

	imageExist, err := provider.registryClient.DoesImageExist(downstreamRepo, tag)
	if err != nil {
		return diag.FromErr(err)
	}
	if !imageExist {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Image do not exist",
			Detail:   "The image with such tag was not found in the registry.",
		})
		return diags
	}

	err = d.Set("upstream_repo", upstreamRepo)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("downstream_repo", downstreamRepo)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("tag", tag)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDownstreamImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDownstreamImageCreate(ctx, d, m)
}

func resourceDownstreamImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
