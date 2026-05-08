provider "mws" {
  endpoint                            = "https://api.mwsapis.ru" # Endpoint for the provider API calls
  service_account_authorized_key_path = "path/to/authorized_key" # Path to the service account authorized key file
  project                             = "my-project"             # Project name
  zone                                = "ru-central1-a"          # Zone name
}
