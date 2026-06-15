resource "mws_mkafka_cluster" "kafka" {
  cluster = "{{.Name}}"
  version = "4.0"

  metadata = {
    display_name = "Test Kafka Cluster"
    description  = "Test Kafka cluster for acctest"
  }

  active = true

  endpoints = [
    {
      name    = "vpc-endpoint"
      network = "{{.NetworkID}}"
      broker_addresses = [
        { ref = "{{.Broker1AddressRef}}" }
      ]
    }
  ]

  instances = {
    broker = {
      vm_type = "compute/vmTypes/gen-2-8"
      disk = {
        size = "10Gb"
        type = "NETWORK_STANDARD_SSD"
      }
      allocation = [
        {
          zone  = "ru-central1-a"
          count = 1
        }
      ]
    }
    controller = {
      combined_with_broker = true
    }
  }

  maintenance_window = {
    weekly = {
      days = ["TUESDAY"]
      hour = 4
    }
  }
}
