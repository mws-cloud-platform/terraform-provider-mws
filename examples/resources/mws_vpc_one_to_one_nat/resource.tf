resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet" {
  subnet  = var.subnet_name
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr
}

resource "mws_vpc_address" "internal_address" {
  address = var.internal_address_name
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet.metadata.id
}

resource "mws_vpc_external_address" "address" {
  external_address = var.external_address_name
}

resource "mws_vpc_one_to_one_nat" "example" {
  network        = mws_vpc_network.network.network
  one_to_one_nat = var.one_to_one_name

  external = {
    address = {
      ref = mws_vpc_external_address.address.metadata.id
    }
  }

  internal = {
    address = {
      ref = mws_vpc_address.internal_address.metadata.id
    }
  }
}

variable "network_name" {
  type        = string
  default     = "my-test-network"
  description = "Network name"
}

variable "subnet_name" {
  type        = string
  default     = "my-test-subnet"
  description = "Subnet name"
}

variable "subnet_cidr" {
  type        = string
  default     = "192.168.0.0/16"
  description = "Subnet CIDR"
}

variable "internal_address_name" {
  type        = string
  default     = "my-test-internal-address"
  description = "Internal address name"
}

variable "external_address_name" {
  type        = string
  default     = "my-test-external-address"
  description = "External address name"
}

variable "one_to_one_name" {
  type        = string
  default     = "my-test-one-to-one"
  description = "One-to-One NAT name"
}

