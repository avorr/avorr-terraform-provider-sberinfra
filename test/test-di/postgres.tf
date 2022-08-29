resource "di_postgres" "postgres" {
	group_id        = data.di_group.devices.id
	project_id      = di_project.terraformtest.id
	service_name    = "TF SmartAppIDE PG SDFD edz (KSSD-907)"
	ir_group        = "postgres"
	flavor          = "m1.small"
	disk            = 37
	region          = "skolkovo"
	zone            = "edz"
	virtualization  = "openstack"
	os_name         = "rhel"
	os_version      = "8.5.2"
	fault_tolerance = "stand-alone"
	public_ssh_name = "CAB-SA-CI000160"
	app_params      = {
		joindomain           = "sigma.sbrf.ru"
		max_connections      = 200
		version              = 13
		postgres_db_name     = "smmjaisdfd"
		postgres_db_user     = "smmjaisdfd"
		postgres_db_password = <<-EOF
      $ANSIBLE_VAULT;1.1;AES256
      65346162336565383039303663316639363062363339323135616530373131343737633166366531
      6331353633383531633530333335636638343131306530660a306530373232363365356338663537
      39383266356666303464363934326239636661323235633932646231306136663961376436623830
      6466373231393130370a396233333031313565316263313632356130376638363031363965396438
      33623434616563383537346332336633623939303230393766306334396566623938
    EOF
	}
	volume {
		size = 50
		path = "/pgarclogs"
	}
	volume {
		size = 50
		path = "/pgdata"
	}
	volume {
		size = 50
		path = "/pgerrorlogs"
	}
	volume {
		size = 5
		path = "/pgbackup"
	}
	count           = 0
	tag_ids         = [
#		di_tag.mytag0.id,
#		di_tag.mytag1.id,
#		di_tag.mytag2.id,
#		di_tag.mytag3.id,
#		di_tag.mytag4.id,
	]
}
