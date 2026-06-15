resource "mws_compute_virtual_machine" "vm" {
  virtual_machine = "%[1]s"
  vm_type         = "vmTypes/gen-2-8"

  hardware = {
    power                     = "ON"
    graceful_shutdown_timeout = "1m 30s"
  }

  os = {
    hostname     = "my"
    local_domain = "vm"
    metadata = {
      attributes = {
        env       = "dev"
        user-data = "string"
        test      = "test"
      }
    }
  }

  storage = {
    disks = [
      {
        name = "boot"
        boot = true
        disk = {
          ref = "%[2]s"
        }
      },
      {
        name = "data"
        boot = false
        disk = {
          ref = "%[3]s"
        }
      },
      {
        name = "other"
        boot = "false" # quoted bool is also valid here
        disk = {
          spec = {
            zone = "ru-central1-a"
            disk_type = "nbs-pl2"
            size = "10 GB"
            iops = "1000"
            block_size = "4 KB"
          }
        }
      }
    ]
  }

  network = {
    network_interfaces = [
      {
        name    = "%[1]s-network-interface-primary"
        primary = true
        addresses = [
          {
            address = {
              ref = "%[4]s"
            }
            one_to_one_nat = {
              external = {
                address = {
                  spec = {}
                }
              }
            }
          }
        ]
      }
    ]
  }
}
