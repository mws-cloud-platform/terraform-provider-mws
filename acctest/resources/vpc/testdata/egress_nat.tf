resource "mws_vpc_egress_nat" "{{.Name}}" {
  egress_nat = "{{.Name}}"

  external = {
    addresses = [
      {
        ref = "{{.ExternalAddress}}"
      }
    ]
  }

  internal = {
    subnets = ["{{.Subnet}}"]
  }

  network = "{{.Network}}"

  metadata = {
    description  = "This is an example Egress NAT resource"
    display_name = "Example Egress NAT"
  }
}
