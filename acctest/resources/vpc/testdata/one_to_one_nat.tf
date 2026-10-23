resource "mws_vpc_one_to_one_nat" "%[1]s" {
  one_to_one_nat = "%[1]s"
  network        = "%[2]s"

  external = {
    address = {
      ref = "%[3]s"
    }
  }

  internal = {
    address = {
      ref = "%[4]s"
    }
  }
}
