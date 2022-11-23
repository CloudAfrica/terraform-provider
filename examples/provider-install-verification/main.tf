terraform {
  required_providers {
    cloudafrica = {
      source = "registry.terraform.io/cloudafrica/cloudafrica"
    }
  }
}

provider "cloudafrica" {}

data "cloudafrica_servers" "example" {}
