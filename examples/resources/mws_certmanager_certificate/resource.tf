resource "mws_certmanager_certificate" "certificate" {
  name = var.certificate_name
  self_managed = {
    certificate = file(var.certificate_file_path)
    private_key = file(var.private_key_file_path)
  }
  self_managed_version = 1
}

variable "certificate_name" {
  type        = string
  default     = "my-test-certificate"
  description = "Certificate name"
}

variable "certificate_file_path" {
  type        = string
  description = "Certificate file path"
}

variable "private_key_file_path" {
  type        = string
  description = "Private key file path"
}
