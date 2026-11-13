resource "mws_auditlogs_collector" "collector" {
  collector_name = "%s"

  metadata = {
    display_name = "Test Audit Logs Collector"
  }

  destination = {
    storage = "%s"
  }

  enabled = true

  events_sources = [
    { project = "%s" }
  ]

  services = [
    {
      service = "compute"
      filter  = {}
    },
  ]

  service_account = "%s"
}
