resource "mws_mpostgres_cluster_database" "database" {
  cluster             = "%s"
  database            = "%s"
  owner               = "%s"
  deletion_protection = false
}
