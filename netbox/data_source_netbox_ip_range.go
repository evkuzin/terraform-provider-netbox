package netbox

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
)

func dataSourceNetboxIpRange() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetboxIpRangeRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"contains": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},
		},
	}
}

func dataSourceNetboxIpRangeRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	contains := d.Get("contains").(string)

	params := ipam.NewIpamIPRangesListParams()
	params.Contains = &contains

	limit := int64(2) // Limit of 2 is enough
	params.Limit = &limit

	res, err := api.Ipam.IpamIPRangesList(params, nil)
	if err != nil {
		return err
	}

	if *res.GetPayload().Count > int64(1) {
		return errors.New("More than one result. Specify a more narrow filter")
	}
	if *res.GetPayload().Count == int64(0) {
		return errors.New("No result")
	}
	result := res.GetPayload().Results[0]
	d.Set("id", result.ID)
	d.SetId(strconv.FormatInt(result.ID, 10))
	return nil
}
