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
  name = "andre-test-name"
  cpus = 1
  ram_mib = 1024
  ssh_keys = [{body= "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC1zUv7uXSKY0WRWHQoO2jwT+dllMw1HEpKy7lSG8ve57X+IiJlD6qdHfIYce2RK5a7VVkM5Mi9ddHYD8iv7beXoDGLu1aa6WnTAOiOPHR24RgtkZVur8wsGUcE0JKTWUBmvTHG+k5x1UYG9NU7e8QooiprWgvolQleXVkfXGlse0+DY0LYnmhDFUpj2cnOLyR30z4IRvRC7Plx+cSEaY8r/WQAhC1soOrKeMLgkrXePNoqjTes6uLTd4RE5xXYsIA7Fxp9X8a6CNbV3AV2fpfmqbiBqD+bClr7i219KOOQBb6BZEDjfNVlKEGzxVsZ+ZVNj/bifKZtV9RHZufTTelR user@host"
               name="user-key"}]
  state = "running"

  disks = [ {position = 0
             size_mb = 10000} ]

}

output "server" {
  value = resource.cloudafrica_server.example
}

