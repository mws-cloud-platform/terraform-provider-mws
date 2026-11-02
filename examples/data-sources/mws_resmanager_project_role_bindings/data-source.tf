data "mws_resmanager_project_role_bindings" "all" {
}

output "role_bindings" {
  value = data.mws_resmanager_project_role_bindings.all.items.*.metadata.id
}
