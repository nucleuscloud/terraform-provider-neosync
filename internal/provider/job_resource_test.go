package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Job_Pg_Pg(t *testing.T) {
	name := acctest.RandString(10)

	config := fmt.Sprintf(`
resource "neosync_connection" "source" {
	name = "%s-src"

	postgres = {
		url = "test-url"
	}
}
resource "neosync_connection" "destination" {
	name = "%s-dest"

	postgres = {
		url = "test-url2"
	}
}

resource "neosync_job" "job1" {
	name = "%s"
	source = {
		postgres = {
			halt_on_new_column_addition = false
			connection_id = neosync_connection.source.id
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
	`, name, name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_job.job1", "id"),
				),
			},
		},
	})
}
