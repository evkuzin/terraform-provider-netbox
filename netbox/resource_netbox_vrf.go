package netbox

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxVrf() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVrfCreate,
		Read:   resourceNetboxVrfRead,
		Update: resourceNetboxVrfUpdate,
		Delete: resourceNetboxVrfDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set:      schema.HashString,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceNetboxVrfCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	data := models.WritableVRF{}

	name := d.Get("name").(string)
	tenant_id := int64(d.Get("tenant_id").(int))

	data.Name = &name
	if tenant_id != 0 {
		data.Tenant = &tenant_id
	}

	data.Tags, _ = getNestedTagListFromResourceDataSet(api, d.Get("tags"))

	data.ExportTargets = []int64{}
	data.ImportTargets = []int64{}

	params := ipam.NewIpamVrfsCreateParams().WithData(&data)

	res, err := api.Ipam.IpamVrfsCreate(params, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxVrfRead(d, m)
}

func resourceNetboxVrfRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamVrfsReadParams().WithID(id)

	res, err := api.Ipam.IpamVrfsRead(params, nil)
	if err != nil {

		return err
	}

	d.Set("name", res.GetPayload().Name)
	if res.GetPayload().Tenant != nil {
		d.Set("tenant_id", res.GetPayload().Tenant.ID)
	} else {
		d.Set("tenant_id", nil)
	}
	return nil
}

func resourceNetboxVrfUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.WritableVRF{}

	name := d.Get("name").(string)

	tags, _ := getNestedTagListFromResourceDataSet(api, d.Get("tags"))

	data.Name = &name
	data.Tags = tags
	data.ExportTargets = []int64{}
	data.ImportTargets = []int64{}

	if tenantID, ok := d.GetOk("tenant_id"); ok {
		data.Tenant = int64ToPtr(int64(tenantID.(int)))
	}
	params := ipam.NewIpamVrfsPartialUpdateParams().WithID(id).WithData(&data)

	_, err := api.Ipam.IpamVrfsPartialUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxVrfRead(d, m)
}

func resourceNetboxVrfDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamVrfsDeleteParams().WithID(id)

	_, err := api.Ipam.IpamVrfsDelete(params, nil)
	if err != nil {
		return err
	}
	return nil
}
