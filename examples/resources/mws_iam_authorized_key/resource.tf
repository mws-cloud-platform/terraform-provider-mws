resource "mws_iam_service_account" "example" {
  service_account = var.service_account_name
}

resource "mws_iam_authorized_key" "example" {
  authorized_key  = var.authorized_key_name
  service_account = mws_iam_service_account.example.service_account
  key_algorithm   = "ES256"

  active          = true
  expiration_time = "2027-01-01T00:00:00Z"

  metadata = {
    display_name = "Example Authorized Key"
  }
}

variable "service_account_name" {
  type        = string
  default     = "my-test-service-account"
  description = "Service account name"
}

variable "authorized_key_name" {
  type        = string
  default     = "my-test-auth-key"
  description = "Authorized key name"
}
