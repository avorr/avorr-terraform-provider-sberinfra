#data "di_domain" "iac" {
#  name = "IaC"
#}
#
#data "di_group" "devices" {
#  name           = "Devices"
#  domain_id      = data.di_domain.iac.id
#}
#
#data "di_stand_type" "dev" {
#  name = "DEV"
#}
#
#data "di_as" "ec" {
#  code           = "CI02718748"
#  domain_id      = data.di_domain.iac.id
#}

data "di_domain" "domain" {
  name = "SberDevices"
}

data "di_group" "devices" {
  name           = "TestGroupSD"
  domain_id      = data.di_domain.domain.id
}

data "di_stand_type" "dev" {
  name = "DEV"
}

data "di_as" "ec" {
  code           = "CI01808661"
  domain_id      = data.di_domain.domain.id
}


#data "vault_generic_secret" "test_vault_secrets" {
##  path = "di-terraform/t-sberdevices/devops/dev/terraformtest/di_ignite/test-ignite-cluster"
##  path = "kv/team/test"
##  path = "CI01808661_CI01875672/A/MAIN/OSH/MAIN/KV/approle"
#  path = "A/MAIN/JEN/MAIN/KV/ALL_0.1"
#}
#
#output "test_output_secret" {
##  value = data.vault_generic_secret.test111.data.gg_client_password
##  value = nonsensitive(data.vault_generic_secret.test111.data.gg_client_password)
#  value = nonsensitive(data.vault_generic_secret.test_vault_secrets.data.NEXUS_USERNAME)
##  sensitive = true
#}