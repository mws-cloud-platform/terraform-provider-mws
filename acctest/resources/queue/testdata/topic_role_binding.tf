resource "mws_queue_topic_role_binding" "topic_role_binding" {
  topic        = "%s"
  role_binding = "%s"
  role         = "iam/roles/queue.consumer"

  subject = {
    service_account = "%s"
  }

  metadata = {
    display_name = "Test Queue Topic Role Binding"
    description  = "Test role binding for acctest"
  }
}
