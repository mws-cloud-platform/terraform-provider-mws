resource "mws_auditlogs_storage" "example" {
  storage_name = var.storage_name

  metadata = {
    display_name = "Example Audit Logs Storage"
  }

  retention = {
    period = "30d"
    size   = "10GB"
  }
}

variable "storage_name" {
  type        = string
  default     = "example-storage"
  description = "Audit logs storage name"
}
