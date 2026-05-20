resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet_a" {
  subnet  = "${var.subnet_name}-a"
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr_a
}

resource "mws_vpc_subnet" "subnet_b" {
  subnet  = "${var.subnet_name}-b"
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr_b
}

resource "mws_vpc_address" "address" {
  address = var.address_name
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet_a.metadata.id
}

resource "mws_iam_service_account" "sa" {
  service_account = var.service_account_name
}

resource "mws_mk8s_cluster" "example" {
  availability = {
    zonal_ha = { // zonal high available
      zone = "ru-central1-b"
    }
  }

  cluster_name = var.cluster_name

  metadata = {
    description  = "Zonal HA Kubernetes cluster"
    display_name = "HA Example Cluster"
  }

  network = {
    pods_cidr     = var.pods_cidr
    services_cidr = var.services_cidr
    primary_endpoint = {
      ref = mws_vpc_address.address.id
    }
  }

  version_control = {
    release_channel = "stable"
    maintenance_window = {
      weekly = {
        days = ["MONDAY", "WEDNESDAY"]
        hour = 3
      }
    }
    version = "v1.34.1-mws.1"
  }
}

resource "mws_mk8s_node_group" "example" {
  cluster_name    = mws_mk8s_cluster.example.cluster_name
  node_group_name = "${var.cluster_name}-example-node-group"

  service_account = {
    ref = mws_iam_service_account.sa.metadata.id
  }

  subnet = {
    ref = mws_vpc_subnet.subnet_a.metadata.id
  }

  vm_type = {
    ref = "compute/vmTypes/gen-2-8"
  }

  scale = {
    autoscaling = {
      min = 1
      max = 3
    }
  }

  rollout_strategy = {
    max_surge       = 1
    max_unavailable = 0
  }

  version_control = {
    auto_update = true
    version     = "v1.34.1-mws.1"
    maintenance_window = {
      weekly = {
        days     = ["MONDAY", "WEDNESDAY"]
        hour     = 3
        duration = "4h"
      }
    }
  }

  image_storage_size = "20Gb"

  labels = [
    { key = "environment", value = "production" },
    { key = "team", value = "devops" }
  ]

  taints = []

  zone = "ru-central1-a"
}

variable "network_name" {
  type        = string
  default     = "my-ha-network"
  description = "VPC network name"
}

variable "subnet_name" {
  type        = string
  default     = "my-ha-subnet"
  description = "Base name for subnets"
}

variable "subnet_cidr_a" {
  type        = string
  default     = "192.168.0.0/17"
  description = "CIDR for subnet A"
}

variable "subnet_cidr_b" {
  type        = string
  default     = "192.168.128.0/17"
  description = "CIDR for subnet B"
}

variable "address_name" {
  type        = string
  default     = "k8s-primary-ip"
  description = "Name for the primary endpoint address"
}

variable "service_account_name" {
  type        = string
  default     = "my-service-account"
  description = "Service account name"
}

variable "cluster_name" {
  type        = string
  default     = "ha-k8s-cluster"
  description = "Name of the Kubernetes cluster"
}

variable "pods_cidr" {
  type        = string
  default     = "10.244.0.0/16"
  description = "Pods CIDR block"
}

variable "services_cidr" {
  type        = string
  default     = "10.96.0.0/16"
  description = "Services CIDR block"
}
