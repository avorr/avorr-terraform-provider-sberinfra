# SberCloud Portal properties

variable "flavor" {
  description = "VM flavor name"
  type        = string
  default     = "m1.tiny"
}

variable "disk" {
  description = "VM root disk size"
  type        = number
  default     = 50
}

#variable "os_name" {
#  description = "VM OS name"
#  type        = string
#  default     = "altlinux"
#}

#variable "os_version" {
#  description = "VM OS version"
#  type        = string
#  default     = "8sp"
#}

variable "os_name" {
  description = "VM OS name"
  type        = string
  default     = "rhel"
}

variable "os_version" {
  description = "VM OS version"
  type        = string
  default     = "7.9"
}
