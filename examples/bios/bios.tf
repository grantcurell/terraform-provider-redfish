terraform {
  required_providers {
    redfish = {
      source = "dell/redfish"
    }
  }
}

provider "redfish" {
    //user = "admin"
    //password = "passw0rd"
}

data "redfish_bios" "bios" {
  for_each = var.rack1

  redfish_server {
    user = each.value.user
    password = each.value.password
    endpoint = each.value.endpoint
    ssl_insecure = each.value.ssl_insecure
  }

}

resource "redfish_bios" "bios" {
  for_each = var.rack1

  redfish_server {
    user = each.value.user
    password = each.value.password
    endpoint = each.value.endpoint
    ssl_insecure = each.value.ssl_insecure
  }

  attributes = data.redfish_bios.bios[each.key].attributes

  settings_apply_time = "OnReset"
  action_after_apply = "ForceRestart"
  //wait_to_finish = true
}

output "bios_attributes" {
  value = data.redfish_bios.bios
}

