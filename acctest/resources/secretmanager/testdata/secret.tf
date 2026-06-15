resource "mws_secretmanager_secret" "secret" {
  name   = "%s"
  active = true
}
