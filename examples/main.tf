terraform {
  required_providers {
    docker = {
      version = "0.1.0"
      source = "tunein.com/devops/docker"
    }
  }
}

provider "docker" {
  docker_hub_username = var.docker_hub_username
  docker_hub_password = var.docker_hub_password
}

data "docker_upstream_image" "newrelic_k8s_metadata_injection" {
  repo = "docker.io/newrelic/k8s-metadata-injection"
  tag = "1.6.0"
}

resource "docker_downstream_image" "newrelic_k8s_metadata_injection" {
  upstream_repo = data.docker_upstream_image.newrelic_k8s_metadata_injection.repo
  downstream_repo = "809163537096.dkr.ecr.us-west-2.amazonaws.com/docker.io/newrelic/k8s-metadata-injection"
  tag = data.docker_upstream_image.newrelic_k8s_metadata_injection.tag
}