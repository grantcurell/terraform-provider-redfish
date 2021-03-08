terraform {
  required_providers {
    redfish = {
      source = "dell/redfish"
    }
  }
}

// For servers without a uniquely defined username/password these values will be used
provider "redfish" {
  user = "root"
  password = "password"
}

resource "redfish_power" "system_power" {
  for_each = var.rack1

  redfish_server {
    user = each.value.user
    password = each.value.password
    endpoint = each.value.endpoint
    ssl_insecure = each.value.ssl_insecure
  }

  // The valid options are defined below.
  // Taken from the Redfish specification at: https://redfish.dmtf.org/schemas/DSP2046_2019.4.html
  /*
  | string           | Description                                                                             |
  |------------------|-----------------------------------------------------------------------------------------|
  | ForceOff         | Turn off the unit immediately (non-graceful shutdown).                                  |
  | ForceOn          | Turn on the unit immediately.                                                           |
  | ForceRestart     | Shut down immediately and non-gracefully and restart the system.                        |
  | GracefulRestart  | Shut down gracefully and restart the system.                                            |
  | GracefulShutdown | Shut down gracefully and power off.                                                     |
  | Nmi              | Generate a diagnostic interrupt, which is usually an NMI on x86 systems, to stop normal |
  |                  | operations, complete diagnostic actions, and, typically, halt the system.               |
  | On               | Turn on the unit.                                                                       |
  | PowerCycle       | Power cycle the unit.                                                                   |
  | PushPowerButton  | Simulate the pressing of the physical power button on this unit.                        |
  */
  desired_power_state = "On"
}

output "current_power_state" {
  value = redfish_power.system_power
}
