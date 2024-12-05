package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test creating a single hook.
func TestAcc_JobHook_Create(t *testing.T) {
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

resource "neosync_job_hook" "jh1" {
	name = "%s"
	description = "this is a description"
	job_id = neosync_job.job1.id
	enabled = true
	priority = 5
	config = {
		sql = {
			query = "select 1;"
			connection_id = neosync_connection.destination.id
			timing = {
				pre_sync = {}
			}
		}
	}
}
	`, name, name, name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_job_hook.jh1", "id"),
				),
			},
		},
	})
}

func TestAcc_JobHook_Update(t *testing.T) {
	name := acctest.RandString(10)

	configInit := fmt.Sprintf(`
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

resource "neosync_job_hook" "hook1" {
	name = "%[1]s"
	description = "Initial description"
	job_id = neosync_job.job1.id
	enabled = true
	priority = 1
	config = {
		sql = {
			query = "SELECT 1;"
			connection_id = neosync_connection.source.id
			timing = {
				pre_sync = {}
			}
		}
	}
}`, name)

	configUpdate := fmt.Sprintf(`
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

resource "neosync_job_hook" "hook1" {
	name = "%[1]s-updated"
	description = "Updated description"
	job_id = neosync_job.job1.id
	enabled = false
	priority = 2
	config = {
		sql = {
			query = "SELECT 2;"
			connection_id = neosync_connection.source.id
			timing = {
				post_sync = {}
			}
		}
	}
}`, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configInit,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "name", name),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "description", "Initial description"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "enabled", "true"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "priority", "1"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "config.sql.query", "SELECT 1;"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "name", name+"-updated"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "description", "Updated description"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "enabled", "false"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "priority", "2"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "config.sql.query", "SELECT 2;"),
				),
			},
		},
	})
}

// Test multiple hooks on same job.
func TestAcc_JobHook_Multiple(t *testing.T) {
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

resource "neosync_job_hook" "hook1" {
	name = "%[1]s-1"
	description = "First hook"
	job_id = neosync_job.job1.id
	enabled = true
	priority = 1
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

resource "neosync_job_hook" "hook2" {
	name = "%[1]s-2"
	description = "Second hook"
	job_id = neosync_job.job1.id
	enabled = true
	priority = 2
	config = {
		sql = {
			query = "SELECT 2;"
			connection_id = neosync_connection.source.id
			timing = {
				post_sync = {}
			}
		}
	}
}`, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "name", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttr("neosync_job_hook.hook1", "priority", "1"),
					resource.TestCheckResourceAttr("neosync_job_hook.hook2", "name", fmt.Sprintf("%s-2", name)),
					resource.TestCheckResourceAttr("neosync_job_hook.hook2", "priority", "2"),
				),
			},
		},
	})
}

// Test importing job hook.
func TestAcc_JobHook_Import(t *testing.T) {
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

resource "neosync_job_hook" "hook1" {
	name = "%[1]s"
	description = "Importable hook"
	job_id = neosync_job.job1.id
	enabled = true
	priority = 1
	config = {
		sql = {
			query = "SELECT 1;"
			connection_id = neosync_connection.source.id
			timing = {
				pre_sync = {}
			}
		}
	}
}`, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "neosync_job_hook.hook1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
