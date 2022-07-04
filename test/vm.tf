
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
#		path = "/test1"
#	}
#	volume {
#		size = 5
#		path = "/test2"
#	}
	count           = 0
#	tag_ids         = [
#		di_tag.mytag0.id,
#		di_tag.mytag1.id,
#		di_tag.mytag2.id,
#		di_tag.mytag3.id,
#		di_tag.mytag4.id,
#	]
}

resource "di_nginx" "nginx1" {
  ir_group        = "nginx"
  service_name    = "smsaide test nginx"
  group_id        = data.di_group.devices.id
  project_id      = di_project.terraformtest.id
  region          = "skolkovo"
  zone            = "edz"
  fault_tolerance = "stand-alone"
  virtualization  = "openstack"
  os_name         = "rhel"
  os_version      = "8.5.2"
	public_ssh_name = "CAB-SA-CI000160"
	flavor          = "m4.tiny"
  disk            = 30
	app_params = {
		version      = "1.20.2-1"
		nginx_geoip  = "No"
		nginx_brotli = "No"
		joindomain   = "sigma.sbrf.ru"
	}
	volume {
		size = 1
		path = "/test1"
	}
	volume {
		size = 2
		path = "/test2"
	}
	volume {
		size = 3
		path = "/test3"
	}
	volume {
		size = 4
		path = "/test4"
	}
	volume {
		size = 10
		path = "/opt/nginx"
	}
  count          = 0
	tag_ids         = [
#		di_tag.mytag0.id,
#		di_tag.mytag1.id,
#		di_tag.mytag2.id,
#		di_tag.mytag3.id,
#		di_tag.mytag4.id,
	]
}

resource "di_openshift" "osproject" {
	service_name    = "Тестовый Проект"
	ir_group        = "project"
	os_name         = "oc4_project"
#	os_version      = "2019"
	virtualization  = "openshift"
	fault_tolerance = "stand-alone"
#	region          = "skolkovo"
	region          = "acod-5"
	cpu             = 2
	ram             = 5
#	disk            = 50
	zone            = "edz"
#	public_ssh_name = "CAB-SA-CI000160"
	group_id        = data.di_group.devices.id
	project_id      = di_project.terraformtest.id

	app_params = {
		name_project = "testproject1"
		admin_user   = "cab-sa-ci000160"
	}
	count      = 0
#	tag_ids         = [
#		di_tag.mytag0.id,
#		di_tag.mytag1.id,
#		di_tag.mytag2.id,
#		di_tag.mytag3.id,
#	]
}

# postgres SE
resource "di_postgres_se" "postgres_se" {
	group_id        = data.di_group.devices.id
	project_id      = di_project.terraformtest.id
	service_name    = "postgres se test server"
	ir_group        = "postgres_se"
	flavor          = "m4.tiny"
	disk            = 35
	region          = "skolkovo"
	zone            = "edz"
	virtualization  = "openstack"
	os_name         = "rhel"
	os_version      = "8.5.2"
	fault_tolerance = "stand-alone"
	public_ssh_name = "CAB-SA-CI000160"
	app_params      = {
		as_tuz               = "CAB-SA-CI000160" # Список технологических учетных записей
		version              = "4.2.6" # Версия PostgreSQL Sber Edition
		as_admins            = "17756439" # Список администраторов АС
		joindomain           = "sigma.sbrf.ru"
		schema_name          = "schema1"
		database_name        = "db1"
		security_level       = "K4"
		fault_tolerance      = "stand-alone" # хз зачем оно в ди 2 раза
		installation_type    = "standalone-postgresql-only"
		tablespace_name      = "ts1"
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
		size = 1
		path = "/pgbackup"
	}
	count           = 0
}

# ELK
resource "di_elk" "elk" {
	group_id        = data.di_group.devices.id
	project_id      = di_project.terraformtest.id
	service_name    = "test elk 4"
	ir_group        = "elk"
	flavor          = "m2.tiny"
	disk            = 35
	region          = "skolkovo"
	zone            = "edz"
	virtualization  = "openstack"
	os_name         = "rhel"
	os_version      = "8.5.2"
	fault_tolerance = "stand-alone"
	public_ssh_name = "CAB-SA-CI000160"
	app_params      = {
		joindomain   = "sigma.sbrf.ru"
#		elk_set      = "Elasticsearch + Logstash + Kibana"
		elk_set      = "Kibana"
		version      = "7.13.2"
		java_version = "1.8.0"
	}
	volume {
		path    = "/usr"
		size    = 40
	}
	volume {
		size    = 40
		path    = "/opt/elastic"
	}
	volume {
		size    = 40
		path    = "/home"
	}
	volume {
		size    = 10
		path    = "/var"
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
