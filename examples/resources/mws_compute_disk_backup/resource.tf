resource "mws_compute_disk" "source_disk" {
  disk      = var.source_disk_name
  disk_type = var.disk_type
  iops      = 1000
  size      = "10GB"
}

resource "mws_compute_disk_backup" "backup" {
  disk_backup = var.backup_name

  source = {
    disk = {
      id = mws_compute_disk.source_disk.id
    }
  }

  metadata = {
    description  = "Backup of test disk"
    display_name = "Test Disk Backup"
  }
}

variable "source_disk_name" {
  type        = string
  default     = "test-source-disk"
  description = "Source disk name"
}

variable "backup_name" {
  type        = string
  default     = "test-disk-backup"
  description = "Backup name"
}

variable "disk_type" {
  type        = string
  default     = "diskTypes/nbs-pl2"
  description = "Disk type"
}
