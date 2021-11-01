terraform {
  required_providers {
    docker = {
      version = "0.1.0"
      source  = "tunein.com/devops/docker"
    }
  }
}

provider "docker" {
  docker_hub_username = "login"
  docker_hub_password = "password"
}