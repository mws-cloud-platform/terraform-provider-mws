resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet" {
  subnet  = var.subnet_name
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr
}

resource "mws_mclickhouse_cluster" "example" {
  cluster = var.cluster_name
  version = "25.3"
  active  = true

  metadata = {
    display_name = "Standalone ClickHouse Cluster Example"
    description  = "Standalone ClickHouse Cluster"
  }

  shards = [{
    name = "shard"

    resources = {
      vm_type = "vmTypes/gen-4-8"
      disk = {
        type = "NETWORK_STANDARD_SSD"
        size = "10Gb"
      }
    }

    weight = 1

    instances = [{
      name  = "instance-1"
      zone  = "ru-central1-a"
      count = 1

      endpoints = [{
        address = {
          spec = {
            subnet = mws_vpc_subnet.subnet.metadata.id
          }
        }
        external_address = {
          spec = {}
        }
      }]
    }]
  }]

  bootstrap_admin = {
    username         = var.cluster_admin_username
    password_version = 1
    password         = var.cluster_admin_password
  }

  maintenance_window = {
    weekly = {
      days = ["MONDAY"]
      hour = 3
    }
  }

  backup = {
    hour               = 2
    retain_period_days = 7
  }
}

variable "network_name" {
  type        = string
  default     = "standalone-ch-network"
  description = "Network name"
}

variable "subnet_name" {
  type        = string
  default     = "standalone-ch-subnet"
  description = "Subnet name"
}

variable "subnet_cidr" {
  type        = string
  default     = "192.168.1.0/24"
  description = "Subnet CIDR"
}

variable "cluster_name" {
  type        = string
  default     = "standalone-ch"
  description = "ClickHouse cluster name"
}

variable "cluster_admin_username" {
  type        = string
  default     = "admin"
  description = "ClickHouse cluster admin username"
}

variable "cluster_admin_password" {
  type        = string
  sensitive   = true
  default     = "securePassword123!"
  description = "ClickHouse cluster admin password"
}
