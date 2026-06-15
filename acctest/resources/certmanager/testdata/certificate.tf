resource "mws_certmanager_certificate" "certificate" {
  name = "%s"
  self_managed = {
    certificate = file("%s")
    private_key = file("%s")
  }
  self_managed_version = %d
}
