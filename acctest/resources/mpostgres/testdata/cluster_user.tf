resource "mws_mpostgres_cluster_user" "user" {
  cluster          = "%s"
  user             = "%s"
  password         = "%s"
  password_version = 1
  role             = "DB_OWNER_USER"
}
