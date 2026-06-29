resource "mws_gpt_deployment" "deployment" {
  deployment_name = var.deployment_name
  model           = "glm-4.6-357b"
  is_active       = true

  metadata = {
    display_name = "Example GPT Deployment"
    description  = "An example GPT model deployment"
  }
}

variable "deployment_name" {
  type        = string
  default     = "my-gpt-deployment"
  description = "GPT deployment name"
}
