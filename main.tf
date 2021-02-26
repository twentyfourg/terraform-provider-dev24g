terraform {
  required_providers {
    dev24g = {
      source = "twentyfourg/dev24g"
    }
  }
  required_version = "~>0.14.3"
}

# Configure the Provider
provider "dev24g" {
  workspace = "24g"
}

data "dev24g_bitbucket_repository" "api" {
  name      = "2222-13-billing-tool"
}

resource "dev24g_bitbucket_deployment" "evan" {
  name       = "dev-api"
  stage      = "Test"
  repository = "${data.dev24g_bitbucket_repository.api.workspace}/${data.dev24g_bitbucket_repository.api.name}"
}

resource "dev24g_bitbucket_deployment_variable" "foobar" {
  key        = "foo"
  value      = "bar"
  secured    = false
  deployment = dev24g_bitbucket_deployment.evan.id
}