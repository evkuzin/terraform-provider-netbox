package netbox

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/tenancy"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxTenantGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxTenantGroupCreate,
		Read:   resourceNetboxTenantGroupRead,
		Update: resourceNetboxTenantGroupUpdate,
		Delete: resourceNetboxTenantGroupDelete,

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
			"parent_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceNetboxTenantGroupCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	name := d.Get("name").(string)
	parent_id := int64(d.Get("parent_id").(int))
	description := d.Get("description").(string)

	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}

	data := &models.WritableTenantGroup{}
	data.Name = &name
	data.Slug = &slug
	data.Description = description
	if parent_id != 0 {
		data.Parent = &parent_id
	}

	params := tenancy.NewTenancyTenantGroupsCreateParams().WithData(data)

	res, err := api.Tenancy.TenancyTenantGroupsCreate(params, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxTenantGroupRead(d, m)
}

func resourceNetboxTenantGroupRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)

	params := tenancy.NewTenancyTenantGroupsReadParams().WithID(id)

	res, err := api.Tenancy.TenancyTenantGroupsRead(params, nil)
	if err != nil {

		return err
	}

	d.Set("name", res.GetPayload().Name)
	d.Set("slug", res.GetPayload().Slug)
	d.Set("description", res.GetPayload().Description)
	if res.GetPayload().Parent != nil {
		d.Set("parent", res.GetPayload().Parent.ID)
	}
	return nil
}

func resourceNetboxTenantGroupUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.WritableTenantGroup{}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	parent_id := int64(d.Get("parent_id").(int))

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
	data.Description = description
	if parent_id != 0 {
		data.Parent = &parent_id
	}
	params := tenancy.NewTenancyTenantGroupsPartialUpdateParams().WithID(id).WithData(&data)

	_, err := api.Tenancy.TenancyTenantGroupsPartialUpdate(params, nil)
	if err != nil {
		return err
	}

	return resourceNetboxTenantGroupRead(d, m)
}

func resourceNetboxTenantGroupDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := tenancy.NewTenancyTenantGroupsDeleteParams().WithID(id)

	_, err := api.Tenancy.TenancyTenantGroupsDelete(params, nil)
	if err != nil {
		return err
	}
	return nil
}
