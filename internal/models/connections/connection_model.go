package connection_model

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
)

type ConnectionResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	AccountId types.String `tfsdk:"account_id"`

	Postgres *Postgres `tfsdk:"postgres"`
	Mysql    *Mysql    `tfsdk:"mysql"`
	AwsS3    *AwsS3    `tfsdk:"aws_s3"`
}

type AwsS3 struct {
	Bucket      types.String    `tfsdk:"bucket"`
	PathPrefix  types.String    `tfsdk:"path_prefix"`
	Region      types.String    `tfsdk:"region"`
	Endpoint    types.String    `tfsdk:"endpoint"`
	Credentials *AwsCredentials `tfsdk:"credentials"`
}

type AwsCredentials struct {
	Profile         types.String `tfsdk:"profile"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	SessionToken    types.String `tfsdk:"session_token"`
	FromEc2Role     types.Bool   `tfsdk:"from_ec2_role"`
	RoleArn         types.String `tfsdk:"role_arn"`
	RoleExternalId  types.String `tfsdk:"role_external_id"`
}

type Postgres struct {
	Url types.String `tfsdk:"url"`

	Host    types.String `tfsdk:"host"`
	Port    types.Int64  `tfsdk:"port"`
	Name    types.String `tfsdk:"name"`
	User    types.String `tfsdk:"user"`
	Pass    types.String `tfsdk:"pass"`
	SslMode types.String `tfsdk:"ssl_mode"`

	Tunnel *SSHTunnel `tfsdk:"tunnel"`

	ClientTls         *ClientTlsConfig      `tfsdk:"client_tls"`
	ConnectionOptions *SqlConnectionOptions `tfsdk:"connection_options"`
}

type Mysql struct {
	Url types.String `tfsdk:"url"`

	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
	Name     types.String `tfsdk:"name"`
	User     types.String `tfsdk:"user"`
	Pass     types.String `tfsdk:"pass"`
	Protocol types.String `tfsdk:"protocol"`

	Tunnel *SSHTunnel `tfsdk:"tunnel"`

	ClientTls         *ClientTlsConfig      `tfsdk:"client_tls"`
	ConnectionOptions *SqlConnectionOptions `tfsdk:"connection_options"`
}

type SSHTunnel struct {
	Host               types.String `tfsdk:"host"`
	Port               types.Int64  `tfsdk:"port"`
	User               types.String `tfsdk:"user"`
	KnownHostPublicKey types.String `tfsdk:"known_host_public_key"`

	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
}

type SqlConnectionOptions struct {
	MaxOpenConnections *int32  `tfsdk:"max_open_connections"`
	MaxIdleConnections *int32  `tfsdk:"max_idle_connections"`
	MaxIdleDuration    *string `tfsdk:"max_idle_duration"`
	MaxOpenDuration    *string `tfsdk:"max_open_duration"`
}

func (c *SqlConnectionOptions) ToDto() (*mgmtv1alpha1.SqlConnectionOptions, error) {
	if c == nil {
		return nil, errors.New("sql connection options is nil")
	}

	return &mgmtv1alpha1.SqlConnectionOptions{
		MaxConnectionLimit: c.MaxOpenConnections,
		MaxIdleConnections: c.MaxIdleConnections,
		MaxIdleDuration:    c.MaxIdleDuration,
		MaxOpenDuration:    c.MaxOpenDuration,
	}, nil
}

func (c *SqlConnectionOptions) FromDto(dto *mgmtv1alpha1.SqlConnectionOptions) error {
	if c == nil {
		return errors.New("sql connection options is nil")
	}

	if dto == nil {
		return errors.New("sql connection options dto is nil")
	}

	c.MaxOpenConnections = dto.MaxConnectionLimit
	c.MaxIdleConnections = dto.MaxIdleConnections
	c.MaxIdleDuration = dto.MaxIdleDuration
	c.MaxOpenDuration = dto.MaxOpenDuration

	return nil
}

type ClientTlsConfig struct {
	RootCert   *string `tfsdk:"root_cert"`
	ClientCert *string `tfsdk:"client_cert"`
	ClientKey  *string `tfsdk:"client_key"`
	ServerName *string `tfsdk:"server_name"`
}

func (c *ClientTlsConfig) ToDto() (*mgmtv1alpha1.ClientTlsConfig, error) {
	if c == nil {
		return nil, errors.New("client tls config is nil")
	}

	return &mgmtv1alpha1.ClientTlsConfig{
		RootCert:   c.RootCert,
		ClientCert: c.ClientCert,
		ClientKey:  c.ClientKey,
		ServerName: c.ServerName,
	}, nil
}

func (c *ClientTlsConfig) FromDto(dto *mgmtv1alpha1.ClientTlsConfig) error {
	if c == nil {
		return errors.New("client tls config is nil")
	}

	c.RootCert = dto.RootCert
	c.ClientCert = dto.ClientCert
	c.ClientKey = dto.ClientKey
	c.ServerName = dto.ServerName

	return nil
}

func (c *ConnectionResourceModel) ToCreateConnectionDto() (*mgmtv1alpha1.CreateConnectionRequest, error) {
	if c == nil {
		return nil, errors.New("connection resource model is nil")
	}

	connConfig, err := c.ToConnectionConfigDto()
	if err != nil {
		return nil, err
	}

	return &mgmtv1alpha1.CreateConnectionRequest{
		Name:             c.Name.ValueString(),
		AccountId:        c.AccountId.ValueString(),
		ConnectionConfig: connConfig,
	}, nil
}

func (c *ConnectionResourceModel) ToUpdateConnectionDto() (*mgmtv1alpha1.UpdateConnectionRequest, error) {
	if c == nil {
		return nil, errors.New("connection resource model is nil")
	}

	connConfig, err := c.ToConnectionConfigDto()
	if err != nil {
		return nil, err
	}

	return &mgmtv1alpha1.UpdateConnectionRequest{
		Id:               c.Id.ValueString(),
		Name:             c.Name.ValueString(),
		ConnectionConfig: connConfig,
	}, nil
}

func (c *ConnectionResourceModel) FromDto(dto *mgmtv1alpha1.Connection) error {
	if c == nil {
		return errors.New("connection resource model is nil")
	}

	if dto == nil {
		return errors.New("connection dto is nil")
	}

	c.Id = types.StringValue(dto.Id)
	c.Name = types.StringValue(dto.Name)
	c.AccountId = types.StringValue(dto.AccountId)

	if dto.ConnectionConfig != nil {
		if err := c.FromConnectionConfigDto(dto.GetConnectionConfig()); err != nil {
			return err
		}
	}

	return nil
}

func (c *ConnectionResourceModel) ToConnectionConfigDto() (*mgmtv1alpha1.ConnectionConfig, error) {
	if c == nil {
		return nil, errors.New("connection resource model is nil")
	}

	if c.Postgres != nil {
		var tunnel *mgmtv1alpha1.SSHTunnel
		if c.Postgres.Tunnel != nil {
			tunnelDto, err := c.Postgres.Tunnel.ToDto()
			if err != nil {
				return nil, err
			}
			tunnel = tunnelDto
		}
		var clientTls *mgmtv1alpha1.ClientTlsConfig
		if c.Postgres.ClientTls != nil {
			clientTlsDto, err := c.Postgres.ClientTls.ToDto()
			if err != nil {
				return nil, err
			}
			clientTls = clientTlsDto
		}
		var connectionOptions *mgmtv1alpha1.SqlConnectionOptions
		if c.Postgres.ConnectionOptions != nil {
			connectionOptionsDto, err := c.Postgres.ConnectionOptions.ToDto()
			if err != nil {
				return nil, err
			}
			connectionOptions = connectionOptionsDto
		}
		if c.Postgres.Url.ValueString() != "" {
			return &mgmtv1alpha1.ConnectionConfig{
				Config: &mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Url{
							Url: c.Postgres.Url.ValueString(),
						},
						Tunnel:            tunnel,
						ConnectionOptions: connectionOptions,
						ClientTls:         clientTls,
					},
				},
			}, nil
		} else {
			pg := c.Postgres
			if pg.Host.ValueString() == "" || pg.Port.ValueInt64() == 0 || pg.Name.ValueString() == "" || pg.User.ValueString() == "" || pg.Pass.ValueString() == "" {
				return nil, fmt.Errorf("invalid postgres config")
			}
			return &mgmtv1alpha1.ConnectionConfig{
				Config: &mgmtv1alpha1.ConnectionConfig_PgConfig{
					PgConfig: &mgmtv1alpha1.PostgresConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.PostgresConnectionConfig_Connection{
							Connection: &mgmtv1alpha1.PostgresConnection{
								Host:    pg.Host.ValueString(),
								Port:    int32(pg.Port.ValueInt64()),
								Name:    pg.Name.ValueString(),
								User:    pg.User.ValueString(),
								Pass:    pg.Pass.ValueString(),
								SslMode: pg.SslMode.ValueStringPointer(),
							},
						},
						Tunnel:            tunnel,
						ConnectionOptions: connectionOptions,
						ClientTls:         clientTls,
					},
				},
			}, nil
		}
	}
	if c.Mysql != nil {
		var tunnel *mgmtv1alpha1.SSHTunnel
		if c.Mysql.Tunnel != nil {
			tunnelDto, err := c.Mysql.Tunnel.ToDto()
			if err != nil {
				return nil, err
			}
			tunnel = tunnelDto
		}
		var clientTls *mgmtv1alpha1.ClientTlsConfig
		if c.Mysql.ClientTls != nil {
			clientTlsDto, err := c.Mysql.ClientTls.ToDto()
			if err != nil {
				return nil, err
			}
			clientTls = clientTlsDto
		}
		var connectionOptions *mgmtv1alpha1.SqlConnectionOptions
		if c.Mysql.ConnectionOptions != nil {
			connectionOptionsDto, err := c.Mysql.ConnectionOptions.ToDto()
			if err != nil {
				return nil, err
			}
			connectionOptions = connectionOptionsDto
		}
		if c.Mysql.Url.ValueString() != "" {
			return &mgmtv1alpha1.ConnectionConfig{
				Config: &mgmtv1alpha1.ConnectionConfig_MysqlConfig{
					MysqlConfig: &mgmtv1alpha1.MysqlConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.MysqlConnectionConfig_Url{
							Url: c.Mysql.Url.ValueString(),
						},
						Tunnel:            tunnel,
						ConnectionOptions: connectionOptions,
						ClientTls:         clientTls,
					},
				},
			}, nil
		} else {
			mysql := c.Mysql
			if mysql.Host.ValueString() == "" || mysql.Port.ValueInt64() == 0 || mysql.Name.ValueString() == "" || mysql.User.ValueString() == "" || mysql.Pass.ValueString() == "" || mysql.Protocol.ValueString() == "" {
				return nil, fmt.Errorf("invalid mysql config")
			}
			return &mgmtv1alpha1.ConnectionConfig{
				Config: &mgmtv1alpha1.ConnectionConfig_MysqlConfig{
					MysqlConfig: &mgmtv1alpha1.MysqlConnectionConfig{
						ConnectionConfig: &mgmtv1alpha1.MysqlConnectionConfig_Connection{
							Connection: &mgmtv1alpha1.MysqlConnection{
								Host:     mysql.Host.ValueString(),
								Port:     int32(mysql.Port.ValueInt64()),
								Name:     mysql.Name.ValueString(),
								User:     mysql.User.ValueString(),
								Pass:     mysql.Pass.ValueString(),
								Protocol: mysql.Protocol.ValueString(),
							},
						},
						Tunnel:            tunnel,
						ConnectionOptions: connectionOptions,
						ClientTls:         clientTls,
					},
				},
			}, nil
		}
	}
	if c.AwsS3 != nil {
		var creds *mgmtv1alpha1.AwsS3Credentials
		if c.AwsS3.Credentials != nil {
			creds = &mgmtv1alpha1.AwsS3Credentials{
				Profile:         c.AwsS3.Credentials.Profile.ValueStringPointer(),
				AccessKeyId:     c.AwsS3.Credentials.AccessKeyId.ValueStringPointer(),
				SecretAccessKey: c.AwsS3.Credentials.SecretAccessKey.ValueStringPointer(),
				SessionToken:    c.AwsS3.Credentials.SessionToken.ValueStringPointer(),
				FromEc2Role:     c.AwsS3.Credentials.FromEc2Role.ValueBoolPointer(),
				RoleArn:         c.AwsS3.Credentials.RoleArn.ValueStringPointer(),
				RoleExternalId:  c.AwsS3.Credentials.RoleExternalId.ValueStringPointer(),
			}
		}
		return &mgmtv1alpha1.ConnectionConfig{
			Config: &mgmtv1alpha1.ConnectionConfig_AwsS3Config{
				AwsS3Config: &mgmtv1alpha1.AwsS3ConnectionConfig{
					Bucket:      c.AwsS3.Bucket.ValueString(),
					PathPrefix:  c.AwsS3.PathPrefix.ValueStringPointer(),
					Region:      c.AwsS3.Region.ValueStringPointer(),
					Endpoint:    c.AwsS3.Endpoint.ValueStringPointer(),
					Credentials: creds,
				},
			},
		}, nil
	}

	return nil, errors.New("unable to find a config to hydrate connection resource model")
}

func (c *ConnectionResourceModel) FromConnectionConfigDto(dto *mgmtv1alpha1.ConnectionConfig) error {
	switch config := dto.GetConfig().(type) {
	case *mgmtv1alpha1.ConnectionConfig_PgConfig:
		switch pgcc := config.PgConfig.ConnectionConfig.(type) {
		case *mgmtv1alpha1.PostgresConnectionConfig_Connection:
			c.Postgres = &Postgres{
				Host:    types.StringValue(pgcc.Connection.Host),
				Port:    types.Int64Value(int64(pgcc.Connection.Port)),
				Name:    types.StringValue(pgcc.Connection.Name),
				User:    types.StringValue(pgcc.Connection.User),
				Pass:    types.StringValue(pgcc.Connection.Pass),
				SslMode: types.StringPointerValue(pgcc.Connection.SslMode),
				Tunnel:  nil,
			}
			if config.PgConfig.Tunnel != nil {
				c.Postgres.Tunnel = &SSHTunnel{}
				if err := c.Postgres.Tunnel.FromDto(config.PgConfig.Tunnel); err != nil {
					return err
				}
			}
			if config.PgConfig.ConnectionOptions != nil {
				c.Postgres.ConnectionOptions = &SqlConnectionOptions{}
				if err := c.Postgres.ConnectionOptions.FromDto(config.PgConfig.ConnectionOptions); err != nil {
					return err
				}
			}
			if config.PgConfig.ClientTls != nil {
				c.Postgres.ClientTls = &ClientTlsConfig{}
				if err := c.Postgres.ClientTls.FromDto(config.PgConfig.ClientTls); err != nil {
					return err
				}
			}
			return nil
		case *mgmtv1alpha1.PostgresConnectionConfig_Url:
			c.Postgres = &Postgres{
				Url:    types.StringValue(pgcc.Url),
				Tunnel: nil,
			}
			if config.PgConfig.Tunnel != nil {
				c.Postgres.Tunnel = &SSHTunnel{}
				if err := c.Postgres.Tunnel.FromDto(config.PgConfig.Tunnel); err != nil {
					return err
				}
			}
			if config.PgConfig.ConnectionOptions != nil {
				c.Postgres.ConnectionOptions = &SqlConnectionOptions{}
				if err := c.Postgres.ConnectionOptions.FromDto(config.PgConfig.ConnectionOptions); err != nil {
					return err
				}
			}
			if config.PgConfig.ClientTls != nil {
				c.Postgres.ClientTls = &ClientTlsConfig{}
				if err := c.Postgres.ClientTls.FromDto(config.PgConfig.ClientTls); err != nil {
					return err
				}
			}
			return nil
		default:
			return errors.New("unable to find a config to hydrate connection resource model")
		}
	case *mgmtv1alpha1.ConnectionConfig_MysqlConfig:
		switch mycc := config.MysqlConfig.ConnectionConfig.(type) {
		case *mgmtv1alpha1.MysqlConnectionConfig_Connection:
			c.Mysql = &Mysql{
				Host:     types.StringValue(mycc.Connection.Host),
				Port:     types.Int64Value(int64(mycc.Connection.Port)),
				Name:     types.StringValue(mycc.Connection.Name),
				User:     types.StringValue(mycc.Connection.User),
				Pass:     types.StringValue(mycc.Connection.Pass),
				Protocol: types.StringValue(mycc.Connection.Protocol),
				Tunnel:   nil,
			}
			if config.MysqlConfig.Tunnel != nil {
				c.Mysql.Tunnel = &SSHTunnel{}
				if err := c.Mysql.Tunnel.FromDto(config.MysqlConfig.Tunnel); err != nil {
					return err
				}
			}
			if config.MysqlConfig.ConnectionOptions != nil {
				c.Mysql.ConnectionOptions = &SqlConnectionOptions{}
				if err := c.Mysql.ConnectionOptions.FromDto(config.MysqlConfig.ConnectionOptions); err != nil {
					return err
				}
			}
			if config.MysqlConfig.ClientTls != nil {
				c.Mysql.ClientTls = &ClientTlsConfig{}
				if err := c.Mysql.ClientTls.FromDto(config.MysqlConfig.ClientTls); err != nil {
					return err
				}
			}
			return nil
		case *mgmtv1alpha1.MysqlConnectionConfig_Url:
			c.Mysql = &Mysql{
				Url:    types.StringValue(mycc.Url),
				Tunnel: nil,
			}
			if config.MysqlConfig.Tunnel != nil {
				c.Mysql.Tunnel = &SSHTunnel{}
				if err := c.Mysql.Tunnel.FromDto(config.MysqlConfig.Tunnel); err != nil {
					return err
				}
			}
			if config.MysqlConfig.ConnectionOptions != nil {
				c.Mysql.ConnectionOptions = &SqlConnectionOptions{}
				if err := c.Mysql.ConnectionOptions.FromDto(config.MysqlConfig.ConnectionOptions); err != nil {
					return err
				}
			}
			if config.MysqlConfig.ClientTls != nil {
				c.Mysql.ClientTls = &ClientTlsConfig{}
				if err := c.Mysql.ClientTls.FromDto(config.MysqlConfig.ClientTls); err != nil {
					return err
				}
			}
			return nil
		default:
			return errors.New("unable to findconfig to hydrate connection resource model")
		}
	case *mgmtv1alpha1.ConnectionConfig_AwsS3Config:
		c.AwsS3 = &AwsS3{
			Bucket:     types.StringValue(config.AwsS3Config.Bucket),
			PathPrefix: types.StringPointerValue(config.AwsS3Config.PathPrefix),
			Region:     types.StringPointerValue(config.AwsS3Config.Region),
			Endpoint:   types.StringPointerValue(config.AwsS3Config.Endpoint),
		}
		if !isAwsCredentialsEmpty(config.AwsS3Config.Credentials) {
			c.AwsS3.Credentials = &AwsCredentials{}
			c.AwsS3.Credentials.Profile = types.StringPointerValue(config.AwsS3Config.Credentials.Profile)
			c.AwsS3.Credentials.AccessKeyId = types.StringPointerValue(config.AwsS3Config.Credentials.AccessKeyId)
			c.AwsS3.Credentials.SecretAccessKey = types.StringPointerValue(config.AwsS3Config.Credentials.SecretAccessKey)
			c.AwsS3.Credentials.SessionToken = types.StringPointerValue(config.AwsS3Config.Credentials.SessionToken)
			c.AwsS3.Credentials.FromEc2Role = types.BoolPointerValue(config.AwsS3Config.Credentials.FromEc2Role)
			c.AwsS3.Credentials.RoleArn = types.StringPointerValue(config.AwsS3Config.Credentials.RoleArn)
			c.AwsS3.Credentials.RoleExternalId = types.StringPointerValue(config.AwsS3Config.Credentials.RoleExternalId)
		}
		return nil
	default:
		return errors.New("unable to find a config to hydrate connection resource model")
	}
}

func (s *SSHTunnel) FromDto(dto *mgmtv1alpha1.SSHTunnel) error {
	if s == nil {
		return errors.New("ssh tunnel is nil")
	}

	if dto == nil {
		return errors.New("ssh tunnel dto is nil")
	}

	s.Host = types.StringValue(dto.Host)
	s.Port = types.Int64Value(int64(dto.Port))
	s.User = types.StringValue(dto.User)
	s.KnownHostPublicKey = types.StringPointerValue(dto.KnownHostPublicKey)

	if dto.Authentication != nil {
		switch auth := dto.Authentication.AuthConfig.(type) {
		case *mgmtv1alpha1.SSHAuthentication_PrivateKey:
			s.PrivateKey = types.StringValue(auth.PrivateKey.Value)
			s.Passphrase = types.StringPointerValue(auth.PrivateKey.Passphrase)
		case *mgmtv1alpha1.SSHAuthentication_Passphrase:
			s.Passphrase = types.StringValue(auth.Passphrase.Value)
		}
	}

	return nil
}

func (s *SSHTunnel) ToDto() (*mgmtv1alpha1.SSHTunnel, error) {
	if s == nil {
		return nil, errors.New("ssh tunnel is nil")
	}

	var auth *mgmtv1alpha1.SSHAuthentication
	if s.PrivateKey.ValueString() != "" {
		auth = &mgmtv1alpha1.SSHAuthentication{
			AuthConfig: &mgmtv1alpha1.SSHAuthentication_PrivateKey{
				PrivateKey: &mgmtv1alpha1.SSHPrivateKey{
					Value:      s.PrivateKey.ValueString(),
					Passphrase: s.Passphrase.ValueStringPointer(),
				},
			},
		}
	} else if s.Passphrase.ValueString() != "" {
		auth = &mgmtv1alpha1.SSHAuthentication{
			AuthConfig: &mgmtv1alpha1.SSHAuthentication_Passphrase{
				Passphrase: &mgmtv1alpha1.SSHPassphrase{
					Value: s.Passphrase.ValueString(),
				},
			},
		}
	}

	return &mgmtv1alpha1.SSHTunnel{
		Host:               s.Host.ValueString(),
		Port:               int32(s.Port.ValueInt64()),
		User:               s.User.ValueString(),
		KnownHostPublicKey: s.KnownHostPublicKey.ValueStringPointer(),
		Authentication:     auth,
	}, nil
}

func isAwsCredentialsEmpty(creds *mgmtv1alpha1.AwsS3Credentials) bool {
	if creds == nil {
		return true
	}
	return creds.Profile == nil && creds.AccessKeyId == nil && creds.SecretAccessKey == nil && creds.SessionToken == nil &&
		creds.FromEc2Role == nil && creds.RoleArn == nil && creds.RoleExternalId == nil
}
