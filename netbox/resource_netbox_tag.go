package netbox

import (
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/extras"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxTagCreate,
		Read:   resourceNetboxTagRead,
		Update: resourceNetboxTagUpdate,
		Delete: resourceNetboxTagDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"color_hex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "9e9e9e",
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9a-f]{6}$"), "Must be hex color string"),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
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

func resourceNetboxTagCreate(d *schema.ResourceData, m interface{}) error {
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

	color := d.Get("color_hex").(string)
	description := d.Get("description").(string)
	params := extras.NewExtrasTagsCreateParams().WithData(
		&models.Tag{
			Name:        &name,
			Slug:        &slug,
			Color:       color,
			Description: description,
		},
	)

	res, err := api.Extras.ExtrasTagsCreate(params, nil)
	if err != nil {
		//return errors.New(getTextFromError(err))
		return err
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxTagRead(d, m)
}

func resourceNetboxTagRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := extras.NewExtrasTagsReadParams().WithID(id)

	res, err := api.Extras.ExtrasTagsRead(params, nil)
	if err != nil {

		return err
	}

	d.Set("name", res.GetPayload().Name)
	d.Set("slug", res.GetPayload().Slug)
	d.Set("color_hex", res.GetPayload().Color)
	d.Set("description", res.GetPayload().Description)
	return nil
}

func resourceNetboxTagUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.Tag{}

	name := d.Get("name").(string)
	color := d.Get("color_hex").(string)
	description := d.Get("description").(string)

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
	data.Color = color
	data.Description = description

	params := extras.NewExtrasTagsUpdateParams().WithID(id).WithData(&data)

	_, err := api.Extras.ExtrasTagsUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxTagRead(d, m)
}

func resourceNetboxTagDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := extras.NewExtrasTagsDeleteParams().WithID(id)

	_, err := api.Extras.ExtrasTagsDelete(params, nil)
	if err != nil {
		return err
	}
	return nil
}
