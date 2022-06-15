resource "di_sowa" "sowa" {
	group_id        = data.di_group.devices.id
	project_id      = di_project.terraformtest.id
	service_name    = "sowa1"
	ir_group        = "sowa"
	flavor          = "m2.tiny"
	disk            = 31
	region          = "skolkovo"
	zone            = "edz"
	virtualization  = "openstack"
	os_name         = "rhel"
	os_version      = "8.3"
	fault_tolerance = "stand-alone"
	public_ssh_name = "CAB-SA-CI000160"
	app_params      = {
		version    = "2.3-1602-8"
		joindomain = "sigma.sbrf.ru"
	}
	volume {
		size = 50
		path = "/sowalogs"
	}
	volume {
		size = 12
		path = "/var"
	}
	volume {
		size = 10
		path = "/sowa"
	}
	volume {
		size = 10
		path = "/sowarun"
	}
	volume {
		size = 6
		path = "/usr"
	}
	volume {
		size = 3
		path = "/usr/local/sowa"
	}
	count           = 0
	tag_ids         = [
		#		di_tag.mytag0.id,
		#		di_tag.mytag1.id,
		#		di_tag.mytag2.id,
		#		di_tag.mytag3.id,
	]
}
