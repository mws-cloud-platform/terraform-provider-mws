resource "mws_vpc_address_group" "address_group" {
  address_group = "%s"
  network       = "%s"

  metadata = {
    display_name = "Example Address Group"
    description  = "This is an example Address Group resource"
  }

  addresses = [
    {
      ref = "%s"
    }
  ]
}
