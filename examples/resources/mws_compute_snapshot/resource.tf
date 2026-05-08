data "mws_compute_image" "image" {
  image   = "mws-ubuntu-2204-lts-v20250529"
  project = "mws-ubuntu"
}

resource "mws_compute_disk" "disk" {
  disk      = var.disk_name
  disk_type = "diskTypes/nbs-pl2"
  iops      = 1000
  size      = "10GB"
  source = {
    image = data.mws_compute_image.image.metadata.id
  }
}

resource "mws_compute_snapshot" "snapshot" {
  snapshot = var.snapshot_name
  source = {
    disk = {
      id = mws_compute_disk.disk.metadata.id
    }
  }
}

variable "snapshot_name" {
  type        = string
  default     = "my-test-snapshot"
  description = "Snapshot name"
}

variable "disk_name" {
  type        = string
  default     = "my-test-disk"
  description = "Disk name"
}
