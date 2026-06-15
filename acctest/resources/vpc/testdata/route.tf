resource "mws_vpc_route" "route" {
  route   = "%s"
  network = "%s"

  destination = {
    spec = {
      cidrs = ["10.0.0.0/8"]
    }
  }

  next_hop = {
    address = {
      ref = "%s"
    }
  }
}
