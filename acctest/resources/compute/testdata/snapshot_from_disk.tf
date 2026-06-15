resource "mws_compute_snapshot" "snapshot" {
  snapshot = "%s"
  source = {
    disk = {
      id = "%s"
    }
  }
}
