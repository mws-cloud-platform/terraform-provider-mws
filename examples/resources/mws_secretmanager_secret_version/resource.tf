resource "mws_secretmanager_secret" "example" {
  name   = var.secret_name
  active = true
}

resource "mws_secretmanager_secret_version" "example" {
  name   = mws_secretmanager_secret.example.name
  active = true

  data = {
    foo = "bar"
  }
}

variable "secret_name" {
  type        = string
  default     = "my-secret"
  description = "Name of the secret"
}

variable "version_name" {
  type        = string
  default     = "my-secret-version"
  description = "Name of the secret version"
}
