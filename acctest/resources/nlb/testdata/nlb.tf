resource "mws_nlb_nlb" "nlb" {
  network = "%s"
  nlb     = "%s"

  listener = {
    internal = {
      address = {
        ref = "%s"
      }
    }
  }

  rules = [
    {
      proto_port  = "TCP:80"
      target_port = 8080

      target_address_groups = [
        {
          ref = "%s"
        }
      ]

      health_check = {
        protocol = {
          tcp = {
            port = 8080
          }
        }

        interval            = "30s"
        timeout             = "3s"
        healthy_threshold   = 3
        unhealthy_threshold = 3
      }
    }
  ]
}
