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
)

var _ resource.Resource = &UserDefinedTransformerResource{}
var _ resource.ResourceWithImportState = &UserDefinedTransformerResource{}

func NewUserDefinedTransformerResource() resource.Resource {
	return &UserDefinedTransformerResource{}
}

type UserDefinedTransformerResource struct {
	client    mgmtv1alpha1connect.TransformersServiceClient
	accountId *string
}

type UserDefinedTransformerResourceModel struct {
	Id          types.String       `tfsdk:"id"`
	Name        types.String       `tfsdk:"name"`
	Description types.String       `tfsdk:"description"`
	Datatype    types.String       `tfsdk:"datatype"`
	Source      types.String       `tfsdk:"source"`
	Config      *TransformerConfig `tfsdk:"config"`
	AccountId   types.String       `tfsdk:"account_id"`
}

func (r *UserDefinedTransformerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_defined_transformer"
}

func (r *UserDefinedTransformerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Neosync User Defined Transformer resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique friendly name of the transformer",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the transformer",
				Required:    true,
			},
			"datatype": schema.StringAttribute{
				Description: "The datatype of the transformer",
				Required:    true,
			},
			"source": schema.StringAttribute{
				Description: "The system transformer that this user defined transformer is based off of. This is the same value that is specified as the key in the config",
				Required:    true,
			},
			"config": transformerSchema,
			"account_id": schema.StringAttribute{
				Description:   "The unique identifier of the account. Can be pulled from the API Key if present, or must be specified if using a user access token",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the transformer",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

func (r *UserDefinedTransformerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = providerData.TransformerClient
	r.accountId = providerData.AccountId
}

func (r *UserDefinedTransformerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserDefinedTransformerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var accountId string
	if data.AccountId.ValueString() == "" {
		if r.accountId != nil {
			accountId = *r.accountId
		}
	} else {
		accountId = data.AccountId.ValueString()
	}
	if accountId == "" {
		resp.Diagnostics.AddError("no account id", "must provide account id either on the resource or provide through environment configuration")
		return
	}

	data.AccountId = types.StringValue(accountId)
	dto, err := toTransformerDto(&data)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert transformer model to dto", err.Error())
		return
	}
	transResp, err := r.client.CreateUserDefinedTransformer(ctx, connect.NewRequest(&mgmtv1alpha1.CreateUserDefinedTransformerRequest{
		AccountId:         dto.AccountId,
		Name:              dto.Name,
		Description:       dto.Description,
		Type:              dto.DataType,
		Source:            dto.Source,
		TransformerConfig: dto.Config,
	}))
	if err != nil {
		resp.Diagnostics.AddError("create transformer error", err.Error())
		return
	}

	transformer := transResp.Msg.Transformer
	updatedModel, err := fromTransformerDto(transformer)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert dto to transformer model", err.Error())
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created transformer resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
}

func (r *UserDefinedTransformerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserDefinedTransformerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connResp, err := r.client.GetUserDefinedTransformerById(ctx, connect.NewRequest(&mgmtv1alpha1.GetUserDefinedTransformerByIdRequest{
		TransformerId: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get transformer", err.Error())
		return
	}

	transformer := connResp.Msg.Transformer
	updatedModel, err := fromTransformerDto(transformer)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert dto to transformer model", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
}

func (r *UserDefinedTransformerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserDefinedTransformerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dto, err := toTransformerDto(&data)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert transformer model to dto", err.Error())
		return
	}

	connResp, err := r.client.UpdateUserDefinedTransformer(ctx, connect.NewRequest(&mgmtv1alpha1.UpdateUserDefinedTransformerRequest{
		TransformerId:     dto.Id,
		Name:              dto.Name,
		Description:       dto.Description,
		TransformerConfig: dto.Config,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to update transformer", err.Error())
		return
	}

	transformer := connResp.Msg.Transformer
	updatedModel, err := fromTransformerDto(transformer)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert dto to transformer model", err.Error())
		return
	}

	tflog.Trace(ctx, "updated transformer")
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
}

func (r *UserDefinedTransformerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserDefinedTransformerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteUserDefinedTransformer(ctx, connect.NewRequest(&mgmtv1alpha1.DeleteUserDefinedTransformerRequest{
		TransformerId: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete transformer", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted transformer")
}

func (r *UserDefinedTransformerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Unable to import", "must provide ID")
		return
	}

	transResp, err := r.client.GetUserDefinedTransformerById(ctx, connect.NewRequest(&mgmtv1alpha1.GetUserDefinedTransformerByIdRequest{
		TransformerId: req.ID,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get transformer", err.Error())
		return
	}

	transformer := transResp.Msg.Transformer

	updatedModel, err := fromTransformerDto(transformer)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert dto to transformer model", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
}

func fromTransformerDto(dto *mgmtv1alpha1.UserDefinedTransformer) (*UserDefinedTransformerResourceModel, error) {
	if dto == nil {
		return nil, errors.New("dto was nil")
	}
	configModel, err := toTransformerConfigFromDto(dto.Config)
	if err != nil {
		return nil, err
	}
	model := &UserDefinedTransformerResourceModel{
		Id:          types.StringValue(dto.Id),
		Name:        types.StringValue(dto.Name),
		AccountId:   types.StringValue(dto.AccountId),
		Description: types.StringValue(dto.Description),
		Datatype:    types.StringValue(dto.DataType),
		Source:      types.StringValue(dto.Source),
		Config:      configModel,
	}
	return model, nil
}

func toTransformerDto(model *UserDefinedTransformerResourceModel) (*mgmtv1alpha1.UserDefinedTransformer, error) {
	if model == nil {
		return nil, errors.New("model was nil")
	}
	configDto, err := fromModelTransformerConfig(model.Config)
	if err != nil {
		return nil, err
	}
	dto := &mgmtv1alpha1.UserDefinedTransformer{
		Id:          model.Id.ValueString(),
		Name:        model.Name.ValueString(),
		AccountId:   model.AccountId.ValueString(),
		Description: model.Description.ValueString(),
		DataType:    model.Datatype.ValueString(),
		Source:      model.Source.ValueString(),
		Config:      configDto,
	}
	return dto, nil
}
