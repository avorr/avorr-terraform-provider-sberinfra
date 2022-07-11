
resource "di_tag" "tags" {
	count = length(var.all_tags)
	name = element(var.all_tags, count.index)
}

variable "vm_tags" {
	description = "VM tags"
	type = list(string)
	default = [
		"jenkins",
		"wildfly"
	]
}

variable "all_tags" {
	description = "all tags"
	type = list(string)
	default = [
		"vpn",
		"jenkins",
		"wildfly",
		"kibana"
	]
}



#variable "disks" {
#  type = map
#  default = {}
#}

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


resource "di_vm" "vm1" {
	group_id        = data.di_si-group.group.id
	project_id      = data.di_si-project.project.id
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
}




#data "di_stand_type" "dev" {
#	name = "DEV"
#}



#data "di_stand_type" "dev" {
#	name = "DEV"
#}

#data "di_as" "ec" {
#	code           = "CI01808661"
#	domain_id      = data.di_si-domain.domain.id
#}


data "di_si-domain" "domain" {
	name = "ГосТех"
}

data "di_si-group" "group" {
	name           = "Common"
	domain_id      = data.di_si-domain.domain.id
}

data "di_si-project" "project" {
	group_id = data.di_si-group.group.id
	name = "gt-common-admins-uat-junior"
#	name = "gt-common-admins-prod-junior"
}

output "domain_id" {
	value = data.di_si-domain.domain.id
}

output "group_id" {
	value = data.di_si-group.group.id
}

output "projects_id" {
	value = data.di_si-project.project.id
}
