package netbox

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxPlatform() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxPlatformCreate,
		Read:   resourceNetboxPlatformRead,
		Update: resourceNetboxPlatformUpdate,
		Delete: resourceNetboxPlatformDelete,

		Description: `From the [official documentation](https://docs.netbox.dev/en/stable/core-functionality/devices/#platforms):

> A platform defines the type of software running on a device or virtual machine. This can be helpful to model when it is necessary to distinguish between different versions or feature sets. Note that two devices of the same type may be assigned different platforms: For example, one Juniper MX240 might run Junos 14 while another runs Junos 15.`,

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

func resourceNetboxPlatformCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	name := d.Get("name").(string)

	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}

	params := dcim.NewDcimPlatformsCreateParams().WithData(
		&models.WritablePlatform{
			Name: &name,
			Slug: &slug,
		},
	)

	res, err := api.Dcim.DcimPlatformsCreate(params, nil)
	if err != nil {
		//return errors.New(getTextFromError(err))
		return err
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxPlatformRead(d, m)
}

func resourceNetboxPlatformRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimPlatformsReadParams().WithID(id)

	res, err := api.Dcim.DcimPlatformsRead(params, nil)

	if err != nil {

		return err
	}

	d.Set("name", res.GetPayload().Name)
	d.Set("slug", res.GetPayload().Slug)
	return nil
}

func resourceNetboxPlatformUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.WritablePlatform{}

	name := d.Get("name").(string)

	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}

	data.Slug = &slug
	data.Name = &name

	params := dcim.NewDcimPlatformsPartialUpdateParams().WithID(id).WithData(&data)

	_, err := api.Dcim.DcimPlatformsPartialUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxPlatformRead(d, m)
}

func resourceNetboxPlatformDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := dcim.NewDcimPlatformsDeleteParams().WithID(id)

	_, err := api.Dcim.DcimPlatformsDelete(params, nil)
	if err != nil {
		return err
	}
	return nil
}
