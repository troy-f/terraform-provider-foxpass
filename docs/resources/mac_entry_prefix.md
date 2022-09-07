# MAC Entry Prefix (Resource)

```hcl
resource "foxpass_mac_entry_prefix" "resource-name" {
  entryname        = "example"
  prefix           = "aa:bb:cc:dd:ee:ff"
}
```

### Required

- `entryname` (String)
- `prefix` (String)

### Read-Only

- `id` (String) The ID of this resource.
