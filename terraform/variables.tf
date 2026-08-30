variable "project_name" {
  type    = string
  default = "gobank"
}

variable "location" {
  type    = string
  default = "eastus"
}

variable "node_count" {
  type    = number
  default = 1
}

variable "node_vm_size" {
  type    = string
  default = "Standard_B2s"
}

variable "environment" {
  type    = string
  default = "dev"
}
