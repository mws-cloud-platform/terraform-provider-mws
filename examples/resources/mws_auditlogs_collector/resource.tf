resource "mws_auditlogs_collector" "example" {
  collector_name = var.collector_name

  metadata = {
    display_name = "Example Audit Logs Collector"
  }

  destination = {
    storage = mws_auditlogs_storage.example.id
  }

  enabled = true

  events_sources = []

  services = [
    {
      service = "compute"
      filter  = {}
    },
  ]
}

resource "mws_auditlogs_storage" "example" {
  storage_name = var.storage_name

  retention = {
    period = "30d"
  }
}

variable "collector_name" {
  type        = string
  default     = "example-collector"
  description = "Audit logs collector name"
}

variable "storage_name" {
  type        = string
  default     = "example-storage"
  description = "Audit logs storage name"
}
