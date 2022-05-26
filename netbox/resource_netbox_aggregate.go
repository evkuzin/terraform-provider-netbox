package netbox

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxAggregate() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxAggregateCreate,
		Read:   resourceNetboxAggregateRead,
		Update: resourceNetboxAggregateUpdate,
		Delete: resourceNetboxAggregateDelete,

		Description: `From the [official documentation](https://docs.netbox.dev/en/stable/core-functionality/ipam/#aggregates):

> NetBox allows us to specify the portions of IP space that are interesting to us by defining aggregates. Typically, an aggregate will correspond to either an allocation of public (globally routable) IP space granted by a regional authority, or a private (internally-routable) designation.`,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rir_id": {
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
func resourceNetboxAggregateCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	data := models.WritableAggregate{}

	prefix := d.Get("prefix").(string)
	description := d.Get("description").(string)

	data.Prefix = &prefix
	data.Description = description

	if tenantID, ok := d.GetOk("tenant_id"); ok {
		data.Tenant = int64ToPtr(int64(tenantID.(int)))
	}

	if rirID, ok := d.GetOk("rir_id"); ok {
		data.Rir = int64ToPtr(int64(rirID.(int)))
	}

	data.Tags, _ = getNestedTagListFromResourceDataSet(api, d.Get("tags"))

	params := ipam.NewIpamAggregatesCreateParams().WithData(&data)
	res, err := api.Ipam.IpamAggregatesCreate(params, nil)
	if err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxAggregateRead(d, m)
}

func resourceNetboxAggregateRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamAggregatesReadParams().WithID(id)

	res, err := api.Ipam.IpamAggregatesRead(params, nil)
	if err != nil {

		return err
	}

	d.Set("description", res.GetPayload().Description)
	if res.GetPayload().Prefix != nil {
		d.Set("prefix", res.GetPayload().Prefix)
	}

	if res.GetPayload().Tenant != nil {
		d.Set("tenant_id", res.GetPayload().Tenant.ID)
	} else {
		d.Set("tenant_id", nil)
	}

	if res.GetPayload().Rir != nil {
		d.Set("rir_id", res.GetPayload().Rir.ID)
	} else {
		d.Set("rir_id", nil)
	}

	d.Set("tags", getTagListFromNestedTagList(res.GetPayload().Tags))

	return nil
}

func resourceNetboxAggregateUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.WritableAggregate{}
	prefix := d.Get("prefix").(string)
	description := d.Get("description").(string)

	data.Prefix = &prefix
	data.Description = description

	if tenantID, ok := d.GetOk("tenant_id"); ok {
		data.Tenant = int64ToPtr(int64(tenantID.(int)))
	}

	if rirID, ok := d.GetOk("rir_id"); ok {
		data.Rir = int64ToPtr(int64(rirID.(int)))
	}

	data.Tags, _ = getNestedTagListFromResourceDataSet(api, d.Get("tags"))

	params := ipam.NewIpamAggregatesUpdateParams().WithID(id).WithData(&data)
	_, err := api.Ipam.IpamAggregatesUpdate(params, nil)
	if err != nil {
		return err
	}
	return resourceNetboxAggregateRead(d, m)
}

func resourceNetboxAggregateDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamAggregatesDeleteParams().WithID(id)
	_, err := api.Ipam.IpamAggregatesDelete(params, nil)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
