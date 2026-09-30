resource "mws_queue_topic" "example" {
  topic = var.topic_name

  metadata = {
    display_name = "Example Serverless Queue Topic"
    description  = "Serverless Queue topic example"
  }

  partition_count = 3

  config = {
    "cleanup.policy" = "compact"
    "retention.ms"   = "86400000"
  }
}

variable "topic_name" {
  type        = string
  default     = "example-topic"
  description = "Queue topic name"
}
