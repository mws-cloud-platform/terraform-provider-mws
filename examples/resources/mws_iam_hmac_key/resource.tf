resource "mws_iam_service_account" "example" {
  service_account = var.service_account_name
}

resource "mws_iam_hmac_key" "example" {
  key_name        = var.hmac_key_name
  service_account = mws_iam_service_account.example.service_account

  metadata = {
    display_name = "Example HMAC Key"
  }
}

variable "service_account_name" {
  type        = string
  default     = "my-test-service-account"
  description = "Service account name"
}

variable "hmac_key_name" {
  type        = string
  default     = "my-test-hmac-key"
  description = "HMAC Key name"
}
