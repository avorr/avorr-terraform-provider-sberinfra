resource "si_vm" "vm" {
  service_name = "terraform-test-${format("%02d", count.index + 1)}"
  group_id     = data.si_group.group.id
  vdc_id       = si_vdc.vdc.id
  ir_type      = "os_alt"
  os_name      = "altlinux" # "rhel"
  os_version   = "8sp" # "7.9"
  flavor       = "m1.tiny"
  description  = "testing vm"
  #  public_ssh_name = "id_rsa.pub"
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
    name = "postgres"
  }
  volume {
    size         = 50
    name         = "kafka"
    storage_type = "iscsi_common"
  }
  count = 1
}
