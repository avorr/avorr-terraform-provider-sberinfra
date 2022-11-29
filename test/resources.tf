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

data "si_domain" "domain" {
  name = "ГосТех"
}

data "si_group" "group" {
  domain_id = data.si_domain.domain.id
  name      = "Common"
}

resource si_vdc "vdc" {
  name        = "terraform-test-vdc" //requared false
  group_id    = data.si_group.group.id
  datacenter  = "PD24R3PROM" //"okvm1"
  description = "si.dns.zone"
  limits      = {
    vcpu    = 100
    ram     = 10000
    storage = 1000
  }
  network {
    name    = "internal-network"
    cidr    = "172.31.0.0/20"
    dns     = ["8.8.8.8", "8.8.4.4"]
    dhcp    = true
    default = true
  }
  network {
    name = "internal-network2"
    cidr = "172.30.100.0/30"
    dns  = ["8.8.8.8", "8.8.4.4"]
    dhcp = true
  }
}

locals {
  networks = {for k, v in si_vdc.vdc.network : k.name => v.id}
}

resource "si_vm" "vm" {
  service_name = "terraform-test-${format("%02d", count.index + 1)}.${si_vdc.vdc.description}"
  group_id     = data.si_group.group.id
  vdc_id       = si_vdc.vdc.id
  os_name      = "rhel"
  os_version   = "7.9"
  flavor       = "m1.tiny"
  disk         = {
    size = 50
    #    storage_type = "iscsi-fast-01"
  }
  network_id = local.networks["internal-network"]
  tag_ids    = [
    si_tag.nolabel.id
  ]
  security_groups = [
    si_security_group.iam.id,
    si_security_group.kafka.id,
  ]
  volume {
    size = 50
  }
  volume {
    size         = 50
    storage_type = "iscsi_common"
  }
  volume {
    size         = 100
    storage_type = "rbd-1"
  }
  count = 1
}

resource "si_security_group" "iam" {
  name   = "iam"
  vdc_id = si_vdc.vdc.id
  security_rule {
    ethertype        = "IPv4"
    direction        = "ingress"
    protocol         = "tcp"
    remote_ip_prefix = "172.21.21.10/28"
    port_range_min   = 443
    port_range_max   = 444
  }
  security_rule {
    ethertype        = "IPv4"
    direction        = "ingress"
    protocol         = "tcp"
    remote_ip_prefix = "172.21.21.10/28"
    port_range_min   = 80
    port_range_max   = 80
  }
}

resource "si_security_group" "kafka" {
  name   = "kafka"
  vdc_id = si_vdc.vdc.id
  security_rule {
    ethertype      = "IPv4"
    direction      = "ingress"
    protocol       = "tcp"
    port_range_min = 9092
    port_range_max = 9092
  }
  security_rule {
    ethertype      = "IPv4"
    direction      = "ingress"
    protocol       = "tcp"
    port_range_min = 2181
    port_range_max = 2181
  }
}

#resource "si_vdc" "import" {
#}

#resource "si_vm" "import" {
#}