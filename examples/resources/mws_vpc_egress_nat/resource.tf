resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet" {
  subnet  = var.subnet_name
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr
}

resource "mws_vpc_external_address" "external_address" {
  external_address = var.external_address_name
}

resource "mws_vpc_egress_nat" "example_egress_nat" {
  egress_nat = "example-egress-nat"
  network    = mws_vpc_network.network.network

  metadata = {
    description  = "This is an example Egress NAT resource"
    display_name = "Example Egress NAT"
  }

  external = {
    addresses = [
      {
        ref = mws_vpc_external_address.external_address.metadata.id
      }
    ]
  }

  internal = {
    subnets = [mws_vpc_subnet.subnet.metadata.id]
  }
}

variable "network_name" {
  type        = string
  default     = "my-network"
  description = "VPC network name"
}

variable "subnet_name" {
  type        = string
  default     = "my-subnet"
  description = "Base name for subnets"
}

variable "subnet_cidr" {
  type        = string
  default     = "192.168.0.0/17"
  description = "CIDR for subnet"
}

variable "external_address_name" {
  type        = string
  default     = "my-test-external-address"
  description = "External address name"
}
