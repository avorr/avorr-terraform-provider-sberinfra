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

resource si_project "project" {
  name       = "Test-terraform-project" //requared false
  group_id   = data.si_group.group.id
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
  service_name = "terraform-test-${format("%02d", count.index + 1)}.${si_project.project.desc}"
  group_id     = data.si_group.group.id
  project_id   = si_project.project.id
  os_name      = "rhel"
  os_version   = "7.9"
  flavor       = "m1.tiny"
  #  disk         = 50 // Optional param | Temporary
  network_uuid = local.networks["internal-network"]
  hdd {
    size = 50
    #    storage_type = "iscsi-fast-01"
  }
  tag_ids = [
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
    size         = 100
    storage_type = "rbd-1"
  }
  count = 1
}

resource "si_security_group" "group" {
  group_name = "test-group"
  project_id = si_project.project.id
  security_rule {
    ethertype        = "IPv4"            //TypeString "IPv4", "IPv6"
    direction        = "ingress"         //TypeString "ingress", "egress"
    protocol         = "tcp"             //TypeString "tcp", "udp", "icmp"
    remote_ip_prefix = "172.21.21.0/28"  //TypeString
    port_range_min   = 443               //TypeInt
    port_range_max   = 444               //TypeInt
  }
  security_rule {
    ethertype        = "IPv4"
    direction        = "egress"
    protocol         = "tcp"
    remote_ip_prefix = "172.21.21.0/28"
    port_range_min   = 443
    port_range_max   = 444
  }
  security_rule {
    ethertype        = "IPv4"
    direction        = "egress"
    protocol         = "udp"
    remote_ip_prefix = "172.21.21.0/0"
    port_range_min   = 443
    port_range_max   = 444
  }
}



#resource "si_project" "import" {
#}

#resource "si_vm" "import" {
#}