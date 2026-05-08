resource "mws_compute_image" "image" {
  image  = var.image_name
  family = var.image_family
  source = {
    external_url = var.image_source_url
  }
}

variable "image_name" {
  type        = string
  default     = "mws-ubuntu-2204-lts-v20250529"
  description = "Image name"
}

variable "image_family" {
  type        = string
  default     = "mws-ubuntu-2204-lts"
  description = "Image family"
}

variable "image_source_url" {
  type        = string
  default     = "https://storage.mwsapis.ru/mws-ubuntu/mws-ubuntu-2204-lts-v20250529.334787.raw"
  description = "Image external source URL"
}
