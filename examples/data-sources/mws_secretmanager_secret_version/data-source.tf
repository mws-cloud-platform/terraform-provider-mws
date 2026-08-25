data "mws_secretmanager_secret_version" "secretmanager_secret_version" {
  name    = "my-secret"
  version = "current"
}
