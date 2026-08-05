resource "mws_iam_authorized_key" "authorized_key" {
  authorized_key  = "%s"
  service_account = "%s"
  key_algorithm   = "ES256"

  active          = true
  expiration_time = "2027-01-01T00:00:00Z"
}
