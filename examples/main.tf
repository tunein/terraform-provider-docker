terraform {
  required_providers {
    docker = {
      version = "0.1.0"
      source  = "tunein.com/devops/docker"
    }
  }
}
#
provider "docker" {
  docker_hub_username = var.docker_hub_username
  docker_hub_password = var.docker_hub_password
}

data "docker_upstream_image" "alpine" {
  repo = "alpine"
  tag  = "3.12.8"
}

data "docker_upstream_image" "kube_state_metric" {
  repo = "quay.io/coreos/kube-state-metrics"
  tag  = "v1.9.6"
}

data "docker_upstream_image" "newrelic_k8s_metadata_injection" {
  repo = "docker.io/newrelic/k8s-metadata-injection"
  tag  = "1.6.0"
}

data "docker_upstream_image" "metric_server" {
  repo = "k8s.gcr.io/metrics-server/metrics-server"
  tag  = "v0.5.0"
}

locals {
  alpine_downstream_repo                          = "${var.aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${data.docker_upstream_image.alpine.repo}"
  kube_state_metric_downstream_repo               = "${var.aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${data.docker_upstream_image.kube_state_metric.repo}"
  newrelic_k8s_metadata_injection_downstream_repo = "${var.aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${data.docker_upstream_image.newrelic_k8s_metadata_injection.repo}"
  metric_server_downstream_repo                   = "${var.aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${data.docker_upstream_image.metric_server.repo}"
}

resource "docker_downstream_image" "alpine" {
  upstream_repo   = data.docker_upstream_image.alpine.repo
  downstream_repo = local.alpine_downstream_repo
  tag             = data.docker_upstream_image.alpine.tag
}

resource "docker_downstream_image" "kube_state_metric" {
  upstream_repo   = data.docker_upstream_image.kube_state_metric.repo
  downstream_repo = local.kube_state_metric_downstream_repo
  tag             = data.docker_upstream_image.kube_state_metric.tag
}

resource "docker_downstream_image" "newrelic_k8s_metadata_injection" {
  upstream_repo   = data.docker_upstream_image.newrelic_k8s_metadata_injection.repo
  downstream_repo = local.newrelic_k8s_metadata_injection_downstream_repo
  tag             = data.docker_upstream_image.newrelic_k8s_metadata_injection.tag
}

resource "docker_downstream_image" "metric_server" {
  upstream_repo   = data.docker_upstream_image.metric_server.repo
  downstream_repo = local.metric_server_downstream_repo
  tag             = data.docker_upstream_image.metric_server.tag
}