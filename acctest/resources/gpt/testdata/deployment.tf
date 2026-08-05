resource "mws_gpt_deployment" "deployment" {
  deployment_name = "%s"
  is_active       = true
  model           = "glm-5.2"
}
