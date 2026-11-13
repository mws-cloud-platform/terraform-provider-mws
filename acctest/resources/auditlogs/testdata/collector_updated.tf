resource "mws_auditlogs_collector" "collector" {
  collector_name = "%s"

  metadata = {
    display_name = "Test Audit Logs Collector Updated"
  }

  destination = {
    storage = "%s"
  }

  enabled = false

  events_sources = [
    { project = "%s" }
  ]

  services = [
    {
      service = "compute"
      filter  = {}
    },
    {
      service = "iam"
      filter  = {}
    },
  ]

  service_account = "%s"
}
