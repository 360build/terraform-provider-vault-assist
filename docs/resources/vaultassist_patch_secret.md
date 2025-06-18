# vaultassist_patch_secret Resource

The `vaultassist_patch_secret` resource is used to patch a key-value pair into a Vault secret at a specified path and mount.

## Example Usage

```hcl
resource "vaultassist_patch_secret" "test" {
  path  = ""
  mount = ""
  key   = ""
  value = jsonencode({
    "testkey" = "testvalue"
  })
}
```

## Argument Reference

- `path` (String, **Required**) – The path to the secret within the mount.
- `mount` (String, **Required**) – The mount point in Vault.
- `key` (String, **Required**) – The key to patch in the secret.
- `value` (String, **Required**, Sensitive) – The value to set for the key.

## Attributes Reference

- `created` – Boolean indicating if the secret was created or patched.