provider "dev24g" {}

data "dev24g_bitbucket_repository" "api" {
  name      = "796-4-1-vxp-api"
  workspace = "24g"
}

resource "dev24g_bitbucket_deployment" "evan" {
  name       = "evan"
  stage      = "Test"
  repository = "${data.dev24g_bitbucket_repository.api.workspace}/${data.dev24g_bitbucket_repository.api.name}"
}

resource "dev24g_bitbucket_deployment_variable" "foobar" {
  key        = "foo"
  value      = "bar"
  secured    = false
  deployment = dev24g_bitbucket_deployment.evan.id
}

output "repo" {
  value = data.dev24g_bitbucket_repository.api.name
}
