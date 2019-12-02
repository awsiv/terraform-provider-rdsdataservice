# terraform-provider-rdsdataservice
Manage Postgres db resources using the AWS Data API - Heavily inspired by [terraform-provider-postgresql] (https://github.com/terraform-providers/terraform-provider-postgresql)

## Requirements ##
Terraform 0.12+
Go 1.13 (to build the provider plugin)

## Install ##

You will need to install the binary as a [terraform third party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).  Terraform will then pick up the binary from the local filesystem when you run `terraform init`.

```sh
curl -s https://raw.githubusercontent.com/awsiv/terraform-provider-rdsdataservice/master/install.sh | bash
```

## Usage ##

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