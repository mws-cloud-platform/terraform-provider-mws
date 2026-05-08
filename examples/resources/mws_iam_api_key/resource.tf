resource "mws_iam_service_account" "example" {
  service_account = var.service_account_name
}

resource "mws_iam_api_key" "example" {
  api_key         = var.api_key_name
  service_account = mws_iam_service_account.example.service_account

  metadata = {
    display_name = "Example Api-Key"
  }
}

variable "service_account_name" {
  type        = string
  default     = "my-test-service-account"
  description = "Service account name"
}

variable "api_key_name" {
  type        = string
  default     = "my-test-api-key"
  description = "Api-Key name"
}
