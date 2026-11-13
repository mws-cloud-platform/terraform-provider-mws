resource "mws_auditlogs_storage" "storage" {
  storage_name = "%s"

  metadata = {
    display_name = "Test Audit Logs Storage"
  }

  retention = {
    period = "30d"
    size   = "10GB"
  }
}
