package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Job_DataSource(t *testing.T) {
	name := acctest.RandString(10)

	config := fmt.Sprintf(`
resource "neosync_connection" "destination" {
	name = "%s-src"

	postgres = {
		url = "test-url"
	}
}

resource "neosync_job" "job1" {
	name = "%s"
	source = {
		generate = {
			fk_source_connection_id = neosync_connection.destination.id
			schemas = [
				{
					schema = "public"
					tables = [
						{
							table = "users"
							row_count = 10
						}
					]
				},
			]
		}
	}
	destinations = [
		{
			connection_id = neosync_connection.destination.id
			postgres = {
				init_table_schema = false
				truncate_table = {
					truncate_before_insert = true
					cascade = true
				}
			}
		}
	]
	mappings = [
		{
			schema = "public"
			table = "users"
			column = "id"
			transformer = {
				source = "passthrough"
				config = {
					passthrough = {}
				}
			}
		}
	]
}

data "neosync_job" "job1" {
	id = neosync_job.job1.id
}
	`, name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.neosync_job.job1", "id"),
					resource.TestCheckResourceAttr("data.neosync_job.job1", "name", name),
				),
			},
		},
	})
}
