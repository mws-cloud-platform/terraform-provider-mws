resource "mws_resmanager_enabled_service" "compute" {
  service = "compute"

  timeouts = {
    create = "1h"
  }
}
