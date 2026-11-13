resource "mws_mkafka_connector" "connector" {
  cluster   = "%[1]s"
  connector = "%[2]s"

  metadata = {
    display_name = "Test Kafka Connector"
    description  = "Test connector for acctest"
  }

  s3_sink_connector = {
    topics    = "%[3]s"
    dlq_topic = "%[3]s-dlq"

    s3_properties = {
      bucket_name               = "acctest-kafka-connector"
      endpoint                  = "https://storage.mwsapis.ru"
      access_key_id             = "acctest-access-key-id"
      secret_access_key         = "acctest-secret-access-key"
      secret_access_key_version = 1
    }

    tasks_max = 1
  }
}
