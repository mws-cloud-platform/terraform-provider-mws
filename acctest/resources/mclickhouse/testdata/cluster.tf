resource "mws_mclickhouse_cluster" "cluster" {
  cluster = "%s"
  version = "25.3"
  active  = true

  metadata = {
    display_name = "Standalone ClickHouse Cluster"
    description  = "A standalone ClickHouse Cluster"
  }

  shards = [{
    name = "shard"

    resources = {
      vm_type = "vmTypes/gen-4-8"
      disk = {
        type = "NETWORK_STANDARD_SSD"
        size = "10GB"
      }
    }

    instances = [{
      name  = "instance-1"
      zone  = "ru-central1-a"
      count = 1

      endpoints = [{
        address = {
          spec = {
            subnet = "%s"
          }
        }
        external_address = {
          spec = {}
        }
      }]
    }]
  }]

  bootstrap_admin = {
    username         = "admin"
    password_version = 1
    password         = "%s"
  }

  maintenance_window = {
    weekly = {
      days = ["MONDAY"]
      hour = 3
    }
  }

  backup = {
    hour               = 2
    retain_period_days = 7
  }
}
