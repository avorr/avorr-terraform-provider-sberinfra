/*
####m1
m1.tiny  1/1
m1.small  2/2
m1.medium  4/4
m1.large  8/8
m1.xlarge  16/16


####kasper_n2
kasper_n2.tiny 6/6
kasper_n2.small 6/10
kasper_n2.medium 6/18
kasper_n2.large 6/24
kasper_n2.xlarge  6/32


####kasper_n1
kasper_n1.small 4/6
kasper_n1.medium  4/10
kasper_n1.large  4/18


####kasper_n3
kasper_n3.small  10/18
kasper_n3.medium  10/32
kasper_n3.large 10/50
kasper_n3.xlarge  10/64


####m1
m2.tiny  1/2
m2.small  2/4
m2.medium  4/8
m2.large  8/16
m2.xlarge  16/32
m2.xxlarge  32/64


####m3
m3.medium  4/12
m3.large  8/24


####m4
m4.tiny  1/4
m4.small  2/8
m4.medium  4/16
m4.large  8/32
m4.xlarge  16/64


####m6
m6.tiny  1/6
m6.small  2/12
m6.medium  4/24
m6.large  8/48
m6.xlarge  16/96


####m8
m8.tiny  1/8
m8.small  2/16
m8.medium  4/32
m8.large  8/64
m8.xlarge  16/128


####m12
m12.large  8/96


####m16
m16.tiny  1/16
m16.small  2/32
m16.large  8/128
m16.xxlarge  32/128


#"os_name": "rhel" or "altlinux"
#"os_version": "7.9" or "altlinux"

# storage_type = "rbd-1"           ------> SLOW
# storage_type = "rbd-2"           ------> SLOW BACKUP
# storage_type = "iscsi_common"    ------> FAST
# storage_type = "__DEFAULT__"     ------> DEFAULT TYPE
*/


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
