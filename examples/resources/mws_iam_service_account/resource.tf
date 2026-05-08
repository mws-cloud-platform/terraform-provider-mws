resource "mws_iam_service_account" "example" {
  service_account = var.service_account_name

  metadata = {
    display_name = "Example SA"
  }
}

variable "service_account_name" {
  type        = string
  default     = "my-test-service-account"
  description = "Service account name"
}
