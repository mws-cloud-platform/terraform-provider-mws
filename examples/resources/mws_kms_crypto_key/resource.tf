resource "mws_kms_crypto_key" "example" {
  key               = var.crypto_key_name
  default_algorithm = "AES_256_GCM"

  rotation_policy = {
    enabled                = true
    rotation_interval_days = 90
  }

  usage_policy = {
    enabled = true
  }

  destruction_policy = {
    default_destruction_interval_days = 30
  }

  metadata = {
    description  = "Example crypto key for encryption"
    display_name = "Example Crypto Key"
  }
}

variable "crypto_key_name" {
  type        = string
  default     = "example-crypto-key"
  description = "Name of the crypto key"
}
