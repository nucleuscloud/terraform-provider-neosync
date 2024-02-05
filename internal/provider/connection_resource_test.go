package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAcc_Connection_Postgres_Url(t *testing.T) {
	connectionName := acctest.RandString(10)

	testAccConnectionConfig := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	postgres = {
		url = "test-url"
	}
}
`, connectionName)
	testAccConnectionConfigUpdated := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	postgres = {
		url = "test-url2"
	}
}
`, connectionName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_connection.test1", "id"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "name", connectionName),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.url", "test-url"),
				),
			},
			{
				Config: testAccConnectionConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.url", "test-url2"),
				),
			},
		},
	})
}

func TestAcc_Connection_Postgres_Connection(t *testing.T) {
	connectionName := acctest.RandString(10)

	testAccConnectionConfig := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	postgres = {
		host = "test-url"
		port = 5432
		name = "neosync"
		user = "postgres"
		pass = "postgres123"
		ssl_mode = "disable"

		tunnel = {
			host = "localhost"
			port = 22
			user = "test"
			known_host_public_key = "123"
			private_key = "my-private-key"
			passphrase = "test"
		}
	}
}
`, connectionName)
	testAccConnectionConfigUpdated := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	postgres = {
		host = "test-url"
		port = 5432
		name = "neosync"
		user = "postgres"
		pass = "postgres123"
		ssl_mode = "disable"

		tunnel = {
			host = "localhost"
			port = 22
			user = "test"
			known_host_public_key = "111"
			private_key = "my-private-key"
			passphrase = "test"
		}
	}
}
`, connectionName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_connection.test1", "id"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "name", connectionName),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.host", "test-url"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.port", "5432"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.name", "neosync"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.user", "postgres"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.pass", "postgres123"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.ssl_mode", "disable"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.host", "localhost"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.port", "22"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.user", "test"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.known_host_public_key", "123"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.private_key", "my-private-key"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.passphrase", "test"),
				),
			},
			{
				Config: testAccConnectionConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.known_host_public_key", "111"),
				),
			},
		},
	})
}

func TestAcc_Connection_Import(t *testing.T) {
	connectionName := acctest.RandString(10)
	testAccConnectionConfig := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	postgres = {
		host = "test-url"
		port = 5432
		name = "neosync"
		user = "postgres"
		pass = "postgres123"
		ssl_mode = "disable"

		tunnel = {
			host = "localhost"
			port = 22
			user = "test"
			known_host_public_key = "123"
			private_key = "my-private-key"
			passphrase = "test"
		}
	}
}
`, connectionName)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_connection.test1", "id"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "name", connectionName),
				),
			},
			{
				ResourceName:      "neosync_connection.test1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getPrimaryImportId("neosync_connection.test1"),
			},
		},
	})
}

func TestAcc_Connection_Mysql_Url(t *testing.T) {
	connectionName := acctest.RandString(10)

	testAccConnectionConfig := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	mysql = {
		url = "test-url"
	}
}
`, connectionName)
	testAccConnectionConfigUpdated := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	mysql = {
		url = "test-url2"
	}
}
`, connectionName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_connection.test1", "id"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "name", connectionName),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.url", "test-url"),
				),
			},
			{
				Config: testAccConnectionConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.url", "test-url2"),
				),
			},
		},
	})
}

func TestAcc_Connection_Mysql_Connection(t *testing.T) {
	connectionName := acctest.RandString(10)

	testAccConnectionConfig := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	mysql = {
		host = "test-url"
		port = 5432
		name = "neosync"
		user = "myuser"
		pass = "123"
		protocol = "tcp"

		tunnel = {
			host = "localhost"
			port = 22
			user = "test"
			known_host_public_key = "123"
			private_key = "my-private-key"
			passphrase = "test"
		}
	}
}
`, connectionName)
	testAccConnectionConfigUpdated := fmt.Sprintf(`
resource "neosync_connection" "test1" {
  name = "%s"

	postgres = {
		host = "test-url"
		port = 5432
		name = "neosync"
		user = "myuser"
		pass = "123"
		protocol = "tcp"

		tunnel = {
			host = "localhost"
			port = 22
			user = "test"
			known_host_public_key = "111"
			private_key = "my-private-key"
			passphrase = "test"
		}
	}
}
`, connectionName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_connection.test1", "id"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "name", connectionName),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.host", "test-url"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.port", "5432"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.name", "neosync"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.user", "myuser"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.pass", "123"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.protocol", "tcp"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.tunnel.host", "localhost"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.tunnel.port", "22"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.tunnel.user", "test"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.tunnel.known_host_public_key", "123"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.tunnel.private_key", "my-private-key"),
					resource.TestCheckResourceAttr("neosync_connection.test1", "mysql.tunnel.passphrase", "test"),
				),
			},
			{
				Config: testAccConnectionConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neosync_connection.test1", "postgres.tunnel.known_host_public_key", "111"),
				),
			},
		},
	})
}

func getPrimaryImportId(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("no ID is set")
		}

		return rs.Primary.ID, nil
	}
}
