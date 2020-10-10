---
page_title: "rdsdataservice_postgres_role"
---

# rdsdataservice_postgres_role Resource

Manage postgres roles

## Example Usage

```hcl
resource "rdsdataservice_postgres_role" "test" {
  name         = "test"
  resource_arn = var.db_arn
  secret_arn   = var.secret_arn
  login        = true
}
```

## Argument Reference

- `name` - (Required) The PostgreSQL database name to connect to.
- `resource_arn` - (Required) DB ARN.
- `secret_arn` - (Required) DBA Secret ARN.
- `login` - (Optional) Determine whether a role is allowed to log in. (Default: `false`)
- `inherit` - (Optional) Determine whether a role "inherits" the privileges of roles it is a member of. (Default: `true`)
- `create_database` - (Optional) Define a role's ability to create databases. (Default: `false`)
- `create_role` - (Optional) Determine whether this role will be permitted to create new roles. (Default: `false`)
- `password` - (Optional) Set the role's password.
- `roles` - (Optional) Role(s) to grant to this new role.
- `superuser` - (Optional) Determine whether the new role is a "superuser". (Default: `false`)

## Attribute Reference
