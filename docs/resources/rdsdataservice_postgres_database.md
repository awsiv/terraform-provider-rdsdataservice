---
page_title: "rdsdataservice_postgres_database"
---

# rdsdataservice_postgres_database Resource

Manage postgres databases

## Example Usage

```hcl
resource "rdsdataservice_postgres_database" "test" {
  name         = "test"
  resource_arn = var.db_arn
  secret_arn   = var.secret_arn
  owner        = "postgres"
}
```

## Argument Reference

- `name` - (Required) The PostgreSQL database name.
- `resource_arn` - (Required) DB ARN.
- `secret_arn` - (Required) DBA Secret ARN.
- `owner` - (Optional) The ROLE which owns the database.. (Default: `postgres`)

## Attribute Reference
