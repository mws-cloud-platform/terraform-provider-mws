resource "mws_compute_disk_backup" "backup" {
  disk_backup = "%s"

  source = {
    disk = {
      id = "%s"
    }
  }
}
