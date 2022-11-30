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
  name = "test-2"
  cpus = 1
  ram_mib = 1024
  ssh_keys = [{body= ""
               name="alkalsjd"}]
  state = "running"

  disks = [ {position = 0
             size_mb = 10000} ]

}

output "server" {
  value = resource.cloudafrica_server.example
}

