resource "mws_secretmanager_secret" "example" {
  name   = var.secret_name
  active = true
}

variable "secret_name" {
  type        = string
  default     = "my-test-secret"
  description = "Name of the secret"
}
