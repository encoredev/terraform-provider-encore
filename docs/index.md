---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "encore Provider"
subcategory: ""
description: |-
  
---

# encore Provider



## Example Usage

```terraform
provider "encore" {
  env      = "aws"
  auth_key = "auth_key"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `auth_key` (String) The [Encore Auth Key](https://encore.dev/docs/develop/auth-keys) to use to authenticate with the Encore Platform. Defaults to `ENCORE_AUTH_KEY` env var.
- `env` (String) The default Encore environment to operate on, if not overridden on a resource. Defaults to primary environment.
