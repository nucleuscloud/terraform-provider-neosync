package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Connection_DataSource(t *testing.T) {
	connectionName := acctest.RandString(10)
	testAccExampleDataSourceConfig := fmt.Sprintf(`
resource "neosync_connection" "test" {
  name = "%s"

	postgres = {
		url = "test-url"
	}
}
data "neosync_connection" "test" {
  id = neosync_connection.test.id
}
`, connectionName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.neosync_connection.test", "id"),
					resource.TestCheckResourceAttr("data.neosync_connection.test", "name", connectionName),
				),
			},
		},
	})
}
