resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet" {
  subnet  = var.subnet_name
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr
}

resource "mws_vpc_address" "address" {
  address = var.address_name
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet.metadata.id
}

resource "mws_vpc_route" "route" {
  route   = var.route_name
  network = mws_vpc_network.network.network

  destination = {
    spec = {
      cidrs = var.destination_cidrs
    }
  }

  next_hop = {
    address = {
      ref = mws_vpc_address.address.id
    }
  }

  metadata = {
    description  = "Route to external network"
    display_name = "External Route"
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

variable "address_name" {
  type        = string
  default     = "my-test-address"
  description = "Address name"
}

variable "route_name" {
  type        = string
  default     = "my-test-route"
  description = "Route name"
}

variable "destination_cidrs" {
  type        = list(string)
  default     = ["10.0.0.0/8"]
  description = "Destination CIDRs for the route"
}