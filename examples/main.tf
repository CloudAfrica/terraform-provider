terraform {
  required_providers {
    cloudafrica = {
      source = "registry.terraform.io/cloudafrica/cloudafrica"
    }
  }
}

provider "cloudafrica" {
  host = "http://localhost:3000/api"
}

data "cloudafrica_servers" "example" {}

output "servers" {
  value = data.cloudafrica_servers.example
}

