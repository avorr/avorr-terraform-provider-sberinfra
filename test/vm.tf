/*
resource "di_tag" "tags" {
  count = length(var.all_tags)
  name  = element(var.all_tags, count.index)
}

variable "vm_tags" {
  description = "VM tags"
  type        = list(string)
  default     = [
    "jenkins",
    "wildfly"
  ]
}

variable "all_tags" {
  description = "all tags"
  type        = list(string)
  default     = [
    "vpn",
    "jenkins",
    "wildfly",
    "kibana"
  ]
}

variable "disks" {
  type    = map
  default = {
    disk1 = {
      size : 50
      storage_type = "rbd-1"
    }
    disk2 = {
      size : 100
      storage_type = "iscsi_common"
    }
  }
}

*/
/*
resource "di_vm" "vm1" {
	group_id        = data.di_group.group.id
	project_id      = data.di_project.project.id
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
*/


#data "di_stand_type" "dev" {
#	name = "DEV"
#}

data "di_domain" "domain" {
#  name = "Гостех - авторизация"
    name = "ГосТех"
#    name = "Росимущество"
}

data "di_group" "group" {
  domain_id = data.di_domain.domain.id
#  name      = "IAM для внешних интеграций"
    name      = "Common"
  #  name      = "НТ"
  #  name      = "ПСИ"
}

output "domain_id" {
  value = data.di_domain.domain.id
}

output "group_id" {
  value = data.di_group.group.id
}

/*
resource di_project "project" {
  ir_group       = "vdc"
  type           = "vdc"
  ir_type        = "vdc_openstack"
  virtualization = "openstack"
  name           = "Test-project" //requared false
  group_id       = data.di_group.group.id
  #  group_id = "52ffd9f6-fbc0-4ddc-bf99-b092c6d0351a"
  #  datacenter = "PD23R3PROM"
  datacenter     = "PD20R3PROM"
  #  datacenter = "okvm1"
  jump_host      = false
  desc           = "test-di.dns.zone"
  limits {
    //requared false
    cores_vcpu_count  = 100
    ram_gb_amount     = 10000
    storage_gb_amount = 1000
  }
  network {
    network_name    = "internal-network"
    cidr            = "172.31.0.0/20"
    dns_nameservers = ["8.8.8.8", "8.8.4.4"]
    enable_dhcp     = true
  }
}
*/

/*
resource di_project "project" {
  ir_group       = "vdc"
  type           = "vdc"
  ir_type        = "vdc_openstack"
  virtualization = "openstack"
  name           = "terraform-test-si-project" //requared false
  group_id       = data.di_group.group.id
  datacenter     = "PD24R3PROM"
  #  datacenter     = "openstack"
  jump_host      = false
  desc           = "test-di.dns.zone"
  limits {
    cores_vcpu_count  = 1000 //requared false
    ram_gb_amount     = 10000
    storage_gb_amount = 1000
  }
  network {
    network_name    = "internal-network"
    cidr            = "172.31.0.0/20"
    dns_nameservers = ["8.8.8.8", "8.8.4.4"]
    enable_dhcp     = true
    is_default      = true
  }
  network {
    network_name    = "test-di-network"
    cidr            = "172.31.10.0/29"
    dns_nameservers = ["8.8.8.8", "8.8.4.4"]
    enable_dhcp     = true
#    is_default      = true
  }
}


resource "di_vm" "vm1" {
  group_id        = data.di_group.group.id
  project_id      = di_project.project.id
  service_name    = "terraform-test-di-vm-0${count.index + 1}"
  ir_group        = "vm"
  os_name         = "rhel"
  os_version      = "7.9"
  virtualization  = "openstack"
  fault_tolerance = "Stand-alone"
  flavor          = "m1.tiny"
  disk            = 50
  zone            = di_project.project.datacenter
#  zone            = "internal"
  volume {
    size = 50
    #    storage_type = "rbd-1"
  }
#  tag_ids = [
#    di_tag.jenkins.id
#  ]
  count = 0
}
*/

resource di_project "project" {
  ir_group       = "vdc"
  type           = "vdc"
  ir_type        = "vdc_openstack"
  virtualization = "openstack"
  name           = "Test-terraform-project" //requared false
  group_id       = data.di_group.group.id
#  datacenter     = "okvm1"
  datacenter     = "PD24R3PROM"
  jump_host      = false
  desc           = "test-di.dns.zone"
  limits {
    cores_vcpu_count  = 100    //
    ram_gb_amount     = 10000   // requared false
    storage_gb_amount = 1000    //
  }
  network {
    network_name    = "internal-network"
    cidr            = "172.31.0.0/20"
    dns_nameservers = ["8.8.8.8", "8.8.4.4"]
    enable_dhcp     = true
#    is_default      = true
  }
  network {
    network_name    = "internal-network2"
    cidr            = "172.30.100.0/30"
    dns_nameservers = ["8.8.8.8", "8.8.4.4"]
    enable_dhcp     = true
    is_default      = true
  }
}

#resource "di_tag" "jenkins" {
#  name = "jenkins"
#}

#"cd57f3d7-176f-4295-aae8-c5ee16716a82"

#locals {
#  networks = {for k, v in di_project.project.network : k.network_name => v.network_uuid}
#}

resource "di_vm" "vm1" {
  group_id        = data.di_group.group.id
#  project_id      = "0c479d29-898b-4b1a-b3cb-bc408723aa4a"
  project_id      = di_project.project.id
  service_name    = "terraform-test-di-vm-0${count.index + 1}"
  ir_group        = "vm"
  os_name         = "rhel"
  os_version      = "7.9"
  virtualization  = "openstack"
  fault_tolerance = "Stand-alone"
  flavor          = "m1.tiny"
  disk            = 50
  zone            = "internal"
  #  network_uuid    = "857d0f59-9479-4028-a583-74f1edaf22a0"
  #  network_name    = "internal-network"
#  network_uuid    = local.networks["internal-network"]
  tag_ids         = [
#    di_tag.jenkins.id
  ]
  count = 0
}

#output "ni" {
#  value = local.networks["internal-network"]
#}
