provider "dev24g" {}

data "dev24g_bitbucket_repository" "api" {
  name = "796-4-1-vxp-api"
  workspace = "24g"
}

output "repo" {
  value = data.dev24g_bitbucket_repository.api.name
}
