package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_JobHookDataSource_Basic(t *testing.T) {
	name := acctest.RandString(10)

	config := fmt.Sprintf(`
resource "neosync_connection" "source" {
	name = "%[1]s-src"

	postgres = {
		url = "test-url"
	}
}
resource "neosync_connection" "destination" {
	name = "%[1]s-dest"

	postgres = {
		url = "test-url2"
	}
}

resource "neosync_job" "job1" {
	name = "%[1]s"
	source = {
		postgres = {
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
				config = {
					passthrough = {}
				}
			}
		}
	]
}

resource "neosync_job_hook" "hook1" {
	name = "%[1]s"
	description = "Test hook"
	job_id = neosync_job.job1.id
	enabled = true
	priority = 5
	config = {
		sql = {
			query = "SELECT 1;"
			connection_id = neosync_connection.source.id
			timing = {
				pre_sync = {}
			}
		}
	}
}

data "neosync_job_hook" "test" {
	id = neosync_job_hook.hook1.id
}`, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Verify hook attributes are correctly read
					resource.TestCheckResourceAttrPair(
						"data.neosync_job_hook.test", "id",
						"neosync_job_hook.hook1", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.neosync_job_hook.test", "name",
						"neosync_job_hook.hook1", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.neosync_job_hook.test", "job_id",
						"neosync_job_hook.hook1", "job_id",
					),
				),
			},
		},
	})
}

func TestAcc_JobHookDataSource_NonExistent(t *testing.T) {
	config := `
data "neosync_job_hook" "test" {
	id = "non-existent-id"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Unable to get job hook"),
			},
		},
	})
}
