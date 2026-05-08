resource "mws_vpc_external_address" "vm_external_address" {
  external_address = var.external_address_name
}

variable "external_address_name" {
  type        = string
  default     = "my-test-external-address"
  description = "External address name"
}
