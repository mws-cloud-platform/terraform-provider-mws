resource "mws_secretmanager_secret" "secret" {
  name   = "%s"
  active = true
  encryption = {
    crypto_key_id = "%s"
  }
}
