resource "di_vm" "vmtest" {
  group_id        = data.di_group.devices.id
  project_id      = di_project.terraformtest.id
  service_name    = "vm1"
  ir_group        = "vm"
  os_name         = "rhel"
  os_version      = "8.5.2"
  #	virtualization  = "openstack"
  virtualization  = "vmware"
  fault_tolerance = "stand-alone"
  region          = "skolkovo"
  flavor          = "m1.tiny"
  disk            = 31
  zone            = "edz"
  public_ssh_name = "CAB-SA-CI000160"
  app_params      = {
    joindomain = "delta.sbrf.ru"
  }
  #	volume {
  #		size = 3
  #		path = "/.test-di"
  #	}
  #	volume {
  #		size = 5
  #		path = "/test2"
  #	}
  count = 0
  #	tag_ids         = [
  #		di_tag.mytag0.id,
  #		di_tag.mytag1.id,
  #		di_tag.mytag2.id,
  #		di_tag.mytag3.id,
  #		di_tag.mytag4.id,
  #	]
}
