package redfish

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stmcginnis/gofish"
)

func dataSourceRedfishBios() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedfishBiosRead,
		Schema:      getDataSourceRedfishBiosSchema(),
	}
}

func getDataSourceRedfishBiosSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"redfish_server": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "List of server BMCs and their respective user credentials",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "User name for login",
					},
					"password": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "User password for login",
						Sensitive:   true,
					},
					"endpoint": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Server BMC IP address or hostname",
					},
					"ssl_insecure": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "This field indicates whether the SSL/TLS certificate must be verified or not",
					},
				},
			},
		},
		"odata_id": {
			Type:        schema.TypeString,
			Description: "OData ID for the Bios resource",
			Computed:    true,
		},
		"attributes": {
			Type:        schema.TypeMap,
			Description: "Bios attributes",
			Elem: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			Computed: true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Id",
			Computed:    true,
		},
	}
}

// TODO why not roll this all into one function?
func dataSourceRedfishBiosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	service, err := NewConfig(m.(*schema.ResourceData), d)
	if err != nil {
		return diag.Errorf(err.Error())
	}
	return readRedfishBios(service, d)
}

// readRedfishBios Because there are several hundred BIOS settings that aren't yet part of the state we need a way to
// bring them into the configuration beforehand. This function reaches out to the servers and pulls their BIOS settings
// and puts them into a value in the config called attributes. These attributes are later combined with any user input.
// The user input will override any values pre-existing.
func readRedfishBios(service *gofish.Service, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	systems, err := service.Systems()
	if err != nil {
		return diag.Errorf("error fetching computer systems collection: %s", err)
	}

	bios, err := systems[0].Bios()
	if err != nil {
		return diag.Errorf("error fetching bios: %s", err)
	}

	// TODO: BIOS Attributes' values might be any of several types.
	// terraform-sdk currently does not support a map with different
	// value types. So we will convert int and float values to string
	attributes := make(map[string]interface{})

	// copy from the BIOS attributes to the new bios attributes map
	for key, value := range bios.Attributes {
		attributes[key] = value
	}

	if err := d.Set("odata_id", bios.ODataID); err != nil { // TODO - what is this value?
		return diag.Errorf("error setting bios OData ID: %s", err)
	}

	if err := d.Set("id", bios.ID); err != nil {
		return diag.Errorf("error setting bios ID: %s", err)
	}

	if err := d.Set("attributes", attributes); err != nil {
		return diag.Errorf("error setting bios attributes: %s", err)
	}

	// Set the ID to the redfish endpoint + bios @odata.id
	serverConfig := d.Get("redfish_server").([]interface{})
	endpoint := serverConfig[0].(map[string]interface{})["endpoint"].(string)
	biosResourceId := endpoint + bios.ODataID
	d.SetId(biosResourceId)

	return diags
}
