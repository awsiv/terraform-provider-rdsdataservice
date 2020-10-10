# terraform-provider-rdsdataservice

Manage Postgres db resources using the AWS Data API - Heavily inspired by [terraform-provider-postgresql](https://github.com/terraform-providers/terraform-provider-postgresql)

[AWS Data API](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/data-api.html) allows us to run SQL using HTTP endpoint and AWS SDK. This is awesome because it means that we no longer need to manage connections. This also uses secretsmanager secret so we no longer have to worry about secrets ending up in terraform state.

Since it uses AWS SDK, it might as well belong to terraform-provider-aws itself, but then, the CRUD operations are SQL statements instead of actual API calls - so maybe it has its own place? I am working on porting more resources and more importantly the acceptance tests. Let me know what you think about it :)

API documentation: [package rdsdataservice](https://godoc.org/github.com/aws/aws-sdk-go/service/rdsdataservice)

## Requirements

Terraform 0.12+
Go 1.13 (to build the provider plugin)

## Install

You will need to install the binary as a [terraform third party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins). Terraform will then pick up the binary from the local filesystem when you run `terraform init`.

```sh
curl -s https://raw.githubusercontent.com/awsiv/terraform-provider-rdsdataservice/master/install.sh | bash
```

## Usage

```terraform
provider "rdsdataservice" {
  version = "1.0.0"
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
