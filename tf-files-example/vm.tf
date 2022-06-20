
resource "di_vm" "vm1" {
	group_id        = var.group_id
	project_id      = var.project_id
	service_name    = "apvorobev-tf-test-vm"
	ir_group        = "vm"
	os_name         = "rhel"
	os_version      = "7.9"
	virtualization  = "openstack"
	fault_tolerance = "stand-alone"
#	region          = "pd20-okvm3"
	flavor          = "m1.tiny"
	disk            = 50
	zone            = "okvm1"
	volume {
		size = 50
		path = "/test1"
	}
	tag_ids = [
		for tag in di_tag.tags:
		tag.id
		if contains(var.vm_tags, tag.name)
	]
	count           = 1

#	provisioner "remote-exec" {
#		inline = [
#			"ls -la",
#		]
#	}

}

resource "di_tag" "tags" {
	count = length(var.all_tags)
	name = element(var.all_tags, count.index)
}
variable "vm_tags" {
	description = "VM tags"
	type = list(string)
	default = [
		"jenkins",
		"wildfly"
	]
}

variable "all_tags" {
	description = "all tags"
	type = list(string)
	default = [
		"wildfly",
		"jenkins",
		"kibana"
	]
}