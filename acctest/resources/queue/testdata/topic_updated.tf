resource "mws_queue_topic" "topic" {
  topic = "%s"

  metadata = {
    display_name = "Test Queue Topic"
    description  = "Test topic for acctest"
  }

  partition_count = 6

  config = {
    "retention.ms"   = "604800000"
    "cleanup.policy" = "compact"
  }
}
