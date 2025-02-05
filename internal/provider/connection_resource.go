package provider

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	"github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1/mgmtv1alpha1connect"
	connection_model "github.com/nucleuscloud/terraform-provider-neosync/internal/models/connections"
)

var _ resource.Resource = &ConnectionResource{}
var _ resource.ResourceWithImportState = &ConnectionResource{}

func NewConnectionResource() resource.Resource {
	return &ConnectionResource{}
}

type ConnectionResource struct {
	client    mgmtv1alpha1connect.ConnectionServiceClient
	accountId *string
}

func (r *ConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

var (
	tunnelSchema = schema.SingleNestedAttribute{
		Description: "SSH tunnel that is used to access databases that are not publicly accessible to the internet",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The host name of the server",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The post of the ssh server",
				Required:    true,
			},
			"user": schema.StringAttribute{
				Description: "The name of the user that will be authenticated with",
				Required:    true,
			},
			"known_host_public_key": schema.StringAttribute{
				Description: "The known SSH public key of the tunnel server.",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "If using key authentication, this must be a pem encoded private key",
				Optional:    true,
				Sensitive:   true,
			},
			"passphrase": schema.StringAttribute{
				Description: "If not using key authentication, a password must be provided. If a private key is provided, but encrypted, provide the passphrase here as it will be used to decrypt the private key",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
)

func (r *ConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Neosync Connection resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique friendly name of the connection",
				Required:    true,
			},
			"account_id": schema.StringAttribute{
				Description:   "The unique identifier of the account. Can be pulled from the API Key if present, or must be specified if using a user access token",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"postgres": schema.SingleNestedAttribute{
				Description: "The postgres database that will be associated with this connection",
				Optional:    true,
				// PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description: "Standard postgres url connection string. Must be uri compliant",
						Optional:    true,
					},

					"host": schema.StringAttribute{
						Description: "The host name of the postgres server",
						Optional:    true,
					},
					"port": schema.Int64Attribute{
						Description: "The port of the postgres server",
						Optional:    true,
						// Default:     int64default.StaticInt64(5432),
					},
					"name": schema.StringAttribute{
						Description: "The name of the database that will be connected to",
						Optional:    true,
					},
					"user": schema.StringAttribute{
						Description: "The name of the user that will be authenticated with",
						Optional:    true,
					},
					"pass": schema.StringAttribute{
						Description: "The password that will be authenticated with",
						Optional:    true,
						Sensitive:   true,
					},
					"ssl_mode": schema.StringAttribute{
						Description: "The SSL mode for the postgres server",
						Optional:    true,
					},
					"tunnel": tunnelSchema,
				},
			},
			"mysql": schema.SingleNestedAttribute{
				Description: "The mysql database that will be associated with this connection",
				Optional:    true,
				// PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description: "Standard mysql url connection string.",
						Optional:    true,
					},

					"host": schema.StringAttribute{
						Description: "The host name of the mysql server",
						Optional:    true,
					},
					"port": schema.Int64Attribute{
						Description: "The port of the mysql server",
						Optional:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the database that will be connected to",
						Optional:    true,
					},
					"user": schema.StringAttribute{
						Description: "The name of the user that will be authenticated with",
						Optional:    true,
					},
					"pass": schema.StringAttribute{
						Description: "The password that will be authenticated with",
						Optional:    true,
						Sensitive:   true,
					},
					"protocol": schema.StringAttribute{
						Description: "The protocol of the mysql server",
						Optional:    true,
					},
					"tunnel": tunnelSchema,
				},
			},
			"aws_s3": schema.SingleNestedAttribute{
				Description: "The aws s3 bucket that will be associated with this connection",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"bucket": schema.StringAttribute{
						Description: "The name of the S3 bucket",
						Required:    true,
					},
					"path_prefix": schema.StringAttribute{
						Description: "The folder within the bucket that the connection will be scoped to",
						Optional:    true,
					},
					"region": schema.StringAttribute{
						Description: "The region that will be used by the SDK to access the bucket",
						Optional:    true,
					},
					"endpoint": schema.StringAttribute{
						Description: "The endpoint that will be used by the SDK to access the bucket",
						Optional:    true,
					},
					"credentials": schema.SingleNestedAttribute{
						Description: "Credentials that may be necessary to access the S3 bucket in a R/W fashion",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"profile": schema.StringAttribute{
								Description: "The profile found in the ~/.aws/config that can be used to access credentials",
								Optional:    true,
							},
							"access_key_id": schema.StringAttribute{
								Description: "The AWS access key id",
								Optional:    true,
							},
							"secret_access_key": schema.StringAttribute{
								Description: "The AWS secret access key",
								Optional:    true,
								Sensitive:   true,
							},
							"session_token": schema.StringAttribute{
								Description: "The AWS session token",
								Optional:    true,
							},
							"from_ec2_role": schema.BoolAttribute{
								Description: "Will result in the sync operations pulling from the EC2 role",
								Optional:    true,
							},
							"role_arn": schema.StringAttribute{
								Description: "The role arn that can be assumed",
								Optional:    true,
							},
							"role_external_id": schema.StringAttribute{
								Description: "The external id that will be provided when the role arn is assumed",
								Optional:    true,
							},
						},
					},
				},
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the connection",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

