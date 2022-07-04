/*
resource "di_tag" "tags" {
	count = length(var.all_tags)
	name = element(var.all_tags, count.index)
}

variable "vm_tags" {
	description = "VM tags"
	type = list(string)
	default = [
		"TESTTAG",
		"jenkins",
		"wildfly"
	]
}

variable "all_tags" {
	description = "all tags"
	type = list(string)
	default = [
		"TESTTAG",
		"vpn",
		"jenkins",
		"wildfly",
		"kibana"
	]
}
*/

/*
variable "disks" {
  type = map
  default = {}
}

variable "disks" {
	type = map
	default = {
		disk1 = {
			size: 50
			storage_type = "rbd-1"
		}
		disk2 = {
			size: 100
			storage_type = "iscsi_common"
		}
	}
}
*/

resource "di_vm" "vm1" {
	group_id        = var.group_id
	project_id      = var.project_id
	service_name    = "TERRAFORM-TEST"
	ir_group        = "vm"
	os_name         = "rhel"
	os_version      = "7.9"
	virtualization  = "openstack"
	fault_tolerance = "stand-alone"
	flavor          = "m1.tiny"
	disk            = 50
	zone            = "okvm1"
	count           = 1
#	public_ssh_name = "avorr"
/*
	provisioner "remote-exec" {
		inline = [
			"ls -la /",
			"sudo touch /opt/TESTFILE"
		]
		connection {
			type     = "ssh"
			user     = self.user
			password = self.password
			host     = self.ip
			port     = 9022
			private_key = file("~/.ssh/id_rsa")
		}
	}

	tag_ids = [
		for tag in di_tag.tags:
		tag.id
		if contains(var.vm_tags, tag.name)
	]
	dynamic volume {
		for_each = var.disks
		content {
			size = volume.value.size
			storage_type = volume.value.storage_type
		}
	}
*/
}


#data "di_domain" "domain" {
#	name = "ГосТех"
#}


#data "di_stand_type" "dev" {
#	name = "DEV"
#}

#data "di_group" "Common" {
#	name           = "Common"
#	domain_id      = data.di_domain.domain.id
#}


#data "di_stand_type" "dev" {
#	name = "DEV"
#}

#data "di_as" "ec" {
#	code           = "CI01808661"
#	domain_id      = data.di_domain.domain.id
#}


#resource "di_project" "terraformtest" {
#	name                = "TerraformTest"
#	group_id            = var.group_id
#	stand_type_id       = data.di_stand_type.dev.id
#	stand_type_id       = "vdc"
#	app_systems_ci      = "DEV"
#}

/*
resource "di_project" "terraformtest" {
#resource "di_project" "gt-common-admins-uat-junior" {
	name                = "TerraformTest"
	group_id            = var.group_id
	stand_type_id       = data.di_stand_type.dev.id
#	app_systems_ci      = data.di_as.ec.code
#	type                = "vdc"

	datacenter          = "PD23R1PSI"
	app_systems_ci      = "DEV"
	jump_host           = "false"
}
*/


#output "ip" {
#	value = di_vm.vm1[0].ip
#}


#output "password" {
#	value = di_vm.vm1[0].password
#}

#output "tag_id" {
#	value = di_tag.tags
#}

#output "domain" {
#	value = data.di_domain.domain
#}

#output "group" {
#	value = data.di_group.Common
#}
#
#output "stand_type" {
#	value = data.di_stand_type.dev
#}
