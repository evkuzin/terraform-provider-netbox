package netbox

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxRir() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxRirCreate,
		Read:   resourceNetboxRirRead,
		Update: resourceNetboxRirUpdate,
		Delete: resourceNetboxRirDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
func resourceNetboxRirCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	data := models.RIR{}

	name := d.Get("name").(string)
	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}

	data.Name = &name
	data.Slug = &slug

	params := ipam.NewIpamRirsCreateParams().WithData(&data)
	res, err := api.Ipam.IpamRirsCreate(params, nil)
	if err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxRirUpdate(d, m)
}

func resourceNetboxRirRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamRirsReadParams().WithID(id)

	res, err := api.Ipam.IpamRirsRead(params, nil)
	if err != nil {

		return err
	}

	if res.GetPayload().Name != nil {
		d.Set("name", res.GetPayload().Name)
	}

	if res.GetPayload().Slug != nil {
		d.Set("slug", res.GetPayload().Slug)
	}

	return nil
}

func resourceNetboxRirUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.RIR{}

	name := d.Get("name").(string)
	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}

	data.Name = &name
	data.Slug = &slug

	params := ipam.NewIpamRirsUpdateParams().WithID(id).WithData(&data)
	_, err := api.Ipam.IpamRirsUpdate(params, nil)
	if err != nil {
		return err
	}
	return resourceNetboxRirRead(d, m)
}

func resourceNetboxRirDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamRirsDeleteParams().WithID(id)
	_, err := api.Ipam.IpamRirsDelete(params, nil)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
