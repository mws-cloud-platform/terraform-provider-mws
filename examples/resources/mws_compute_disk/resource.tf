resource "mws_compute_disk" "disk" {
  disk      = var.disk_name
  disk_type = var.disk_type
  iops      = 1000
  size      = "10GB"
}

variable "disk_name" {
  type        = string
  default     = "my-test-disk"
  description = "Disk name"
}

variable "disk_type" {
  type        = string
  default     = "diskTypes/nbs-pl2"
  description = "Disk type"
}
