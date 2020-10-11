---
page_title: "Provider: RDS DataService - DataAPI"
---

# RDSDataService Provider

Manage Aurora Serverless databases with Terraform.

[AWS RDSDataService/Data API](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/data-api.html) allows us to run SQL using HTTP endpoint and AWS SDK.

Due to this, we have the following advantages:

- We no longer need to manage connections
- We can use secretsmanager secret and not have to worry about secrets ending up in terraform state.

## Example Usage

```hcl
provider "rdsdataservice" {
    version = "1.0.2"
    region  = var.aws_region
    profile = var.aws_profile
}

resource "rdsdataservice_postgres_database" "test" {
    name         = "test"
    resource_arn = var.db_arn
    secret_arn   = var.secret_arn
    owner        = "postgres"
}

resource "rdsdataservice_postgres_role" "test" {
    name         = "test"
    resource_arn = var.db_arn
    secret_arn   = var.secret_arn
    login        = true
}
```

## Argument Reference

This provider is built to be compatible/similar to [terraform-provider-aws](https://registry.terraform.io/providers/hashicorp/aws/latest/docs), since it uses the AWS SDK and the provider implemenation is inspired by it.
