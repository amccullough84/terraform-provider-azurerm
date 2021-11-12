package desktopvirtualization

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2021-09-03-preview/desktopvirtualization"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualDesktopScalingPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopScalingPlanCreateUpdate,
		Read:   resourceVirtualDesktopScalingPlanRead,
		Update: resourceVirtualDesktopScalingPlanCreateUpdate,
		Delete: resourceVirtualDesktopScalingPlanDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ScalingPlanID(id)
			return err
		}),

		SchemaVersion: 1,
		
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			// To-Do - Add Identity Block

			"friendly_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},

			"exclusion_tag_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),	
			},

			

			"time_zone": {
				Type:        pluginsdk.TypeString,
				Optional:    true,
				Computed:	 true,
				ValidateFunc: computeValidate.VirtualMachineTimeZone(),
			},

			"schedules": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
							Required:     true,
							ForceNew:	  true,
						},

						"days_of_week": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Monday",
									"Tuesday",
									"Wednesday",
									"Thursday",
									"Friday",
									"Saturday",
									"Sunday",
								}, false),
							},
						},

						"ramp_up_start_time": {
							Type: pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,23),

									},
									"Minute": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,59),

									},
								},
							},
						
						},

						"ramp_up_load_balancing": {
							Type: pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},

						"ramp_up_minimum_hosts_percent": {
							Type: pluginsdk.TypeInt,
							Optional: true,
							ValidateFunc: validation.IntBetween(0,100),
					 
						},

						"ramp_up_capacity_threshold_percent" {
							Type: pluginsdk.TypeInt,
							Required: true,
							ValidateFunc: validation.IntBetween(0,100),

						},

						"peak_start_time": {
							Type: pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,23),

									},
									"Minute": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,59),

									},
								},
							},
						
						},

						"peak_load_balancing": {
							Type: pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},

						"ramp_down_start_time": {
							Type: pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,23),

									},
									"Minute": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,59),

									},
								},
							},
						
						},

						"ramp_down_load_balancing": {
							Type: pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},


						"ramp_down_minimum_hosts_percent": {
							Type: pluginsdk.TypeInt,
							Optional: true,
							ValidateFunc: validation.IntBetween(0,100),
					 
						},

						"ramp_down_capacity_threshold_percent" {
							Type: pluginsdk.TypeInt,
							Required: true,
							ValidateFunc: validation.IntBetween(0,100),

						},

						"ramp_down_force_logoff_users": {
							Type: pluginsdk.TypeBool
							Optional: true
							Default: false

						},

						"ramp_down_stop_hosts_when": {
							Type: pluginsdk.TypeString,
							Optional: true,
							Default: string(desktopvirtualization.StopHostsWhenZeroSessions),
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.StopHostsWhenZeroSessions),
								string(desktopvirtualization.StopHostsWhenZeroActiveSessions),
							}, false),

						},

						"ramp_down_wait_time": {
							Type: pluginsdk.TypeInt,
							Optional: true,
							Default: int(30)
						},

						"ramp_down_notifcation_message": {
							Type: pluginsdk.TypeString,
							Optional: true,
							RequiredWith: []string{"ramp_down_force_logoff_users"}
							ValidateFunc: validation.StringLenBetween(1,512)

						},

						"off_peak_start_time": {
							Type: pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,23),

									},
									"Minute": {
										Type: pluginsdk.TypeInt,
										Required: true,
										ValidateFunc: validation.IntBetween(0,59),

									},
								},
							},
						
						},

						"off_peak_load_balancing": {
							Type: pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},

					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}