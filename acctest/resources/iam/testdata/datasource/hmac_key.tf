data "mws_iam_hmac_key" "hmac_key" {
  key_name        = "%s"
  service_account = "%s"
}
