package netbox

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxManufacturer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxManufacturerCreate,
		Read:   resourceNetboxManufacturerRead,
		Update: resourceNetboxManufacturerUpdate,
		Delete: resourceNetboxManufacturerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 30),
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceNetboxManufacturerCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	data := models.Manufacturer{}

	name := d.Get("name").(string)
	data.Name = &name

	slugValue, slugOk := d.GetOk("slug")
	// Default slug to name if not given
	if !slugOk {
		data.Slug = strToPtr(name)
	} else {
		data.Slug = strToPtr(slugValue.(string))
	}

	params := dcim.NewDcimManufacturersCreateParams().WithData(&data)

	res, err := api.Dcim.DcimManufacturersCreate(params, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxManufacturerRead(d, m)
}

func resourceNetboxManufacturerRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimManufacturersReadParams().WithID(id)

	res, err := api.Dcim.DcimManufacturersRead(params, nil)

	if err != nil {

		return err
	}

	d.Set("name", res.GetPayload().Name)
	d.Set("slug", res.GetPayload().Slug)

	return nil
}

func resourceNetboxManufacturerUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.Manufacturer{}

	name := d.Get("name").(string)
	data.Name = &name

	slugValue, slugOk := d.GetOk("slug")
	// Default slug to name if not given
	if !slugOk {
		data.Slug = strToPtr(name)
	} else {
		data.Slug = strToPtr(slugValue.(string))
	}

	params := dcim.NewDcimManufacturersPartialUpdateParams().WithID(id).WithData(&data)

	_, err := api.Dcim.DcimManufacturersPartialUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxManufacturerRead(d, m)
}

func resourceNetboxManufacturerDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimManufacturersDeleteParams().WithID(id)

	_, err := api.Dcim.DcimManufacturersDelete(params, nil)
	if err != nil {
		return err
	}
	return nil
}
