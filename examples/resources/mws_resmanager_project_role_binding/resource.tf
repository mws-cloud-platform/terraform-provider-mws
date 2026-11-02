resource "mws_iam_service_account" "admin" {
  service_account = var.service_account_name
}

resource "mws_resmanager_project_role_binding" "example" {
  role_binding = var.role_binding_name
  role         = var.role

  subject = {
    service_account = mws_iam_service_account.admin.metadata.id
  }
}

variable "service_account_name" {
  type        = string
  default     = "example-project-admin"
  description = "Service account name"
}

variable "role_binding_name" {
  type        = string
  default     = "example-project-role-binding"
  description = "Project role binding name"
}

variable "role" {
  type        = string
  default     = "iam/roles/admin"
  description = "Service account project role"
}
