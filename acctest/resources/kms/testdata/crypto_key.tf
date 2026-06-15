resource "mws_kms_crypto_key" "test_kms_key" {
  key               = "%s"
  default_algorithm = "AES_256_GCM"

  rotation_policy = {
    enabled               = true
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

