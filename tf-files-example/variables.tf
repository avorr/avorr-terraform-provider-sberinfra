variable "allow_unverified_ssl" {
  description = "Strict ssl verify"
  type = bool
  default = true
}

variable "project_name" {
  description = "Project Name"
  type = string
  default = "efs"
}

variable "env" {
  description = "Project environment"
  type = string
  default = "u"
}

# SberCloud Portal properties
variable "enable_sc" {
  description = "Flag for cloud type"
  type = bool
  default = false
}

variable "api_base_url" {
  description = "Cloud base url"
  type = string
  default = "https://portal.gostech.novalocal/api/v1"
}

variable "token" {
  description = "Cloud AUTH token"
  type = string
  default = ""
}

variable "project_id" {
  description = "Project ID"
  type = string
#  default = "e17039cb-c596-40c8-838e-f73256091947"
#  default = "37f80fa9-bd9f-479e-ba49-f2d66376545b" #PD20
#  default = "e5801255-d4e6-41fb-a492-615fdfb5764c" #PD20 gt-common-admins-uat-junior
#  default = "08e3fb48-c212-486e-a77a-809a73caa440" #PD15 dvp-dev-admin
  default = "db95faf0-f1b8-46ef-b371-3ee40795c432" #PD24 gt-rosim-nt-dmz
}

variable "group_id" {
  description = "Group ID"
  type = string
#  default = "ca0519ee-ab9b-45f8-8caa-16993b273627"
#  default = "d92966b6-0345-43a1-8f68-bd5ae85640f2" #PD20
#  default = "86e75f8c-24c2-448b-b8d4-eb614ad82234" #PD15
  default = "58194ebe-56a7-4cde-9e3b-03731993a25e" #PD24 ros nt
}

variable "region" {
  // Must be empty for Publick Cloud
  description = "Virtualization region"
  type = string
  default = ""
}

variable "zone" {
  description = "Network Zone"
  type = string
  default = "okvm1"
}

variable "virtualization" {
  description = "Virtualization type"
  type = string
  default = "openstack"
}

variable "flavor" {
  description = "VM flavor name"
  type = string
  default = "m1.tiny"
}

variable "disk_size" {
  description = "VM root disk size"
  type = number
  default = 50
}

variable "extra_disk_size" {
  description = "VM extra disk size"
  type = string
  default = ""
}




// variable "extra_disk_size_1" {
//   default = ""
// }

// variable "extra_disk_size_2" {
//   default = "10"
// }

variable "ir_group" {
  description = "Information resource type"
  type = string
  default = "vm"
}

// variable "os_name" {
//   description = "VM OS name"
//   type = string
//   default = "altlinux"
// }

// variable "os_version" {
//   description = "VM OS version"
//   type = string
//   default = "8sp"
// }

variable "os_name" {
  description = "VM OS name"
  type = string
  default = "rhel"
}

variable "os_version" {
  description = "VM OS version"
  type = string
  default = "7.9"
}
