resource "mws_resmanager_enabled_service" "service" {
  service = "%s"

  timeouts = {
    create = "1h"
  }
}
