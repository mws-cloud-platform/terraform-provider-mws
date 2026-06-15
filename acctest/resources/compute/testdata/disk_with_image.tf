resource "mws_compute_disk" "disk" {
  disk      = "%s"
  disk_type = "diskTypes/nbs-pl2"
  iops      = 1000
  size      = "%s"
  source = {
    image = "%s"
  }
}
