resource "mws_vpc_network" "network" {
  network         = "%s"
  mtu             = "%d"
  internet_access = "%t"
}
