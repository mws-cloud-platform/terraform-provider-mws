resource "mws_vpc_firewall_rule" "firewall_rule" {
  metadata = {
    display_name = "Allow SSH"
  }

  network       = "%s"
  firewall_rule = "%s"

  priority  = 1000
  direction = "INGRESS"
  action    = "ALLOW"
  source = {
    spec = {
      cidrs = ["0.0.0.0/0"]
    }
  }
  destination = {
    spec = {
      cidrs = ["192.168.0.4/32"]
    }
  }
  proto_ports = ["TCP:22"]
  active      = true
}
