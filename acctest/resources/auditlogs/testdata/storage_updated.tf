resource "mws_auditlogs_storage" "storage" {
  storage_name = "%s"

  metadata = {
    display_name = "Test Audit Logs Storage Updated"
  }

  retention = {
    period = "60d"
    size   = "20GB"
  }
}
