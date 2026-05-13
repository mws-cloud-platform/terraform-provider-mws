data "mws_iam_hmac_key" "iam_hmac_key" {
  key_name        = "my-hmac-key"
  service_account = "my-service-account"
}
