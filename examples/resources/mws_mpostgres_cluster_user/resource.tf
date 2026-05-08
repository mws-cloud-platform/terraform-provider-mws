resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet" {
  subnet  = var.subnet_name
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr
}

resource "mws_vpc_address" "cluster_primary_address" {
  address = "${var.cluster_name}-primary-address"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet.metadata.id
}

resource "mws_mpostgres_cluster" "cluster" {
  cluster = var.cluster_name
  version = "17"
  active  = true

  metadata = {
    display_name = "Standalone Postgres Cluster"
    description  = "A standalone PostgreSQL cluster"
  }

  endpoints = [
    {
      name    = "primary-endpoint"
      network = mws_vpc_network.network.metadata.id
      primary_addresses = [
        {
          ref = mws_vpc_address.cluster_primary_address.metadata.id
        }
      ]
    }
  ]

  instance_template = {
    vm_type = "vmTypes/gen-2-8"
    disk = {
      size = "20GB"
      type = "NETWORK_STANDARD_SSD"
    }
  }

  instances = [
    {
      count = 1
      zone  = "ru-central1-a"
    }
  ]
}

resource "mws_mpostgres_cluster_user" "user" {
  cluster          = mws_mpostgres_cluster.cluster.cluster
  user             = var.user_name
  password         = var.user_password
  password_version = 1
  role             = "DB_OWNER_USER"
}

variable "cluster_name" {
  type        = string
  default     = "my-pg-cluster"
  description = "Postgres cluster name"
}

variable "user_name" {
  type        = string
  default     = "my-pg-cluster-user"
  description = "Database user name"
}

variable "user_password" {
  type        = string
  description = "Database user password"
  sensitive   = true
}

variable "network_name" {
  type        = string
  default     = "my-pg-cluster-network"
  description = "Network name"
}

variable "subnet_name" {
  type        = string
  default     = "my-pg-cluster-subnet"
  description = "Subnet name"
}

variable "subnet_cidr" {
  type        = string
  default     = "192.168.10.0/24"
  description = "Subnet CIDR"
}
