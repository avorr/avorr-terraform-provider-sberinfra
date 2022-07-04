#export DI_URL=https://cs.cloud.sberbank.ru/api/v1
#https://iftpdi.ift.esrt.cloud.sbrf.ru/client/orders/stands
#export DI_TOKEN=supersecrettoken
#export DI_DEBUG=0
#export AWS_S3_ENDPOINT=http://miniohost:9000
#export AWS_DEFAULT_REGION=no-region
#export AWS_ACCESS_KEY_ID=
#export AWS_SECRET_ACCESS_KEY=
#export VAULT_ADDR=http://vaulthost:8200
#export VAULT_TOKEN=

terraform {
  required_version = ">= 0.14.9"
  required_providers {
    di = {
      source  = "sberbank/devops/di"
      version = "0.3.13"
    }
#    vault = {
#      version = "3.3.0"
#    }
  }
#  backend "s3" {
#    key                         = "sberdevices/smartapp/asdsadsads.tfstate"
#    bucket                      = "tfstates"
#    skip_credentials_validation = true
#    skip_region_validation      = true
#    skip_metadata_api_check     = true
#    force_path_style            = true
#  }
}
provider "di" {}

#provider "vault" {
#  skip_child_token = true
#  skip_tls_verify = true
##	address = "http://tkles-pcb000207.vm.esrt.cloud.sbrf.ru:8200"
#	address = "https://ift.secrets.sigma.sbrf.ru" # VAULT_ADDR
#	auth_login {
##		namespace = "CI01808661_CI01875672"
#		namespace  = "CI01808661_CI02063031_DEV"
#		path       = "auth/approle/login"
#		parameters = {
#			role_id = "${var.role_id}"
#			secret_id = "${var.secret_id}"
#		}
#	}
#}
#
#variable "role_id" {}
#variable "secret_id" {}