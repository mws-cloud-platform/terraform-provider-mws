resource "mws_queue_topic" "topic" {
  topic = "%s"

  metadata = {
    display_name = "Test Queue Topic"
    description  = "Test topic for acctest"
  }

  partition_count = 3

  config = {
    "retention.ms" = "604800000"
  }
}
