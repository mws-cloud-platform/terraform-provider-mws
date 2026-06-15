resource "mws_mkafka_topic" "topic" {
  cluster = "{{.KafkaName}}"
  topic   = "{{.Name}}"

  metadata = {
    display_name = "Test Kafka Topic"
    description  = "Test topic for acctest"
  }

  partitions         = 3
  replication_factor = 1
}
