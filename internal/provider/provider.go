package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	"github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1/mgmtv1alpha1connect"
	http_client "github.com/nucleuscloud/terraform-provider-neosync/internal/http/client"
)

const (
	endpointEnvVarKey  = "NEOSYNC_ENDPOINT"
	apiTokenEnvVarKey  = "NEOSYNC_API_TOKEN"
	accountIdEnvVarKey = "NEOSYNC_ACCOUNT_ID"
)

// Ensure NeosyncProvider satisfies various provider inferfaces.
var _ provider.Provider = &NeosyncProvider{}

type NeosyncProvider struct {
	version         string
	defaultEndpoint string
}

type NeosyncProviderModel struct {
	ApiToken  types.String `tfsdk:"api_token"`
	Endpoint  types.String `tfsdk:"endpoint"`
	AccountId types.String `tfsdk:"account_id"`
}

func (p *NeosyncProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "neosync"
	resp.Version = p.version
}

func (p *NeosyncProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The URL to the backend Neosync API server",
				Optional:    true,
			},
			"api_token": schema.StringAttribute{
				Description: "The account-level API token that will be used to authenticate with the API server",
				Optional:    true,
			},
			"account_id": schema.StringAttribute{
				Description: "The account id that should be associated with this provider and any resources that utilize it",
				Optional:    true,
			},
		},
	}
}

type ConfigData struct {
	AccountId        *string
	ConnectionClient mgmtv1alpha1connect.ConnectionServiceClient
	JobClient        mgmtv1alpha1connect.JobServiceClient
}

func (p *NeosyncProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiToken := os.Getenv(apiTokenEnvVarKey)
	endpoint := os.Getenv(endpointEnvVarKey)
	accountId := os.Getenv(accountIdEnvVarKey)
	// todo: add support for specifying account name along with a path to the location of a user jwt file

	var data NeosyncProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.ApiToken.ValueString() != "" {
		apiToken = data.ApiToken.ValueString()
	}

	if data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	}
	if endpoint == "" {
		endpoint = p.defaultEndpoint
	}

	if data.AccountId.ValueString() != "" {
		accountId = data.AccountId.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddWarning(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				fmt.Sprintf("the %s environment variable or provider ", apiTokenEnvVarKey)+
				"configuration block api_token attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Missing Endpoint Configuration",
			"While configuring the provider, the endpoint was not found in "+
				fmt.Sprintf("the %s environment variable or provider ", endpointEnvVarKey)+
				"configuration block endpoint attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpclient := http.DefaultClient
	if apiToken != "" {
		httpclient = http_client.NewWithHeaders(
			map[string]string{"Authorization": fmt.Sprintf("Bearer %s", apiToken)},
		)
	}

	connclient := mgmtv1alpha1connect.NewConnectionServiceClient(
		httpclient,
		endpoint,
	)
	if apiToken != "" && accountId == "" {
		userclient := mgmtv1alpha1connect.NewUserAccountServiceClient(httpclient, endpoint)
		userAccountsResp, err := userclient.GetUserAccounts(ctx, connect.NewRequest(&mgmtv1alpha1.GetUserAccountsRequest{}))
		if err != nil {
			resp.Diagnostics.AddError("user account error", err.Error())
			return
		}
		accounts := userAccountsResp.Msg.Accounts
		if len(accounts) == 0 {
			resp.Diagnostics.AddError("user account error", "unable to find any accounts associated with provided api token")
			return
		}
		accountId = accounts[0].Id
	}

	configData := &ConfigData{
		ConnectionClient: connclient,
	}
	if accountId != "" {
		configData.AccountId = &accountId
	}

	resp.DataSourceData = configData
	resp.ResourceData = configData
}

func (p *NeosyncProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewConnectionResource,
	}
}

func (p *NeosyncProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConnectionDataSource,
		NewJobDataSource,
	}
}

func New(version string, defaultEndpoint string) func() provider.Provider {
	return func() provider.Provider {
		return &NeosyncProvider{
			version:         version,
			defaultEndpoint: defaultEndpoint,
		}
	}
}
