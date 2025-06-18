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

### **Resource Documentation**
Each resource should have its own Markdown file in `docs/resources/<resource_name>.md`. The file should include:

1. **Title**: The name of the resource.
2. **Description**: A brief overview of what the resource does.
3. **Example Usage**: A code snippet showing how to use the resource.
4. **Argument Reference**: A list of all arguments for the resource, including whether they are required or optional.
5. **Attributes Reference**: Any attributes exported by the resource.

Example:
```markdown
# vaultassist_bootstrap_secret Resource

The `vaultassist_bootstrap_secret` resource is used to bootstrap a secret in Vault.

## Example Usage

```hcl
resource "vaultassist_bootstrap_secret" "example" {
  path  = ""
  mount = ""
}