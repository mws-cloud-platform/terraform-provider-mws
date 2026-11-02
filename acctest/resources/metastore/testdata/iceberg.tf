resource "mws_metastore_iceberg" "iceberg" {
  catalog_name = "%[1]s"

  service_account = {
    ref = "%[2]s"
  }

  properties = {
    location = "s3://%[3]s/%[1]s"
  }
}
