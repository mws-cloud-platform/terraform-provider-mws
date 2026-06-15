resource "mws_mkafka_cluster_user" "user" {
  cluster  = "{{.KafkaName}}"
  user     = "{{.Name}}"
  password = "{{.AdminPassword}}"
  password_version = 1

  roles = [
    {
      name = "CLUSTER_ADMIN"
    }
  ]
}
