data "mws_resmanager_zones" "all" {
}

output "zones" {
  value = data.mws_resmanager_zones.all.items.*.metadata.id
}
