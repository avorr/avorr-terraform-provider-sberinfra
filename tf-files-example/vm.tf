data "di_domain" "domain" {
  name = "ГосТех"
}

data "di_group" "group" {
  name      = "Common"
  domain_id = data.di_domain.domain.id
}

data "di_siproject" "project" {
  group_id = data.di_group.group.id
  name     = "gt-common-admins-uat-junior"
#  name     = "gt-common-admins"
}

output "domain_id" {
  value = data.di_domain.domain.id
}

output "group_id" {
  value = data.di_group.group.id
}

output "projects_id" {
  value = data.di_siproject.project.id
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
    #    disk1 = {
    #      size : 50
    #      storage_type = "rbd-1"
    #      storage_type = "__DEFAULT__"
    #    }
    #    disk2 = {
    #      size : 100
    #      storage_type = "iscsi_common"
    #    }
  }
}

resource "di_tag" "tags" {
  count = length(var.all_tags)
  name  = element(var.all_tags, count.index)
}

resource "di_vm" "vm1" {
  group_id        = data.di_group.group.id
  project_id      = data.di_siproject.project.id
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

  tag_ids = [
  for tag in di_tag.tags :
  tag.id
  if contains(var.vm_tags, tag.name)
  ]
  dynamic volume {
    for_each = var.disks
    content {
      size         = volume.value.size
      storage_type = volume.value.storage_type
    }
  }
  #  provisioner "remote-exec" {
  #    inline = [
  #      "ls -la /",
  #      "sudo touch /opt/TESTFILE"
  #    ]
  #    connection {
  #      type     = "ssh"
  #      user     = self.user
  #      password = self.password
  #      host     = self.ip
  #      port     = 9022
  #    }
  #  }
}
