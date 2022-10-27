data "di_domain" "domain" {
    name = "ГосТех"
}

data "di_group" "group" {
  domain_id = data.di_domain.domain.id
    name      = "Common"
}

output "domain_id" {
  value = data.di_domain.domain.id
}

output "group_id" {
  value = data.di_group.group.id
}

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
  networks = {for k, v in di_project.project.network : k.network_name => v.network_uuid}
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
  zone            = "internal"
  network_uuid    = local.networks["internal-network"]
  tag_ids         = [
#    di_tag.jenkins.id
  ]

  volume {
    size = 50
    #    storage_type = "rbd-1"
  }
  count = 0
}

#output "ni" {
#  value = local.networks["internal-network"]
#}
