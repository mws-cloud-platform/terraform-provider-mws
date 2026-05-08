resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet_a" {
  subnet  = "${var.subnet_name}-a"
  network = mws_vpc_network.network.network
  cidr    = var.subnet_cidr_a
}

resource "mws_vpc_address" "broker_addr_1" {
  address = "${var.kafka_name}-broker-addr-1"
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet_a.metadata.id
}

resource "mws_mkafka_cluster" "example" {
  cluster = var.kafka_name
  version = "3.6.0-mws.1"

  metadata = {
    display_name = "Example Kafka Cluster"
    description  = "Managed Kafka cluster example with VPC endpoints"
  }

  active = true

  endpoints = [
    {
      name    = "vpc-endpoint"
      network = mws_vpc_network.network.network
      broker_addresses = [
        { ref = mws_vpc_address.broker_addr_1.id },
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
          zone  = "ru-central1-a"
          count = 1
        }
      ]
    }
    controller = {
      vm_type = "compute/vmTypes/gen-2-4"
      disk = {
        size = "10Gb"
        type = "NETWORK_STANDARD_SSD"
      }
      allocation = [
        {
          zone  = "ru-central1-a"
          count = 1
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
}

resource "mws_mkafka_cluster_user" "example_user" {
  cluster          = mws_mkafka_cluster.example.cluster
  user             = var.kafka_user_name
  password         = var.kafka_user_password
  password_version = 1 //increase on change password

  metadata = {
    display_name = "Example Kafka User"
    description  = "User for accessing the example Kafka cluster"
  }

  roles = [
    {
      name = "CLUSTER_ADMIN"
    }
  ]
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

variable "kafka_name" {
  type        = string
  default     = "kafka-cluster"
  description = "Kafka cluster name"
}

variable "kafka_user_name" {
  type        = string
  default     = "example-user"
  description = "Kafka user name"
}

variable "kafka_user_password" {
  type        = string
  default     = "securePassword123!"
  description = "Kafka user password"
}
