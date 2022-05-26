package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxDeviceDataSource_basic(t *testing.T) {

	testSlug := "dvrl_ds_basic"
	testName := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_device" "test" {
  name = "%[1]s"
  device_type_id = "1"
  device_role = "1"
  site_id = "1"
}
data "netbox_device" "test" {
  depends_on = [netbox_device.test]
  name = "%[1]s"
}`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_device.test", "id", "netbox_device.test", "id"),
				),
			},
		},
	})
}
