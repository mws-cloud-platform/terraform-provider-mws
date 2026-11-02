resource "mws_metastore_iceberg_role_binding" "role_binding" {
  role_binding = "%[2]s"
  role         = "metastore.icebergViewer"
  subject      = "%[3]s"

  resource = {
    catalog = "%[1]s"
  }
}
