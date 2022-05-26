package netbox

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
)

func dataSourceNetboxDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetboxDeviceRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"site": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetboxDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*client.NetBoxAPI)
	returnDiag := diag.Diagnostics{}
	name := d.Get("name").(string)
	params := dcim.NewDcimDevicesListParams()
	params.Name = &name
	limit := int64(2) // Limit of 2 is enough
	params.Limit = &limit

	res, err := api.Dcim.DcimDevicesList(params, nil)
	if err != nil {
		returnDiag = append(returnDiag, diag.Diagnostic{
			Severity:      0,
			Summary:       "",
			Detail:        err.Error(),
			AttributePath: nil,
		})
	}

	if *res.GetPayload().Count > int64(1) {
		returnDiag = append(returnDiag, diag.Diagnostic{
			Severity:      0,
			Summary:       "More than one result",
			Detail:        "Specify a more narrow filter",
			AttributePath: nil,
		})
	}
	if *res.GetPayload().Count == int64(0) {
		returnDiag = append(returnDiag, diag.Diagnostic{
			Severity:      0,
			Summary:       "No results",
			Detail:        "Specify a more narrow filter",
			AttributePath: nil,
		})
	}
	result := res.GetPayload().Results[0]
	d.SetId(strconv.FormatInt(result.ID, 10))
	d.Set("name", result.Name)
	d.Set("id", result.ID)
	d.Set("status", result.Status.Value)
	d.Set("device_type", result.DeviceType.Display)
	d.Set("site", result.Site.Name)
	if result.PrimaryIP != nil {
		d.Set("primary_ip", result.PrimaryIP.Address)
	}
	return returnDiag
}
