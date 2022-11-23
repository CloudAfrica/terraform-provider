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

data "cloudafrica_billing_accounts" "example" {}

output "bas" {
  value = data.cloudafrica_billing_accounts.example
}

