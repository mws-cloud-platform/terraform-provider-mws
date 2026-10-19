data "mws_resmanager_regions" "all" {
}

output "regions" {
  value = data.mws_resmanager_regions.all.items.*.metadata.id
}
