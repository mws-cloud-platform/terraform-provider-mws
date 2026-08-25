resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet" {
  subnet  = var.subnet_name
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr
}

resource "mws_vpc_address" "internal_address" {
  address = "${var.nlb_name}-internal"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet.metadata.id
}

resource "mws_vpc_address" "backend_address1" {
  address = "${var.nlb_name}-backend1"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet.metadata.id
}

resource "mws_vpc_address" "backend_address2" {
  address = "${var.nlb_name}-backend2"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet.metadata.id
}

resource "mws_vpc_address_group" "backend_group" {
  address_group = "${var.nlb_name}-backend-group"
  network       = mws_vpc_network.network.network

  metadata = {
    display_name = "Backend Address Group"
    description  = "Group of backend servers for NLB"
  }

  addresses = [
    {
      ref = mws_vpc_address.backend_address1.metadata.id
    },
    {
      ref = mws_vpc_address.backend_address2.metadata.id
    }
  ]
}

resource "mws_nlb_nlb" "example" {
  network = mws_vpc_network.network.network
  nlb     = var.nlb_name

  metadata = {
    display_name = "Example Internal Network Load Balancer"
    description  = "NLB example with internal listener and backend group"
  }

  listener = {
    internal = {
      address = {
        ref = mws_vpc_address.internal_address.metadata.id
      }
    }
  }

  rules = [
    {
      proto_port  = "TCP:80"
      target_port = 8080

      target_address_groups = [
        {
          ref = mws_vpc_address_group.backend_group.metadata.id
        }
      ]

      health_check = {
        protocol = {
          tcp = {
            port = 8080
          }
        }

        interval            = "30s"
        timeout             = "3s"
        healthy_threshold   = 3
        unhealthy_threshold = 3
      }
    }
  ]
}

variable "network_name" {
  type        = string
  default     = "nlb-example-network"
  description = "VPC network name"
}

variable "subnet_name" {
  type        = string
  default     = "nlb-example-subnet"
  description = "VPC subnet name"
}

variable "subnet_cidr" {
  type        = string
  default     = "192.168.0.0/16"
  description = "Subnet CIDR block"
}

variable "nlb_name" {
  type        = string
  default     = "example-nlb"
  description = "Network load balancer name"
}
