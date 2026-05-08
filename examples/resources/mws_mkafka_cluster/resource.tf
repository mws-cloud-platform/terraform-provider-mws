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

resource "mws_vpc_address" "broker_addr_1" {
  address = "${var.kafka_name}-broker-addr-1"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet_a.metadata.id
}

resource "mws_vpc_address" "broker_addr_2" {
  address = "${var.kafka_name}-broker-addr-2"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet_b.metadata.id
}

resource "mws_vpc_address" "broker_addr_3" {
  address = "${var.kafka_name}-broker-addr-3"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet_b.metadata.id
}

resource "mws_mkafka_cluster" "example" {
  cluster = var.kafka_name
  version = "4.0"

  metadata = {
    display_name = "Example Kafka Cluster"
    description  = "Managed Kafka cluster example with VPC endpoints"
  }

  active = true

  endpoints = [
    {
      name    = "vpc-endpoint"
      network = mws_vpc_network.network.metadata.id
      broker_addresses = [
        { ref = mws_vpc_address.broker_addr_1.id },
        { ref = mws_vpc_address.broker_addr_2.id },
        { ref = mws_vpc_address.broker_addr_3.id }
      ]
    }
  ]

  instances = {
    broker = {
      vm_type = "compute/vmTypes/gen-2-4"
      disk = {
        size = "10Gb"
        type = "NETWORK_STANDARD_SSD"
      }
      allocation = [
        {
          zone  = "ru-central1-b"
          count = 3
        }
      ]
    }
    controller = {
      combined_with_broker = false
      vm_type              = "compute/vmTypes/gen-2-4"
      disk = {
        size = "10Gb"
        type = "NETWORK_STANDARD_SSD"
      }
      allocation = [
        {
          zone  = "ru-central1-b"
          count = 3
        }
      ]
    }
  }

  maintenance_window = {
    weekly = {
      days = ["TUESDAY"]
      hour = 4
    }
  }
  timeouts = {
    create = "1h"
    update = "1h"
    delete = "1h"
  }
}

variable "network_name" {
  type        = string
  default     = "kafka-vpc-network"
  description = "VPC network name"
}

variable "subnet_name" {
  type        = string
  default     = "kafka-subnet"
  description = "Base name for subnets"
}

variable "subnet_cidr_a" {
  type        = string
  default     = "192.168.1.0/24"
  description = "CIDR for subnet A"
}

variable "subnet_cidr_b" {
  type        = string
  default     = "192.168.2.0/24"
  description = "CIDR for subnet B"
}

variable "kafka_name" {
  type        = string
  default     = "kafka-cluster"
  description = "Kafka cluster name"
}
