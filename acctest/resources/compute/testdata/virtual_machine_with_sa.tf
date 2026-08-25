resource "mws_compute_virtual_machine" "vm" {
  virtual_machine = "%[1]s"
  vm_type         = "vmTypes/gen-2-4"

  hardware = {
    power = "OFF"
  }
  storage = {
    disks = [
      {
        name = "boot"
        boot = true
        disk = {
          spec = {
            diskType = "nbs-pl2"
            size = "10 GB"
            iops = 1000
            source = {
              image = "%[2]s"
            }
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
              spec = {
                subnet = "%[3]s"
              }
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

  service_account = %[4]s
}
