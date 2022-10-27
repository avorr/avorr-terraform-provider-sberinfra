#resource "di_project" "terraformtest" {
#  name                = "(DEV) test-di"
#  group_id            = data.di_group.devices.id
#  app_systems_ci      = data.di_as.ec.code
#  stand_type_id       = data.di_stand_type.dev.id
#}

resource "di_project" "terraformtest" {
  name                = "TerraformTest"
  group_id            = data.di_group.devices.id
  stand_type_id       = data.di_stand_type.dev.id
  app_systems_ci      = data.di_as.ec.code
}

#resource "di_project" "terraformtest2" {}



