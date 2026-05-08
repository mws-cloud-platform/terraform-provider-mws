data "mws_compute_image" "image" {
  image   = "mws-ubuntu-2204-lts-v20250529"
  project = "mws-ubuntu"
}

resource "mws_vpc_network" "network" {
  network = var.network_name
}

resource "mws_vpc_subnet" "subnet" {
  subnet  = var.subnet_name
  network = mws_vpc_network.network.network
  cidr    = "192.168.0.0/16"
}

resource "mws_vpc_address" "vm_primary_network_interface_address" {
  network = mws_vpc_network.network.network
  subnet  = mws_vpc_subnet.subnet.metadata.id
  address = "${var.vm_name}-primary-network-interface-address"
}

resource "mws_vpc_external_address" "vm_external_address" {
  external_address = "${var.vm_name}-external-address"
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

resource "mws_compute_virtual_machine" "vm" {
  virtual_machine = var.vm_name
  vm_type         = "vmTypes/gen-2-8"

  hardware = {
    power                     = "ON"
    graceful_shutdown_timeout = "1m 30s"
  }

  storage = {
    disks = [
      {
        name = "boot"
        boot = true
        disk = {
          ref = mws_compute_disk.disk.metadata.id
        }
      }
    ]
  }

  network = {
    network_interfaces = [
      {
        name    = "${var.vm_name}-network-interface-primary"
        primary = true
        addresses = [
          {
            address = {
              ref = mws_vpc_address.vm_primary_network_interface_address.metadata.id
            }
            one_to_one_nat = {
              external = {
                address = {
                  ref = mws_vpc_external_address.vm_external_address.metadata.id
                }
              }
            }
          }
        ]
      }
    ]
  }
}

variable "vm_name" {
  type        = string
  default     = "my-test-vm"
  description = "Virtual machine name"
}

variable "network_name" {
  type        = string
  default     = "my-test-network"
  description = "Network name"
}

variable "subnet_name" {
  type        = string
  default     = "my-test-subnet"
  description = "Subnet name"
}

variable "disk_name" {
  type        = string
  default     = "my-test-disk"
  description = "Disk name"
}
