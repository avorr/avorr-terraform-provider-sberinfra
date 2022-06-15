package imports

const (
	diProviderTpl = `di = {
      source  = "sberbank/devops/di"
      version = "0.2.15"
    }`
	vaultProviderTpl = `  vault = {
      source  = "hashicorp/vault"
      version = "2.17.0"
    }`
	remoteStateTpl = `
data "terraform_remote_state" "project" {
#  backend = "local"
#  config = {
#    path = "${path.module}/../project/terraform.tfstate"
#  }
#  backend = "remote"
  backend = "s3"
  config = {
    bucket = "%s"
    key    = "%s/project.json"
    skip_credentials_validation = true
    skip_region_validation = true
    force_path_style = true
  }
}`
	s3backendTpl = `backend "s3" {
    bucket = "%s"
    key    = "%s/%s.json"
    skip_credentials_validation = true
    skip_region_validation = true
    force_path_style = true
  }`
	mainTpl = `#export DI_URL=https://cs.cloud.sberbank.ru/api/v1
#export DI_TOKEN=supersecrettoken
#export DI_DEBUG=0
#export AWS_S3_ENDPOINT=http://miniohost:9000
#export AWS_DEFAULT_REGION=no-region
#export AWS_ACCESS_KEY_ID=
#export AWS_SECRET_ACCESS_KEY=
#export VAULT_ADDR=http://vaulthost:8200
#export VAULT_TOKEN=

terraform {
  required_version = ">= 0.14.4"
  required_providers {
    %s
	%s
  }
  %s
}
provider "di" {}

%s`
)
