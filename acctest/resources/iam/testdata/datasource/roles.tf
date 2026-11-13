data "mws_iam_roles" "public" {
  type   = "public"
  filter = "metadata.id.name = \"%s\""
}
