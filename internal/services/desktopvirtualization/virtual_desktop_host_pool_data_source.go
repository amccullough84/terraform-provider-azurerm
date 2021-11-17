package desktopvirtualization

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceVirtualDesktopHostPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualDesktopHostPoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"load_balancer_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"friendly_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"validate_environment": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"custom_rdp_properties": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"personal_desktop_assignment_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"maximum_sessions_allowed": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"start_vm_on_connect": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"preferred_app_group_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"registration_info": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"expiration_date": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"reset_token": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"token": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVirtualDesktopHostPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewHostPoolID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Host Pool %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Making Read request on Virtual Desktop Host Pool %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.HostPoolProperties; props != nil {
		maxSessionLimit := 0
		if props.MaxSessionLimit != nil {
			maxSessionLimit = int(*props.MaxSessionLimit)
		}

		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("maximum_sessions_allowed", maxSessionLimit)
		d.Set("load_balancer_type", string(props.LoadBalancerType))
		d.Set("personal_desktop_assignment_type", string(props.PersonalDesktopAssignmentType))
		d.Set("preferred_app_group_type", string(props.PreferredAppGroupType))
		d.Set("type", string(props.HostPoolType))
		d.Set("validate_environment", props.ValidationEnvironment)
		d.Set("start_vm_on_connect", props.StartVMOnConnect)
		d.Set("custom_rdp_properties", props.CustomRdpProperty)

		if err := d.Set("registration_info", flattenVirtualDesktopHostPoolRegistrationInfo(props.RegistrationInfo)); err != nil {
			return fmt.Errorf("setting `registration_info`: %+v", err)
		}
	}

	return nil
}
