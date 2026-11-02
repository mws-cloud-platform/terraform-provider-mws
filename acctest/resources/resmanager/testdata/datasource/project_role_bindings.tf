data "mws_resmanager_project_role_bindings" "filtered" {
  filter = "metadata.id.name = \"%s\""
}
