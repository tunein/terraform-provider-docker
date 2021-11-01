---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "docker Provider"
subcategory: ""
description: |-
  
---

# docker Provider



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **docker_hub_password** (String)
- **docker_hub_username** (String)