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

func TestAcc_Job_Pg_Pg_Mappings(t *testing.T) {
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

	config1 := fmt.Sprintf(`
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
		},
		{
			schema = "public"
			table = "users"
			column = "name"
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

	config2 := fmt.Sprintf(`
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
			column = "name"
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

	var accountID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_job.job1", "id"),
					resource.TestCheckResourceAttrSet("neosync_job.job1", "account_id"),
					GetAccountIdFromState("neosync_job.job1", func(accountId string) { accountID = accountId }),
				),
			},
			{
				Config: config1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_job.job1", "id"),
					resource.TestCheckResourceAttrSet("neosync_job.job1", "account_id"),
					GetTestAccountIdFromStateFn("neosync_job.job1", func() string { return accountID }),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_job.job1", "id"),
					resource.TestCheckResourceAttrSet("neosync_job.job1", "account_id"),
					GetTestAccountIdFromStateFn("neosync_job.job1", func() string { return accountID }),
				),
			},
		},
	})
}

func TestAcc_Job_Pg_Pg_Destinations(t *testing.T) {
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

	config1 := fmt.Sprintf(`
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
resource "neosync_connection" "destination2" {
	name = "%s-dest2"

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
		},
		{
			connection_id = neosync_connection.destination2.id
			postgres = {
				init_table_schema = false
				truncate_table = {
					truncate_before_insert = true
					cascade = true
				}
			}
		},
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
	`, name, name, name, name)

	config2 := fmt.Sprintf(`
resource "neosync_connection" "source" {
	name = "%s-src"

	postgres = {
		url = "test-url"
	}
}
resource "neosync_connection" "destination2" {
	name = "%s-dest2"

	postgres = {
		url = "test-url2"
	}
}
resource "neosync_connection" "destination3" {
	name = "%s-dest3"

	postgres = {
		url = "test-url3"
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
			connection_id = neosync_connection.destination2.id
			postgres = {
				init_table_schema = true
				truncate_table = {
					truncate_before_insert = true
					cascade = true
				}
			}
		},
		{
			connection_id = neosync_connection.destination3.id
			postgres = {
				init_table_schema = false
				truncate_table = {
					truncate_before_insert = true
					cascade = true
				}
			}
		},
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
	`, name, name, name, name)

	_ = config2

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
			{
				Config: config1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_job.job1", "id"),
				),
			},
			// {
			// 	Config: config2,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttrSet("neosync_job.job1", "id"),
			// 	),
			// },
		},
	})
}

func TestAcc_Job_Mysql_Mysql(t *testing.T) {
	name := acctest.RandString(10)

	config := fmt.Sprintf(`
resource "neosync_connection" "source" {
	name = "%s-src"

	mysql = {
		url = "test-url"
	}
}
resource "neosync_connection" "destination" {
	name = "%s-dest"

	mysql = {
		url = "test-url2"
	}
}

resource "neosync_job" "job1" {
	name = "%s"
	source = {
		mysql = {
			halt_on_new_column_addition = false
			connection_id = neosync_connection.source.id
		}
	}
	destinations = [
		{
			connection_id = neosync_connection.destination.id
			mysql = {
				init_table_schema = false
				truncate_table = {
					truncate_before_insert = true
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

func TestAcc_Job_Mysql_Aws(t *testing.T) {
	name := acctest.RandString(10)

	config := fmt.Sprintf(`
resource "neosync_connection" "source" {
	name = "%s-src"

	mysql = {
		url = "test-url"
	}
}
resource "neosync_connection" "destination" {
	name = "%s-dest"

	aws_s3 = {
		bucket = "test-bucket"
	}
}

resource "neosync_job" "job1" {
	name = "%s"
	source = {
		mysql = {
			halt_on_new_column_addition = false
			connection_id = neosync_connection.source.id
		}
	}
	destinations = [
		{
			connection_id = neosync_connection.destination.id
			aws_s3 = {}
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

func TestAcc_Job_Generate_Pg(t *testing.T) {
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
	`, name, name)

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
