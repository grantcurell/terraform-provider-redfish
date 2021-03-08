package redfish

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stmcginnis/gofish/redfish"
	"log"
)

func resourceRedFishPower() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedfishPowerUpdate,
		ReadContext:   resourceRedfishPowerRead,
		UpdateContext: resourceRedfishPowerUpdate,
		DeleteContext: resourceRedfishPowerDelete,
		Schema:        getResourceRedfishPowerSchema(),
	}
}

func getResourceRedfishPowerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"redfish_server": {
			Type: schema.TypeList,
			Required: true,
			Description: "List of server BMCs and their respective user credentials",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user": {
						Type: schema.TypeString,
						Optional: true,
						Description: "User name for login",
					},
					"password": {
						Type: schema.TypeString,
						Optional: true,
						Description: "User password for login",
						Sensitive: true,
					},
					"endpoint": {
						Type: schema.TypeString,
						Required: true,
						Description: "Server BMC IP address or hostname",
					},
					"ssl_insecure": {
						Type: schema.TypeBool,
						Optional: true,
						Description: "This field indicates whether the SSL/TLS certificate must be verified or not",
					},
				},
			},
		},
		"desired_power_state": {
			Type: schema.TypeString,
			Required: true,
			Description: "Desired power setting. Applicable values 'On','ForceOn','ForceOff','ForceRestart'," +
				"'GracefulRestart','GracefulShutdown','PushPowerButton','PowerCycle','Nmi'.",
		},
	}
}

func resourceRedfishPowerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] %s: Beginning read", d.Id())
	var diags diag.Diagnostics

	service, err := NewConfig(m.(*schema.ResourceData), d)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	// TODO getSystemResource should probably be in some sort of common file. See https://github.com/dell/terraform-provider-redfish/issues/21
	system, err := getSystemResource(service)
	if err != nil {
		log.Printf("[ERROR]: Failed to identify system: %s", err)
		return diag.Errorf(err.Error())
	}

	if err := d.Set("power_state", system.PowerState); err != nil {
		return diag.Errorf("[ERROR]: Could not retrieve system power state. %s", err)
	}

	return diags
}

// TODO - maybe we should centralize power management?
func resourceRedfishPowerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	resetType, ok := d.GetOk("desired_power_state")

	if !ok || resetType == nil {
		log.Printf("[ERROR]: TODO")
	}

	// Takes the m interface and feeds it the user input data d. You can then reference it with X.GetOk("user")
	service, err := NewConfig(m.(*schema.ResourceData), d)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	// TODO getSystemResource should probably be in some sort of common file. See https://github.com/dell/terraform-provider-redfish/issues/21
	system, err := getSystemResource(service)
	if err != nil {
		log.Printf("[ERROR]: Failed to identify system: %s", err)
		return diag.Errorf(err.Error())
	}

	log.Printf("[TRACE]: Performing system.Reset(%s)", resetType)
	if err = system.Reset(redfish.ResetType((resetType).(string))); err != nil {
		log.Printf("[WARN]: system.Reset returned an error: %s", err)
		return diag.Errorf(err.Error())
	}

	log.Printf("[TRACE]: system.Reset successful")
	return diags
}

func resourceRedfishPowerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
