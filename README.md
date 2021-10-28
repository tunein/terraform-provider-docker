# terraform-provider-docker

The main idea of this provider is to copy a public image to private ECR registry. It will pull image from public 
registry to local machine, re-tag it and push to the ECR. 

Currently, it supports this public registries:
- ECR
- Docker Hub
- quay.io
- GCR

## How to build

```sh
$ make build
```

## How to use

```hcl
terraform {
  required_providers {
    docker = {
      version = "{{VERSION}}"
      source  = "tunein.com/devops/docker"
    }
  }
}

provider "docker" {
  docker_hub_username = var.docker_hub_username
  docker_hub_password = var.docker_hub_password
}

//Verify that required image exists during the `plan` 
data "docker_upstream_image" "alpine" {
  repo = "alpine"
  tag  = "3.12.8"
}

locals {
  alpine_downstream_repo = "${var.aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${data.docker_upstream_image.alpine.repo}"
}

resource "docker_downstream_image" "alpine" {
  upstream_repo   = data.docker_upstream_image.alpine.repo
  downstream_repo = local.alpine_downstream_repo
  tag             = data.docker_upstream_image.alpine.tag
}
```

## Development 

### Debug

Add `--debug` as a program argument to be able to debug the provider. 
