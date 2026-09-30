resource "mws_iam_service_account" "reader" {
  service_account = var.service_account_name

  metadata = {
    display_name = "Example Queue Topic Reader"
  }
}

resource "mws_queue_topic" "example" {
  topic = var.topic_name

  metadata = {
    display_name = "Example Queue Topic"
    description  = "Managed Queue topic example"
  }

  partition_count = 3

  config = {
    "cleanup.policy" = "compact"
    "retention.ms"   = "86400000"
  }
}

resource "mws_queue_topic_role_binding" "example" {
  topic        = mws_queue_topic.example.topic
  role_binding = var.role_binding_name
  role         = "iam/roles/queue.consumer"

  subject = {
    service_account = mws_iam_service_account.reader.metadata.id
  }

  metadata = {
    display_name = "Example Queue Topic Role Binding"
    description  = "Grants the service account read access to the topic"
  }
}

variable "service_account_name" {
  type        = string
  default     = "example-topic-reader"
  description = "Service account name"
}

variable "topic_name" {
  type        = string
  default     = "example-topic"
  description = "Queue topic name"
}

variable "role_binding_name" {
  type        = string
  default     = "example-role-binding"
  description = "Queue topic role binding name"
}