func (r *ConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ConfigData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ConfigData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = providerData.ConnectionClient
	r.accountId = providerData.AccountId
}

func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data connection_model.ConnectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId, err := r.getAccountId(&data)
	if err != nil {
		resp.Diagnostics.AddError("no account id", err.Error())
		return
	}
	data.AccountId = types.StringValue(accountId)

	createRequest, err := data.ToCreateConnectionDto()
	if err != nil {
		resp.Diagnostics.AddError("unable to create connection request", err.Error())
		return
	}

	createResp, err := r.client.CreateConnection(ctx, connect.NewRequest(createRequest))
	if err != nil {
		resp.Diagnostics.AddError("create connection error", err.Error())
		return
	}

	connection := createResp.Msg.GetConnection()
	tflog.Trace(ctx, "created connection")

	newModel := connection_model.ConnectionResourceModel{}
	err = newModel.FromDto(connection)
	if err != nil {
		resp.Diagnostics.AddError("connection model hydration error", err.Error())
		return
	}
	tflog.Trace(ctx, "mapped connection to model during creation")

	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data connection_model.ConnectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connResp, err := r.client.GetConnection(ctx, connect.NewRequest(&mgmtv1alpha1.GetConnectionRequest{
		Id: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get connection", err.Error())
		return
	}

	tflog.Trace(ctx, "retrieved connection")

	connection := connResp.Msg.GetConnection()

	newModel := connection_model.ConnectionResourceModel{}
	err = newModel.FromDto(connection)
	if err != nil {
		resp.Diagnostics.AddError("connection model hydration error", err.Error())
		return
	}
	tflog.Trace(ctx, "mapped connection to model during read")

	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data connection_model.ConnectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "read in planned model during update")

	updateRequest, err := data.ToUpdateConnectionDto()
	if err != nil {
		resp.Diagnostics.AddError("unable to map connection model to update request", err.Error())
		return
	}

	updateResp, err := r.client.UpdateConnection(ctx, connect.NewRequest(updateRequest))
	if err != nil {
		resp.Diagnostics.AddError("unable to update connection", err.Error())
		return
	}
	tflog.Trace(ctx, "updated connection")

	updatedConnection := updateResp.Msg.GetConnection()

	newModel := connection_model.ConnectionResourceModel{}
	err = newModel.FromDto(updatedConnection)
	if err != nil {
		resp.Diagnostics.AddError("connection model hydration error", err.Error())
		return
	}
	tflog.Trace(ctx, "mapped connection to model during update")

	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data connection_model.ConnectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteConnection(ctx, connect.NewRequest(&mgmtv1alpha1.DeleteConnectionRequest{
		Id: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete connection", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted connection")
}

func (r *ConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Unable to import", "must provide ID")
		return
	}

	connResp, err := r.client.GetConnection(ctx, connect.NewRequest(&mgmtv1alpha1.GetConnectionRequest{
		Id: req.ID,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get connection", err.Error())
		return
	}
	tflog.Trace(ctx, "retrieved connection during import")

	connection := connResp.Msg.GetConnection()

	var data connection_model.ConnectionResourceModel
	err = data.FromDto(connection)
	if err != nil {
		resp.Diagnostics.AddError("connection config hydration error", err.Error())
		return
	}

	tflog.Trace(ctx, "mapped connection to model during import")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectionResource) getAccountId(data *connection_model.ConnectionResourceModel) (string, error) {
	var accountId string
	if data.AccountId.ValueString() == "" {
		if r.accountId != nil {
			accountId = *r.accountId
		}
	} else {
		accountId = data.AccountId.ValueString()
	}
	if accountId == "" {
		return "", errors.New("must provide account id either on the resource or provide through environment configuration")
	}
	return accountId, nil
}
