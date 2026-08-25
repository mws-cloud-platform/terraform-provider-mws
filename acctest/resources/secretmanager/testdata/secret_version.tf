resource "mws_secretmanager_secret_version" "secret_version" {
  name    = "%s"
  active  = true

  data = {
    foo = "bar"
  }
}
