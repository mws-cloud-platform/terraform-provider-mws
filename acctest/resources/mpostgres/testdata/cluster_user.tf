resource "mws_mpostgres_cluster_user" "user" {
  cluster          = "%s"
  user             = "%s"
  password         = "%s"
  password_version = 1
  role             = "DB_OWNER_USER"
  additional_roles = [
    {
      name = "DB_MIGRATOR_ROLE"
      expires_at = "%s"
    }
  ]
}
