resource "mws_vpc_network" "network" {
  network         = var.network_name
  mtu             = 1500
  internet_access = true
}

resource "mws_vpc_firewall_rule" "firewall_rule" {
  metadata = {
    display_name = "Allow SSH"
  }

  firewall_rule = var.firewall_rule_name
  network       = mws_vpc_network.network.network

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

variable "network_name" {
  type        = string
  default     = "my-test-network"
  description = "Network name"
}

variable "firewall_rule_name" {
  type        = string
  default     = "my-test-network-allow-ssh-firewall-rule"
  description = "Firewall rule name"
}
