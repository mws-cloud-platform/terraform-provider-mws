resource "mws_compute_image" "image" {
  image  = "%s"
  family = "mws-ubuntu-2204-lts"
  source = {
    external_url = "https://storage.mwsapis.ru/mws-ubuntu/mws-ubuntu-2204-lts-v20250529.334787.raw"
  }
}
