# Postgres connection
resource "neosync_connection" "local_pg" {
  name = "local-pg"

  postgres = {
    host     = "localhost"
    port     = 5432
    name     = "postgres"
    user     = "postgres"
    pass     = "postgres"
    ssl_mode = "disable"

    tunnel = {
      host                  = "localhost"
      port                  = 22
      user                  = "test"
      known_host_public_key = "123"
      private_key           = "my-private-key"
      passphrase            = "test"
    }
  }
}

# Postgres connection with tunnel
resource "neosync_connection" "private_db" {
  name = "private-pg"

  postgres = {
    host     = "my-private-db.example.com"
    port     = 5432
    name     = "postgres"
    user     = "postgres"
    pass     = "postgres"
    ssl_mode = "disable"

    tunnel = {
      host                  = "bastion.example.com"
      port                  = 22
      user                  = "test"
      known_host_public_key = "123"
      private_key           = "my-private-key"
      passphrase            = "test"
    }
  }
}

# Mysql connection with tunnel
resource "neosync_connection" "private_mysql_db" {
  name = "private-mysql"

  mysql = {
    host     = "my-private-db.example.com"
    port     = 3306
    name     = "mysql"
    user     = "mysql"
    pass     = "mysql"
    protocol = "tcp"

    tunnel = {
      host                  = "bastion.example.com"
      port                  = 22
      user                  = "test"
      known_host_public_key = "123"
      private_key           = "my-private-key"
      passphrase            = "test"
    }
  }
}
