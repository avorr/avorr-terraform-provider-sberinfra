variable "group_id" {
  description = "Group ID"
  type = string
    default = "52ffd9f6-fbc0-4ddc-bf99-b092c6d0351a" #PD23 gostech common
}

variable "project_id" {
  description = "Project ID"
  type = string
  default = "5639188a-fc11-48fb-a1f9-db86a5eae74b" #PD23 gt-common-admins-uat-junior
}

#variable "disks" {
#  type = map
#  default = {}
#}

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
