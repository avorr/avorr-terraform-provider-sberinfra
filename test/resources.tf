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

#"os_name": "rhel" or "altlinux"
#"os_version": "7.9" or "altlinux"

# storage_type = "rbd-1"           ------> SLOW
# storage_type = "rbd-2"           ------> SLOW BACKUP
# storage_type = "iscsi_common"    ------> FAST
# storage_type = "__DEFAULT__"     ------> DEFAULT TYPE

#data "si_domain" "domain" {
#  name = "ГосТех"
#}

#data "si_group" "group" {
#  domain_id = data.si_domain.domain.id
#  name      = "Common"
#}

resource si_project "project" {
  #  ir_group       = "vdc"
  #  type           = "vdc"
  #  ir_type        = "vdc_openstack"
  #  virtualization = "openstack"

  name       = "Test-terraform-project" //requared false
#  group_id   = data.si_group.group.id
  group_id   = "493afd2f-8547-4d2c-9be0-0c37aba6b08c"
  datacenter = "PD24R3PROM" //"okvm1"
  jump_host  = false
  desc       = "test-di.dns.zone"
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
#  network {
#    network_name    = "internal-network2"
#    cidr            = "172.30.100.0/30"
#    dns_nameservers = ["8.8.8.8", "8.8.4.4"]
#    enable_dhcp     = true
#  }
}

locals {
  networks = {for k, v in si_project.project.network : k.network_name => v.network_uuid}
}

resource "si_vm" "vm1" {
  #  ir_group        = "vm"
  #  virtualization  = "openstack"
  #  fault_tolerance = "Stand-alone"
  #  zone            = "internal"

  service_name = "terraform-test-si-vm-0${count.index + 1}"
#  group_id     = data.si_group.group.id
  group_id     = "493afd2f-8547-4d2c-9be0-0c37aba6b08c"
  project_id   = si_project.project.id
  os_name      = "rhel"
  os_version   = "7.9"
  flavor       = "m1.tiny"
  disk         = 50
#  network_uuid = local.networks["internal-network"]
#  tag_ids      = [
#    si_tag.nolabel.id
#  ]
#  volume {
#    size = 50
#  }
#  volume {
#    size         = 50
#    storage_type = "iscsi_common"
#  }
#  volume {
#    size         = 50
#    storage_type = "rbd-1"
#  }
  count = 0
}

resource "si_security_group" "test_sg" {
#  project_id = si_project.project.id
  project_id = "2dc1c80c-2998-424d-b0ce-5c0295b590ff"
  group_name = "terraform-test-sg"
#  security_rule {
#    ethertype = "IPv4"
#    direction = "ingress"
#    protocol = "tcp"
#  }
#  security_rule {
#    ethertype = "IPv4"
#    direction = "ingress"
#    protocol = "tcp"
#    remote_ip_prefix = "172.30.100.0/30"
#    port_range_min = 443
#    port_range_max = 444
#  }
}

#resource "si_project" "import" {
#}

#resource "si_vm" "import" {
#}


#{
#  "security_group": {
#    "group_name": "test",
#    "security_rules": []
#  }
#}



#{
#"security_group": {
#"group_name": "test",
#"security_rules": [
#{
#"ethertype": "IPv4",
#"id": "1",
#"direction": "ingress",
#"port_range_min": "443",
#"port_range_max": "444",
#"protocol": "tcp"
#}
#]
#}
#}