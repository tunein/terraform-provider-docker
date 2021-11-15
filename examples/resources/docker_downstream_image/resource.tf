resource "docker_downstream_image" "alpine" {
  upstream_repo   = "alpine"
  downstream_repo = "1234567890.dkr.ecr.us-west-2.amazonaws.com/alpine"
  tag             = "latest"
}