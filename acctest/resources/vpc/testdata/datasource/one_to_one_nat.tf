data "mws_vpc_one_to_one_nat" "test" {
  one_to_one_nat = "%[1]s"
  network        = "%[2]s"
}
