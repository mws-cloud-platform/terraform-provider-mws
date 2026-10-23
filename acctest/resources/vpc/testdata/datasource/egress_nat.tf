data "mws_vpc_egress_nat" "test" {
  egress_nat = "%[1]s"
  network    = "%[2]s"
}
