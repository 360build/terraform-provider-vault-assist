terraform {
  required_providers {
    vaultassist = {
      source  = "360build/vaultassist"
    }
  }
}

provider "vaultassist" {
  address = ""
  role = ""
  mountpoint = ""
  headers = {
    "cf-access-token" = var.cftoken
  }
}

resource "vaultassist_bootstrap_secret" "test" {
  path = ""
}


