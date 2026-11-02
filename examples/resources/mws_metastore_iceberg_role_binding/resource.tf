resource "mws_iam_service_account" "example" {
  service_account = var.service_account_name
}

resource "mws_metastore_iceberg" "example" {
  catalog_name = var.catalog_name

  service_account = {
    ref = mws_iam_service_account.example.metadata.id
  }

  properties = {
    location = var.warehouse_location
  }
}

resource "mws_metastore_iceberg_role_binding" "example" {
  role_binding = var.role_binding_name
  role         = "metastore.icebergAdmin"
  subject      = mws_iam_service_account.example.metadata.id

  resource = {
    catalog = mws_metastore_iceberg.example.catalog_name
  }
}

variable "service_account_name" {
  type        = string
  default     = "example-iceberg-reader"
  description = "Service account name"
}

variable "catalog_name" {
  type        = string
  default     = "example-catalog"
  description = "Iceberg catalog name"
}

variable "warehouse_location" {
  type        = string
  default     = "s3://example-bucket/warehouse"
  description = "Object Storage location of the Iceberg warehouse"
}

variable "role_binding_name" {
  type        = string
  default     = "rb__example-iceberg-viewer"
  description = "Iceberg role binding name"
}
