data "mws_secretmanager_secret_version" "secret_version" {
  name    = "%s"
  version = "current"
}
