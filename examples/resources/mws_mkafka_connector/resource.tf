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
  version = "4.0"

  metadata = {
    display_name = "Example Kafka Cluster"
    description  = "Managed Kafka cluster example"
  }

  active = true

  endpoints = [
    {
      name    = "vpc-endpoint"
      network = mws_vpc_network.network.metadata.id
      broker_addresses = [
        { ref = mws_vpc_address.broker_addr_1.id },
      ]
    }
  ]

  instances = {
    broker = {
      vm_type = "compute/vmTypes/gen-2-8"
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
      combined_with_broker = true
    }
  }
}

resource "mws_mkafka_topic" "example" {
  cluster = mws_mkafka_cluster.example.cluster
  topic   = var.topic_name

  partitions         = 3
  replication_factor = 1
}

resource "mws_mkafka_topic" "dlq" {
  cluster = mws_mkafka_cluster.example.cluster
  topic   = "${var.topic_name}-dlq"

  partitions         = 1
  replication_factor = 1
}

resource "mws_iam_service_account" "example" {
  service_account = var.service_account_name
}

resource "mws_iam_hmac_key" "example" {
  key_name        = var.hmac_key_name
  service_account = mws_iam_service_account.example.service_account
}

resource "mws_mkafka_connector" "example" {
  cluster   = mws_mkafka_cluster.example.cluster
  connector = var.connector_name

  metadata = {
    display_name = "Example Kafka Connector"
    description  = "S3 Sink connector for topic backup"
  }

  active = true

  s3_sink_connector = {
    topics    = mws_mkafka_topic.example.topic
    dlq_topic = mws_mkafka_topic.dlq.topic

    s3_properties = {
      bucket_name               = var.s3_bucket
      endpoint                  = var.s3_endpoint
      access_key_id             = mws_iam_hmac_key.example.status.access_key_id
      secret_access_key         = mws_iam_hmac_key.example.status.secret_access_key
      secret_access_key_version = 1
    }

    tasks_max             = 1
    input_format          = "JSON"
    output_format         = "JSONL"
    file_compression_type = "GZIP"
    output_fields         = ["KEY", "VALUE", "TIMESTAMP"]
    output_envelope       = true
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

variable "kafka_name" {
  type        = string
  default     = "kafka-cluster"
  description = "Kafka cluster name"
}

variable "topic_name" {
  type        = string
  default     = "example-topic"
  description = "Kafka topic name"
}

variable "service_account_name" {
  type        = string
  default     = "kafka-connector-sa"
  description = "Service account name for S3 access"
}

variable "hmac_key_name" {
  type        = string
  default     = "kafka-connector-hmac-key"
  description = "HMAC key name for S3 access"
}

variable "connector_name" {
  type        = string
  default     = "example-connector"
  description = "Kafka connector name"
}

variable "s3_bucket" {
  type        = string
  default     = "kafka-connector-backup"
  description = "Existing S3 bucket for topic backup"
}

variable "s3_endpoint" {
  type        = string
  default     = "https://storage.mwsapis.ru"
  description = "S3 endpoint"
}
