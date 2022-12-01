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

resource "cloudafrica_server" "example" {
  //id = 17592186045491
}

output "server" {
  value = resource.cloudafrica_server.example
}

