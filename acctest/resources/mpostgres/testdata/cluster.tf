resource "mws_mpostgres_cluster" "cluster" {
  cluster = "%s"
  version = "17"
  active  = true

  metadata = {
    display_name = "Standalone Postgres Cluster"
    description  = "A standalone PostgreSQL cluster"
  }

  endpoints = [
    {
      name    = "primary-endpoint"
      network = "%s"
      primary_addresses = [
        {
          ref = "%s"
        }
      ]
    }
  ]

  instance_template = {
    vm_type = "vmTypes/gen-2-8"
    disk = {
      size = "20GB"
      type = "NETWORK_STANDARD_SSD"
    }
  }

  instances = [
    {
      count = 1
      zone  = "ru-central1-a"
    }
  ]

  backup = {
    retain_period_days = 7
    daily = {
      hour = 2
    }
  }

  maintenance_window = {
    weekly = {
      days = ["SUNDAY"]
      hour = 4
    }
  }
}
