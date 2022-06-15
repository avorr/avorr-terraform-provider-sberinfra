# kafka
resource "di_kafka" "kafka1" {
  ir_group        = "kafka"
  service_name    = "smsaide test kafka cluster"
  group_id        = data.di_group.devices.id
  project_id      = di_project.terraformtest.id
  region          = "skolkovo"
  zone            = "edz"
  fault_tolerance = "cluster"
  virtualization  = "openstack"
  os_name         = "rhel"
  os_version      = "8.5.2"
	public_ssh_name = "CAB-SA-CI000160"
#	flavor          = "m2.tiny"
	flavor          = "m2.small"
  disk            = 30
  app_params = {
    jdk_version      = "11"
    release_type     = "KafkaSE"
    version          = "2.7.2"
    security         = "PLAINTEXT__ZK_PLAIN_NO_AUTH__KAFKA_PLAINTEXT_NO_AUTH"
	  box_server_count = 3
	  dc_quantity      = 1
	  fault_tolerance  = "cluster"
	  joindomain       = "sigma.sbrf.ru"
  }
  volume {
    size = 15
    path = "/zookeeper"
  }
  volume {
    size = 10
    path = "/opt/Apache"
  }
  volume {
    size = 85
    path = "/KAFKADATA"
  }
	count = 0
	tag_ids         = [
#		di_tag.mytag0.id,
#		di_tag.mytag1.id,
#		di_tag.mytag2.id,
#		di_tag.mytag3.id,
#		di_tag.mytag4.id,
	]
}

# ignite
resource "di_ignite" "ignite_in_mem" {
	group_id        = di_project.terraformtest.group_id
	project_id      = di_project.terraformtest.id
	service_name    = "test ignite in mem"
	zone            = "edz"
	region          = "skolkovo"
	virtualization  = "openstack"
	flavor          = "m4.small"
	disk            = 35
#	ir_group        = "ignite_se_persistence"
	ir_group        = "ignite_se"
	os_name         = "rhel"
	os_version      = "8.5.2"
	fault_tolerance = "cluster"
	public_ssh_name = "CAB-SA-CI000160"
	app_params      = {
		ise_email           = "Solovyev.S.Aleksandr@sberbank.ru"
		joindomain          = "sigma.sbrf.ru"
		fault_tolerance     = "cluster"
		version             = "D-01.040.00-4.2120.6"
		box_server_count    = 1
		ise_client_password = <<-EOF
      $ANSIBLE_VAULT;1.1;AES256
      37636561393162643264363630303635346337303561376531626563356239313535333937383362
      6139666435366364616332613464353964316561636338630a353161343430396534613162386335
      35616331633164613434343036373731313839653839616637396161323430653862363830363362
      3632656366353038660a303530653638643165653763656131316132623733386133353262363130
      64323765306136333838626266623863663635656331316165373365633861326435
    EOF
	}
	volume {
		name    = "lvoptignite"
		size    = 6
		path    = "/opt/ignite"
		fs_type = "ext4"
	}
	volume {
		name    = "lvoptignitelogs"
		size    = 19
		path    = "/opt/ignite/logs"
		fs_type = "ext4"
	}
	count = 0
	tag_ids         = [
#		di_tag.mytag0.id,
#		di_tag.mytag1.id,
#		di_tag.mytag2.id,
#		di_tag.mytag3.id,
#		di_tag.mytag4.id,
	]
}

# ignite2
resource "di_ignite" "ignite_se_persistence" {
	group_id        = di_project.terraformtest.group_id
	project_id      = di_project.terraformtest.id
	service_name    = "test ignite_se_persistence"
	zone            = "edz"
	region          = "skolkovo"
	virtualization  = "openstack"
	flavor          = "m2.small"
	disk            = 35
	ir_group        = "ignite_se_persistence"
	os_name         = "rhel"
	os_version      = "8.5.2"
	fault_tolerance = "cluster"
	public_ssh_name = "CAB-SA-CI000160"
	app_params      = {
		ise_email           = "Solovyev.S.Aleksandr@sberbank.ru"
		joindomain          = "sigma.sbrf.ru"
		fault_tolerance     = "cluster"
		version             = "D-01.040.00-4.2120.6"
		box_server_count    = 1
		ise_client_password = <<-EOF
      $ANSIBLE_VAULT;1.1;AES256
      37636561393162643264363630303635346337303561376531626563356239313535333937383362
      6139666435366364616332613464353964316561636338630a353161343430396534613162386335
      35616331633164613434343036373731313839653839616637396161323430653862363830363362
      3632656366353038660a303530653638643165653763656131316132623733386133353262363130
      64323765306136333838626266623863663635656331316165373365633861326435
    EOF
	}
	volume {
		size    = 10
		path    = "/opt/ignite"
	}
	volume {
		size    = 12
		path    = "/opt/ignite/wal"
	}
	volume {
		size    = 20
		path    = "/opt/ignite/wal_archive"
	}
	volume {
		size    = 8
		path    = "/opt/ignite/logs"
	}
	volume {
		size    = 75
		path    = "/opt/ignite/data"
	}
	count = 0
	tag_ids         = [
#		di_tag.mytag0.id,
#		di_tag.mytag1.id,
#		di_tag.mytag2.id,
#		di_tag.mytag3.id,
#		di_tag.mytag4.id,
	]
}


# patroni
resource "di_patroni" "patroni" {
	group_id        = data.di_group.devices.id
	project_id      = di_project.terraformtest.id
	service_name    = "test patroni 1"
	ir_group        = "patroni"
	flavor          = "m4.tiny"
	disk            = 35
	region          = "skolkovo"
	zone            = "edz"
	virtualization  = "openstack"
	os_name         = "rhel"
	os_version      = "8.5.2"
  fault_tolerance = "cluster"
	public_ssh_name = "CAB-SA-CI000160"
	app_params      = {
		joindomain           = "sigma.sbrf.ru"
		max_connections      = 101
		version              = 13
		fault_tolerance      = "cluster"
		box_server_count     = 2
		dc_quantity          = 1
		postgres_db_password = <<-EOF
      $ANSIBLE_VAULT;1.1;AES256
      65346162336565383039303663316639363062363339323135616530373131343737633166366531
      6331353633383531633530333335636638343131306530660a306530373232363365356338663537
      39383266356666303464363934326239636661323235633932646231306136663961376436623830
      6466373231393130370a396233333031313565316263313632356130376638363031363965396438
      33623434616563383537346332336633623939303230393766306334396566623938
    EOF
		postgres_dbs         = "database123:17756439"
		# Название базы данных и владелец через ':'
		# Имя БД должно быть длинной от 2 до 8 символов и содержать только латинские буквы в нижнем регистре и цифры.
		# Имя пользователя должно содержать только латинские буквы в нижнем регистре и цифры.
    # Не допускается использовать служебные названия.
	}
	volume {
		size    = 50
		path    = "/pgarclogs"
	}
	volume {
		size    = 50
		path    = "/pgdata"
	}
	volume {
		size    = 50
		path    = "/pgerrorlogs"
	}
	volume {
		size    = 1
		path    = "/pgbackup"
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
