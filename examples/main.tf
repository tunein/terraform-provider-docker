//TODO: rename to tunein-docker
terraform {
  required_providers {
    docker = {
      version = "0.1.0"
      source = "tunein.com/devops/docker"
    }
  }
}

//TODO: other docker provider
provider "docker" {
  docker_hub_username = var.docker_hub_username
  docker_hub_password = var.docker_hub_password
}

#OK
data "docker_upstream_image" "kube_state_metric" {
  repo = "quay.io/coreos/kube-state-metrics"
  tag = "v1.9.6"
}

#OK
data "docker_upstream_image" "newrelic_k8s_metadata_injection" {
  repo = "docker.io/newrelic/k8s-metadata-injection"
  tag = "1.6.0"
}

#OK
data "docker_upstream_image" "metric_server" {
  repo = "k8s.gcr.io/metrics-server/metrics-server"
  tag = "v0.5.0"
}

resource "docker_downstream_image" "kube_state_metric" {
  upstream_repo = data.docker_upstream_image.kube_state_metric.repo
  downstream_repo = "809163537096.dkr.ecr.us-west-2.amazonaws.com/quay.io/coreos/kube-state-metrics"
  tag = data.docker_upstream_image.kube_state_metric.tag
}

resource "docker_downstream_image" "newrelic_k8s_metadata_injection" {
  upstream_repo = data.docker_upstream_image.newrelic_k8s_metadata_injection.repo
  downstream_repo = "809163537096.dkr.ecr.us-west-2.amazonaws.com/docker.io/newrelic/k8s-metadata-injection"
  tag = data.docker_upstream_image.newrelic_k8s_metadata_injection.tag
}

resource "docker_downstream_image" "metric_server" {
  upstream_repo = data.docker_upstream_image.metric_server.repo
  downstream_repo = "809163537096.dkr.ecr.us-west-2.amazonaws.com/k8s.gcr.io/metrics-server/metrics-server"
  tag = data.docker_upstream_image.metric_server.tag
}