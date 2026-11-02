import {
  to = mws_compute_disk.imported
  identity = {
    disk = "%[1]s"
  }
}

resource "mws_compute_disk" "imported" {
  disk      = "%[1]s"
  disk_type = "diskTypes/nbs-pl2"
  iops      = 1000
  size      = "1GB"
}
