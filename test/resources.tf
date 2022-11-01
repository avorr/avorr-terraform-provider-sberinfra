# m1.tiny   = 1/1
# m1.medium = 4/4
# m1.large	= 8/8
# m2.small	= 2/4
# m2.medium = 4/8
# m2.large	= 8/16
# m2.xlarge = 16/32
# m4.small	= 2/8
# m4.medium = 4/16
# m4.large	= 8/32
# m4.xlarge = 16/64
# m6.medium = 4/24
# m8.medium = 4/32
# m8.large	= 8/64


# storage_type = "rbd-1"           ------> SLOW
# storage_type = "rbd-2"           ------> SLOW BACKUP
# storage_type = "iscsi_common"    ------> FAST
# storage_type = "__DEFAULT__"     ------> DEFAULT TYPE

data "si_domain" "domain" {
  name = "ГосТех"
}

data "si_group" "group" {
  domain_id = data.si_domain.domain.id
  name      = "Common"
}

output "domain_id" {
  value = data.si_domain.domain.id
}

output "group_id" {
  value = data.si_group.group.id
}

resource si_project "project" {
  name           = "Test-terraform-project" //requared false
  group_id       = data.si_group.group.id
  datacenter     = "PD24R3PROM" //"okvm1"
  jump_host      = false
  desc           = "test-di.dns.zone"
#  ir_group       = "vdc"
#  type           = "vdc"
#  ir_type        = "vdc_openstack"
#  virtualization = "openstack"
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
    is_default      = true
  }
  network {
    network_name    = "internal-network2"
    cidr            = "172.30.100.0/30"
    dns_nameservers = ["8.8.8.8", "8.8.4.4"]
    enable_dhcp     = true
  }
}

locals {
  networks = {for k, v in si_project.project.network : k.network_name => v.network_uuid}
}

resource "si_vm" "vm1" {
  service_name    = "terraform-test-si-vm-0${count.index + 1}"
  group_id        = data.si_group.group.id
  project_id      = si_project.project.id
  os_name         = "rhel"
  os_version      = "7.9"
  flavor          = "m1.tiny"
  disk            = 50
  network_uuid    = local.networks["internal-network"]
#  ir_group        = "vm"
#  virtualization  = "openstack"
#  fault_tolerance = "Stand-alone"
#  zone            = "internal"

  tag_ids         = [
    si_tag.nolabel.id
  ]
  volume {
    size = 50
  }
  volume {
    size         = 50
    storage_type = "iscsi_common"
  }
  volume {
    size         = 50
    storage_type = "rbd-1"
  }
  count = 1
}

#resource "si_vm" "import" {
#}