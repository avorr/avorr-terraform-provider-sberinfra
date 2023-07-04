data "si_domain" "domain" {
  name = "ГосТех - обучение"
}

data "si_group" "group" {
  domain_id = data.si_domain.domain.id
  name      = "EDU"
}

resource si_vdc "vdc" {
  name        = "terraform-test-vdc"
  group_id    = data.si_group.group.id
  datacenter  = "okvm1" #"PD24R3PROM"
  description = "si.dns.zone"
  limits      = {
    cores   = 100
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

#resource "si_vdc" "import" {}