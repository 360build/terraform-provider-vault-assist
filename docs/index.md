# VaultAssist Provider

The `vaultassist` provider is used to interact with Hasicorp Vault to manage secrets. It implents some missing functionallity 

## Example Usage

```hcl
provider "vaultassist" {
  address     = "https://<your vault ur>"
  role        = ""
  mountpoint  = ""
  headers = {
    "cf-access-token" = var.cftoken
  }
}
```
