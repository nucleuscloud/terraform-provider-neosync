package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Connection_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.neosync_connection.test", "id"),
					resource.TestCheckResourceAttr("data.neosync_connection.test", "name", "foo"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
resource "neosync_connection" "test" {
  name = "foo"

	postgres = {
		url = "test-url"
	}
}
data "neosync_connection" "test" {
  id = neosync_connection.test.id
}
`
