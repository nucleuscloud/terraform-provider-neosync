---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neosync_connection Resource - terraform-provider-neosync"
subcategory: ""
description: |-
  Neosync Connection resource
---

# neosync_connection (Resource)

Neosync Connection resource

## Example Usage

```terraform
# Postgres Connection with URL
resource "neosync_connection" "pg_url" {
  name = "local-pg"

  postgres = {
    url = "postgres://postgres:postgres@localhost:5432?sslmode=disable"
  }
}

# Postgres connection via separate parts
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

    connection_options = {
      max_idle_connections = 10
      max_open_connections = 20
      max_idle_duration    = "3600"
      max_open_duration    = "3700"
    }

    client_tls = {
      root_cert   = "my-root-cert"
      client_cert = "my-client-cert"
      client_key  = "my-client-key"
      server_name = "my-server-name"
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

# AWS S3 Connection
resource "neosync_connection" "job_bucket" {
  name = "stage-backups"

  aws_s3 = {
    bucket      = "my-company-bucket"
    path_prefix = "/neosync"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The unique friendly name of the connection

### Optional

- `account_id` (String) The unique identifier of the account. Can be pulled from the API Key if present, or must be specified if using a user access token
- `aws_s3` (Attributes) The aws s3 bucket that will be associated with this connection (see [below for nested schema](#nestedatt--aws_s3))
- `mysql` (Attributes) The mysql database that will be associated with this connection (see [below for nested schema](#nestedatt--mysql))
- `postgres` (Attributes) The postgres database that will be associated with this connection (see [below for nested schema](#nestedatt--postgres))

### Read-Only

- `id` (String) The unique identifier of the connection

<a id="nestedatt--aws_s3"></a>
### Nested Schema for `aws_s3`

Required:

- `bucket` (String) The name of the S3 bucket

Optional:

- `credentials` (Attributes) Credentials that may be necessary to access the S3 bucket in a R/W fashion (see [below for nested schema](#nestedatt--aws_s3--credentials))
- `endpoint` (String) The endpoint that will be used by the SDK to access the bucket
- `path_prefix` (String) The folder within the bucket that the connection will be scoped to
- `region` (String) The region that will be used by the SDK to access the bucket

<a id="nestedatt--aws_s3--credentials"></a>
### Nested Schema for `aws_s3.credentials`

Optional:

- `access_key_id` (String) The AWS access key id
- `from_ec2_role` (Boolean) Will result in the sync operations pulling from the EC2 role
- `profile` (String) The profile found in the ~/.aws/config that can be used to access credentials
- `role_arn` (String) The role arn that can be assumed
- `role_external_id` (String) The external id that will be provided when the role arn is assumed
- `secret_access_key` (String, Sensitive) The AWS secret access key
- `session_token` (String) The AWS session token



<a id="nestedatt--mysql"></a>
### Nested Schema for `mysql`

Optional:

- `client_tls` (Attributes) TLS configuration for the connection (see [below for nested schema](#nestedatt--mysql--client_tls))
- `connection_options` (Attributes) SQL connection options (see [below for nested schema](#nestedatt--mysql--connection_options))
- `host` (String) The host name of the mysql server
- `name` (String) The name of the database that will be connected to
- `pass` (String, Sensitive) The password that will be authenticated with
- `port` (Number) The port of the mysql server
- `protocol` (String) The protocol of the mysql server
- `tunnel` (Attributes) SSH tunnel that is used to access databases that are not publicly accessible to the internet (see [below for nested schema](#nestedatt--mysql--tunnel))
- `url` (String) Standard mysql url connection string.
- `user` (String) The name of the user that will be authenticated with

<a id="nestedatt--mysql--client_tls"></a>
### Nested Schema for `mysql.client_tls`

Optional:

- `client_cert` (String) The client certificate in PEM format
- `client_key` (String, Sensitive) The client key in PEM format
- `root_cert` (String) The root certificate in PEM format
- `server_name` (String) The expected server name


<a id="nestedatt--mysql--connection_options"></a>
### Nested Schema for `mysql.connection_options`

Optional:

- `max_idle_connections` (Number) The maximum number of idle connections to the database
- `max_idle_duration` (String) The maximum amount of time a connection may be reused
- `max_open_connections` (Number) The maximum number of open connections to the database
- `max_open_duration` (String) The maximum amount of time a connection may be idle


<a id="nestedatt--mysql--tunnel"></a>
### Nested Schema for `mysql.tunnel`

Required:

- `host` (String) The host name of the server
- `port` (Number) The post of the ssh server
- `user` (String) The name of the user that will be authenticated with

Optional:

- `known_host_public_key` (String) The known SSH public key of the tunnel server.
- `passphrase` (String, Sensitive) If not using key authentication, a password must be provided. If a private key is provided, but encrypted, provide the passphrase here as it will be used to decrypt the private key
- `private_key` (String, Sensitive) If using key authentication, this must be a pem encoded private key



<a id="nestedatt--postgres"></a>
### Nested Schema for `postgres`

Optional:

- `client_tls` (Attributes) TLS configuration for the connection (see [below for nested schema](#nestedatt--postgres--client_tls))
- `connection_options` (Attributes) SQL connection options (see [below for nested schema](#nestedatt--postgres--connection_options))
- `host` (String) The host name of the postgres server
- `name` (String) The name of the database that will be connected to
- `pass` (String, Sensitive) The password that will be authenticated with
- `port` (Number) The port of the postgres server
- `ssl_mode` (String) The SSL mode for the postgres server
- `tunnel` (Attributes) SSH tunnel that is used to access databases that are not publicly accessible to the internet (see [below for nested schema](#nestedatt--postgres--tunnel))
- `url` (String) Standard postgres url connection string. Must be uri compliant
- `user` (String) The name of the user that will be authenticated with

<a id="nestedatt--postgres--client_tls"></a>
### Nested Schema for `postgres.client_tls`

Optional:

- `client_cert` (String) The client certificate in PEM format
- `client_key` (String, Sensitive) The client key in PEM format
- `root_cert` (String) The root certificate in PEM format
- `server_name` (String) The expected server name


<a id="nestedatt--postgres--connection_options"></a>
### Nested Schema for `postgres.connection_options`

Optional:

- `max_idle_connections` (Number) The maximum number of idle connections to the database
- `max_idle_duration` (String) The maximum amount of time a connection may be reused
- `max_open_connections` (Number) The maximum number of open connections to the database
- `max_open_duration` (String) The maximum amount of time a connection may be idle


<a id="nestedatt--postgres--tunnel"></a>
### Nested Schema for `postgres.tunnel`

Required:

- `host` (String) The host name of the server
- `port` (Number) The post of the ssh server
- `user` (String) The name of the user that will be authenticated with

Optional:

- `known_host_public_key` (String) The known SSH public key of the tunnel server.
- `passphrase` (String, Sensitive) If not using key authentication, a password must be provided. If a private key is provided, but encrypted, provide the passphrase here as it will be used to decrypt the private key
- `private_key` (String, Sensitive) If using key authentication, this must be a pem encoded private key
