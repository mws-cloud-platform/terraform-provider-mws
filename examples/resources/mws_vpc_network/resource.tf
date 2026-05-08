resource "mws_vpc_network" "network" {
  network = var.network_name
}

variable "network_name" {
  type        = string
  default     = "my-test-network"
  description = "Network name"
}
