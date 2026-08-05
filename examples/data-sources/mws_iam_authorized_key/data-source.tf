data "mws_iam_authorized_key" "iam_authorized_key" {
  authorized_key  = "my-authorized-key"
  service_account = "my-service-account"
}
