resource "mws_resmanager_project_role_binding" "project_role_binding" {
  role_binding = "%s"
  role         = "iam/roles/admin"

  subject = {
    service_account = "%s"
  }
}
