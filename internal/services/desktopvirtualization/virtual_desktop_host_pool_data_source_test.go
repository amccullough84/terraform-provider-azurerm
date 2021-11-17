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

type virtualDesktopHostPoolDataSource struct {
}

func TestAccDataSourceVirtualDesktopHostPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_desktop_host_pool", "test")
	r := virtualDesktopHostPoolDataSource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
				check.That(data.ResourceName).Key("friendly_name").HasValue("A Friendly Name!"),
				check.That(data.ResourceName).Key("description").HasValue("A Description!"),
				check.That(data.ResourceName).Key("custom_rdp_properties").HasValue("audiocapturemode:i:1;audiomode:i:0;"),
			),
		},
	})
}

func (virtualDesktopHostPoolDataSource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.HostPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.HostPoolsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Virtual Desktop Host Pool %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.HostPoolProperties != nil), nil
}

func (virtualDesktopHostPoolDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                     = "acctestHP%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  type                     = "Pooled"
  friendly_name            = "A Friendly Name!"
  description              = "A Description!"
  validate_environment     = true
  start_vm_on_connect      = true
  load_balancer_type       = "BreadthFirst"
  maximum_sessions_allowed = 100
  preferred_app_group_type = "Desktop"
  custom_rdp_properties    = "audiocapturemode:i:1;audiomode:i:0;"

  # Do not use timestamp() outside of testing due to https://github.com/hashicorp/terraform/issues/22461
  registration_info {
    expiration_date = timeadd(timestamp(), "48h")
  }
  lifecycle {
    ignore_changes = [
      registration_info[0].expiration_date,
    ]
  }

  tags = {
    Purpose = "Acceptance-Testing"
  }
}

data "azurerm_virtual_desktop_host_pool" "test" {
  name					   = azurerm_virtual_desktop_host_pool.test.name
  resource_group_name      = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomString)
}
