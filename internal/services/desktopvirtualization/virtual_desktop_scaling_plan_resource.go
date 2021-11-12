package desktopvirtualization

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2021-09-03-preview/desktopvirtualization"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualDesktopScalingPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopScalingPlanCreate,
		Read:   resourceVirtualDesktopScalingPlanRead,
		Update: resourceVirtualDesktopScalingPlanUpdate,
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
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
							ForceNew:     true,
						},

						"days_of_week": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
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
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},
									"Minute": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 59),
									},
								},
							},
						},

						"ramp_up_load_balancing": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},

						"ramp_up_minimum_hosts_percent": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},

						"ramp_up_capacity_threshold_percent": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},

						"peak_start_time": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},
									"Minute": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 59),
									},
								},
							},
						},

						"peak_load_balancing": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},

						"ramp_down_start_time": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},
									"Minute": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 59),
									},
								},
							},
						},

						"ramp_down_load_balancing": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},

						"ramp_down_minimum_hosts_percent": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},

						"ramp_down_capacity_threshold_percent": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},

						"ramp_down_force_logoff_users": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"ramp_down_stop_hosts_when": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(desktopvirtualization.StopHostsWhenZeroSessions),
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.StopHostsWhenZeroSessions),
								string(desktopvirtualization.StopHostsWhenZeroActiveSessions),
							}, false),
						},

						"ramp_down_wait_time": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  int(30),
						},

						"ramp_down_notifcation_message": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							RequiredWith: []string{"ramp_down_force_logoff_users"},
							ValidateFunc: validation.StringLenBetween(1, 512),
						},

						"off_peak_start_time": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"Hour": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},
									"Minute": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 59),
									},
								},
							},
						},

						"off_peak_load_balancing": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.LoadBalancerTypeBreadthFirst),
								string(desktopvirtualization.LoadBalancerTypeDepthFirst),
							}, false),
						},
					},
				},
			},

			"hostpool_association": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"hostpool_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.HostPoolID,
						},

						"scaling_plan_enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualDesktopScalingPlanCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan create")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resourceId := parse.NewScalingPlanID(subscriptionId, resourceGroup, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing virtual desktop scaling plan %q (resource group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ScalingPlanProperties != nil {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_scaling_plan", resourceId)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	context := desktopvirtualization.ScalingPlan{
		Location: &location,
		Tags:     tags.Expand(t),
		Name:     utils.String(d.Get("name").(string)),

		ScalingPlanProperties: &desktopvirtualization.ScalingPlanProperties{
			HostPoolType:       desktopvirtualization.ScalingHostPoolTypePooled, // Current implementation only supports pooled Host Pool types
			Description:        utils.String(d.Get("description").(string)),
			FriendlyName:       utils.String(d.Get("friendly_name").(string)),
			TimeZone:           utils.String(d.Get("time_zone").(string)),
			ExclusionTag:       utils.String(d.Get("exclusion_tag_name").(string)),
			Schedules:          expandScalingPlanSchedule(d),
			HostPoolReferences: expandHostPoolAssociation(d),
		},
	}

	if _, err := client.Create(ctx, resourceGroup, name, context); err != nil {
		return fmt.Errorf("Creating Virtual Desktop Host Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}

func resourceVirtualDesktopScalingPlanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan update")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resourceId := parse.NewScalingPlanID(subscriptionId, resourceGroup, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing virtual desktop scaling plan %q (resource group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ScalingPlanProperties != nil {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_scaling_plan", resourceId)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	context := desktopvirtualization.ScalingPlan{
		Location: &location,
		Tags:     tags.Expand(t),
		Name:     utils.String(d.Get("name").(string)),

		ScalingPlanProperties: &desktopvirtualization.ScalingPlanProperties{
			HostPoolType:       desktopvirtualization.ScalingHostPoolTypePooled, // Current implementation only supports pooled Host Pool types
			Description:        utils.String(d.Get("description").(string)),
			FriendlyName:       utils.String(d.Get("friendly_name").(string)),
			TimeZone:           utils.String(d.Get("time_zone").(string)),
			ExclusionTag:       utils.String(d.Get("exclusion_tag_name").(string)),
			Schedules:          expandScalingPlanSchedule(d),
			HostPoolReferences: expandHostPoolAssociation(d),
		},
	}

	if _, err := client.Create(ctx, resourceGroup, name, context); err != nil {
		return fmt.Errorf("Creating Virtual Desktop Host Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}

func resourceVirtualDesktopScalingPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Desktop Scaling Plan %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Making Read request on Virtual Desktop Scaling Plan %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ScalingPlanProperties; props != nil {

		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("exclusion_tag_name", string(*props.ExclusionTag))
		d.Set("time_zone", string(*props.TimeZone))
		if err := d.Set("schedules", flattenVirtualDesktopScalingPlanSchedule(props.Schedules)); err != nil {
			return fmt.Errorf("setting `schedules`: %+v", err)
		}

		if err := d.Set("hostpool_association", flattenVirtualDesktopScalingPlanHostpoolAssociations(props.HostPoolReferences)); err != nil {
			return fmt.Errorf("setting `hostpool_association`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualDesktopScalingPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Virtual Desktop Scaling Plan %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandScalingPlanSchedule(d *pluginsdk.ResourceData) *[]desktopvirtualization.ScalingSchedule {
	configs := d.Get("schedules").([]interface{})
	schedules := make([]desktopvirtualization.ScalingSchedule, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		rampUpTimeRaw := data["ramp_up_start_time"].([]interface{})
		rampUpTime := expandScalingPlanTime(rampUpTimeRaw)
		peakStartRaw := data["peak_start_time"].([]interface{})
		peakStartTime := expandScalingPlanTime(peakStartRaw)
		rampDownTimeRaw := data["ramp_down_start_time"].([]interface{})
		rampDownTime := expandScalingPlanTime(rampDownTimeRaw)
		offPeakStartRaw := data["off_peak_start_time"].([]interface{})
		offPeakStartTime := expandScalingPlanTime(offPeakStartRaw)

		schedule := desktopvirtualization.ScalingSchedule{
			Name:                           utils.String(data["name"].(string)),
			DaysOfWeek:                     utils.ExpandStringSlice(data["days_of_week"].([]interface{})),
			RampUpStartTime:                rampUpTime,
			RampUpLoadBalancingAlgorithm:   desktopvirtualization.SessionHostLoadBalancingAlgorithm(data["ramp_up_load_balancing"].(string)),
			RampUpMinimumHostsPct:          utils.Int32(int32(data["ramp_up_minimum_hosts_percent"].(int))),
			RampUpCapacityThresholdPct:     utils.Int32(int32(data["ramp_up_capacity_threshold_percent"].(int))),
			PeakStartTime:                  peakStartTime,
			PeakLoadBalancingAlgorithm:     desktopvirtualization.SessionHostLoadBalancingAlgorithm(data["peak_load_balancing"].(string)),
			RampDownStartTime:              rampDownTime,
			RampDownLoadBalancingAlgorithm: desktopvirtualization.SessionHostLoadBalancingAlgorithm(data["ramp_down_load_balancing"].(string)),
			RampDownMinimumHostsPct:        utils.Int32(int32(data["ramp_down_minimum_hosts_percent"].(int))),
			RampDownCapacityThresholdPct:   utils.Int32(int32(data["ramp_down_capacity_threshold_percent"].(int))),
			OffPeakStartTime:               offPeakStartTime,
			OffPeakLoadBalancingAlgorithm:  desktopvirtualization.SessionHostLoadBalancingAlgorithm(data["off_peak_load_balancing"].(string)),
		}
		logoffUsers := false
		if v := data["ramp_down_force_logoff_users"].(bool); v {
			logoffUsers = true
			schedule.RampDownWaitTimeMinutes = utils.Int32(int32(data["ramp_down_wait_time"].(int)))
			schedule.RampDownNotificationMessage = utils.String(data["ramp_down_notifcation_message"].(string))
		}
		schedule.RampDownForceLogoffUsers = utils.Bool(logoffUsers)

		// if v := data["next_hop_in_ip_address"].(string); v != "" {
		// 	schedule.RoutePropertiesFormat.NextHopIPAddress = &v
		// },

		schedules = append(schedules, schedule)
	}

	return &schedules
}

func expandHostPoolAssociation(d *pluginsdk.ResourceData) *[]desktopvirtualization.ScalingHostPoolReference {
	configs := d.Get("hostpool_association").([]interface{})
	hostpools := make([]desktopvirtualization.ScalingHostPoolReference, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		hostpool := desktopvirtualization.ScalingHostPoolReference{
			HostPoolArmPath:    utils.String(data["hostpool_id"].(string)),
			ScalingPlanEnabled: utils.Bool(data["scaling_plan_enabled"].(bool)),
		}

		hostpools = append(hostpools, hostpool)
	}

	return &hostpools
}

func expandScalingPlanTime(input []interface{}) *desktopvirtualization.Time {
	scheduletime := desktopvirtualization.Time{}
	if len(input) > 0 {
		raw := input[0].(map[string]interface{})
		scheduletime.Hour = utils.Int32(int32(raw["hour"].(int)))
		scheduletime.Minute = utils.Int32(int32(raw["minute"].(int)))
	}
	return &scheduletime
}

func flattenScalingPlanTime(input *desktopvirtualization.Time) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"hour":   *input.Hour,
			"minute": *input.Minute,
		},
	}
}

func flattenVirtualDesktopScalingPlanSchedule(input *[]desktopvirtualization.ScalingSchedule) []interface{} {
	results := make([]interface{}, 0)

	if schedules := input; schedules != nil {
		for _, schedule := range *schedules {
			r := make(map[string]interface{})

			r["name"] = *schedule.Name
			r["days_of_week"] = *schedule.DaysOfWeek
			r["ramp_up_start_time"] = flattenScalingPlanTime(schedule.RampUpStartTime)
			r["ramp_up_load_balancing"] = schedule.RampUpLoadBalancingAlgorithm
			r["ramp_up_minimum_hosts_percent"] = *schedule.RampUpMinimumHostsPct
			r["ramp_up_capacity_threshold_percent"] = *schedule.RampUpCapacityThresholdPct
			r["peak_start_time"] = flattenScalingPlanTime(schedule.PeakStartTime)
			r["peak_load_balancing"] = schedule.PeakLoadBalancingAlgorithm
			r["ramp_down_start_time"] = flattenScalingPlanTime(schedule.RampDownStartTime)
			r["ramp_down_load_balancing"] = schedule.RampDownLoadBalancingAlgorithm
			r["ramp_down_minimum_hosts_percent"] = schedule.RampDownMinimumHostsPct
			r["ramp_down_capacity_threshold_percent"] = schedule.RampDownCapacityThresholdPct
			r["ramp_down_force_logoff_users"] = schedule.RampDownForceLogoffUsers
			r["ramp_down_stop_hosts_when"] = schedule.RampDownStopHostsWhen
			r["ramp_down_wait_time"] = schedule.RampDownWaitTimeMinutes
			r["ramp_down_notification_message"] = *schedule.RampDownNotificationMessage
			r["off_peak_start_time"] = flattenScalingPlanTime(schedule.OffPeakStartTime)
			r["off_peak_load_balancing"] = schedule.OffPeakLoadBalancingAlgorithm

			results = append(results, r)
		}
	}

	return results
}
func flattenVirtualDesktopScalingPlanHostpoolAssociations(input *[]desktopvirtualization.ScalingHostPoolReference) []interface{} {
	results := make([]interface{}, 0)

	if hostpools := input; hostpools != nil {
		for _, hostpool := range *hostpools {
			r := make(map[string]interface{})

			r["hostpool_id"] = *hostpool.HostPoolArmPath
			r["scaling_plan_enabled"] = *hostpool.ScalingPlanEnabled

			results = append(results, r)
		}
	}

	return results

}
