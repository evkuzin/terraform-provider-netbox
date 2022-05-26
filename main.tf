resource "netbox_device" "test" {
  name = "shauser-evgeny1"
  device_type_id = "1"
  site_id = "1"
  device_role = "1"
}


data "netbox_device" "test" {
  name = netbox_device.test.name
}

terraform {
  required_providers {
    netbox = {
      source  = "evgeny/netbox"
      version = "0.0.1"
    }
  }
}

provider "netbox" {
  server_url = "shauser-evgeny1:8000"
  api_token = "0123456789abcdef0123456789abcdef01234567"
#  skip_version_check = true
}