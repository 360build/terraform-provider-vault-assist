# vaultassist_bootstrap_secret Resource

The `vaultassist_bootstrap_secret` resource is used to create a secret trough terrafom on init but further never change update or do anything with it. It will just create an emphty secret. 

If the secret already exist it will leave it as is. This is used for when a new module need to have a place holder secret that should be created manually. 

## Example Usage

```hcl
resource "vaultassist_bootstrap_secret" "example" {
  path  = ""
  mount = ""
}
```

## Argument Reference

- `path` (String, **Required**) – The path to the secret within the mount.
- `mount` (String, **Required**) – The mount point in Vault.

