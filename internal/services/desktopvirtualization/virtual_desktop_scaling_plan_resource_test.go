package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualDesktopScalingPlanResource struct {
}

func TestAccVirtualDesktopScalingPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopScalingPlan_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccVirtualDesktopScalingPlan_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopScalingPlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_scaling_plan", "test")
	r := VirtualDesktopScalingPlanResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_scaling_plan"),
		},
	})
}

func (VirtualDesktopScalingPlanResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ScalingPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.ScalingPlansClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Virtual Desktop Scaling Plan %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.ScalingPlanProperties != nil), nil
}

func (VirtualDesktopScalingPlanResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"
}

resource "azurerm_virtual_desktop_scaling_plan" "test" {
	name				= "scalingPlan%x"
	location            = "westeurope"
	resource_group_name = azurerm_resource_group.test.name
	friendly_name		= "Scaling Plan Test"
	description			= "Test Scaling Plan"
	schedule 			= {
		name	= "Weekdays"
		days_of_week = ["Monday","Tuesday","Wednesday","Thursday","Friday"]
		ramp_up_start_time = {
			hour = 6
			minute = 0
		}
		ramp_up_load_balancing = "BreadthFirst"
		ramp_up_minimum_hosts_percent = 20
		ramp_up_capacity_threshold = 10
		peak_start_time = {
			hour = 9
			minute = 0
		}
		peak_load_balancing = "BreadthFirst"
		ramp_down_start_time = {
			hour = 18
			minute = 0
		}
		ramp_down_load_balancing = "DepthFirst"
		ramp_down_minimum_hosts_percent = 10
		ramp_down_capacity_threshold_percent = 5
		ramp_down_stop_hosts_when = "ZeroSessions"
		off_peak_start_time = {
			hour = 22
			minute = 0
		}
		off_peak_load_balancing = "DepthFirst"
	}
	hostpool_association = {
		hostpool_id = azurerm_virtual_desktop_host_pool.test.id
		scaling_plan_enabled = true
	}

}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, data.RandomString)
}

func (VirtualDesktopScalingPlanResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
	name     = "acctestRG-vdesktop-%d"
	location = "%s"
  }
  
  resource "azurerm_virtual_desktop_host_pool" "test" {
	name                 = "acctestHP%s"
	location             = azurerm_resource_group.test.location
	resource_group_name  = azurerm_resource_group.test.name
	type                 = "Pooled"
	validate_environment = true
	load_balancer_type   = "BreadthFirst"
  }
  
  resource "azurerm_virtual_desktop_scaling_plan" "test" {
	  name				= "scalingPlan%x"
	  location            = "westeurope"
	  resource_group_name = azurerm_resource_group.test.name
	  friendly_name		= "Scaling Plan Test"
	  description			= "Test Scaling Plan"
	  schedule 			= {
		  name	= "Weekdays"
		  days_of_week = ["Monday","Tuesday","Wednesday","Thursday","Friday"]
		  ramp_up_start_time = {
			  hour = 6
			  minute = 0
		  }
		  ramp_up_load_balancing = "BreadthFirst"
		  ramp_up_minimum_hosts_percent = 20
		  ramp_up_capacity_threshold = 10
		  peak_start_time = {
			  hour = 9
			  minute = 0
		  }
		  peak_load_balancing = "BreadthFirst"
		  ramp_down_start_time = {
			  hour = 18
			  minute = 0
		  }
		  ramp_down_load_balancing = "DepthFirst"
		  ramp_down_minimum_hosts_percent = 10
		  ramp_down_capacity_threshold_percent = 5
		  ramp_down_force_logoff_users = true
		  ramp_down_wait_time = 45
		  ramp_down_notification = "Please save your work and logoff in the next 45 minutes..."
		  ramp_down_stop_hosts_when = "ZeroSessions"
		  off_peak_start_time = {
			  hour = 22
			  minute = 0
		  }
		  off_peak_load_balancing = "DepthFirst"
	  }
	  hostpool_association = {
		  hostpool_id = azurerm_virtual_desktop_host_pool.test.id
		  scaling_plan_enabled = true
	  }
  
  }

`, data.RandomInteger, data.Locations.Secondary, data.RandomString, data.RandomString)
}

func (r VirtualDesktopScalingPlanResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_scaling_plan" "import" {
	name				= azurerm_virtual_desktop_scaling_plan.test.name
	location            = azurerm_virtual_desktop_scaling_plan.test.location
	resource_group_name = azurerm_virtual_desktop_scaling_plan.test.resource_group_name
	friendly_name		= azurerm_virtual_desktop_scaling_plan.test.friendly_name
	description			= azurerm_virtual_desktop_scaling_plan.test.description
}
`, r.basic(data))
}
