---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_scaling_plan"
description: |-
  Manages a Virtual Desktop Scaling Plan .
---

# azurerm_virtual_desktop_scaling_plan

Manages a Virtual Desktop Scaling Plan .

## Example Usage

```hcl
resource "azurerm_virtual_desktop_scaling_plan" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
  time_zone = "TODO"
  hostpool_type = "TODO"

  schedule {
    name = "example"
    ramp_down_notification_message = "TODO"
    off_peak_load_balancing_algorithm = "TODO"
    peak_start_time = "TODO"
    ramp_down_minimum_hosts_percent = 42
    ramp_down_capacity_threshold_percent = 42
    ramp_down_force_logoff_users = false
    days_of_week = [ "example" ]
    peak_load_balancing_algorithm = "TODO"
    ramp_down_load_balancing_algorithm = "TODO"
    ramp_down_wait_time_minutes = 42
    ramp_up_start_time = "TODO"
    ramp_up_load_balancing_algorithm = "TODO"
    ramp_down_start_time = "TODO"
    ramp_down_stop_hosts_when = "TODO"
    off_peak_start_time = "TODO"    
  }
}
```

## Arguments Reference

The following arguments are supported:

* `hostpool_type` - (Required) TODO. Changing this forces a new Virtual Desktop Scaling Plan  to be created.

* `location` - (Required) The Azure Region where the Virtual Desktop Scaling Plan  should exist. Changing this forces a new Virtual Desktop Scaling Plan  to be created.

* `name` - (Required) The name which should be used for this Virtual Desktop Scaling Plan . Changing this forces a new Virtual Desktop Scaling Plan  to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Desktop Scaling Plan  should exist. Changing this forces a new Virtual Desktop Scaling Plan  to be created.

* `schedule` - (Required) One or more `schedule` blocks as defined below.

* `time_zone` - (Required) TODO.

---

* `description` - (Optional) TODO.

* `exclusion_tag` - (Optional) TODO.

* `friendly_name` - (Optional) TODO.

* `hostpool_reference` - (Optional) One or more `hostpool_reference` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Virtual Desktop Scaling Plan .

---

A `hostpool_reference` block supports the following:

* `hostpool_id` - (Required) The ID of the TODO.

* `scaling_plan_enabled` - (Required) Should the TODO be enabled?

---

A `schedule` block supports the following:

* `days_of_week` - (Required) Specifies a list of TODO.

* `name` - (Required) The name which should be used for this TODO.

* `off_peak_load_balancing_algorithm` - (Required) TODO.

* `off_peak_start_time` - (Required) TODO.

* `peak_load_balancing_algorithm` - (Required) TODO.

* `peak_start_time` - (Required) TODO.

* `ramp_down_capacity_threshold_percent` - (Required) TODO.

* `ramp_down_force_logoff_users` - (Required) TODO.

* `ramp_down_load_balancing_algorithm` - (Required) TODO.

* `ramp_down_minimum_hosts_percent` - (Required) TODO.

* `ramp_down_notification_message` - (Required) TODO.

* `ramp_down_start_time` - (Required) TODO.

* `ramp_down_stop_hosts_when` - (Required) TODO.

* `ramp_down_wait_time_minutes` - (Required) TODO.

* `ramp_up_load_balancing_algorithm` - (Required) TODO.

* `ramp_up_start_time` - (Required) TODO.

* `ramp_up_capacity_threshold_percent` - (Optional) TODO.

* `ramp_up_minimum_hosts_percent` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Desktop Scaling Plan .

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Scaling Plan .
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Scaling Plan .
* `update` - (Defaults to 1 hour) Used when updating the Virtual Desktop Scaling Plan .
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Desktop Scaling Plan .

## Import

Virtual Desktop Scaling Plan s can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_scaling_plan.example C:/Program Files/Git/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/scalingPlans/plan1
```