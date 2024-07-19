package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/shopify"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MetafieldDefinitionResource{}
var _ resource.ResourceWithImportState = &MetafieldDefinitionResource{}

func NewMetafieldDefinitionResource() resource.Resource {
	return &MetafieldDefinitionResource{}
}

// MetafieldDefinitionResource defines the resource implementation.
type MetafieldDefinitionResource struct {
	client *shopify.Client
}

// MetafieldDefinitionResourceModel describes the resource data model.
type MetafieldDefinitionResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	OwnerType   types.String `tfsdk:"owner_type"`
	Namespace   types.String `tfsdk:"namespace"`
	Key         types.String `tfsdk:"key"`
	Type        types.String `tfsdk:"type"`
	Pin         types.Bool   `tfsdk:"pin"`
}

func (r *MetafieldDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metafield_definition"
}

func (r *MetafieldDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "MetafieldDefinition definition resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique ID of the metafield.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name for the metafield definition.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description for the metafield definition.",
				Optional:            true,
			},
			"owner_type": schema.StringAttribute{
				MarkdownDescription: `The resource type that the metafield definition is attached to.
Possible values are:
- API_PERMISSION
- ARTICLE
- BLOG
- CARTTRANSFORM
- COLLECTION
- COMPANY
- COMPANY_LOCATION
- CUSTOMER
- DELIVERY_CUSTOMIZATION
- DISCOUNT
- DRAFTORDER
- FULFILLMENT_CONSTRAINT_RULE
- LOCATION
- MARKET
- MEDIA_IMAGE
- ORDER
- ORDER_ROUTING_LOCATION_RULE
- PAGE
- PAYMENT_CUSTOMIZATION
- PRODUCT
- PRODUCTVARIANT
- SHOP
- VALIDATION
- PRODUCTIMAGE
`,
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: `The container for a group of metafields that the metafield is or will be associated with. Used in tandem with ` + "`key`" + ` to lookup a metafield on a resource, preventing conflicts with other metafields with the same ` + "`key.`" + `
					Must be 3-255 characters long and can contain alphanumeric, hyphen, and underscore characters.`,
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for a metafield within its namespace.\nMust be 3-64 characters long and can contain alphanumeric, hyphen, and underscore characters.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: `The type of data that each of the metafields that belong to the metafield definition will store. Refer to the list of [supported types](https://shopify.dev/docs/apps/build/custom-data/metafields/list-of-data-types).`,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"pin": schema.BoolAttribute{
				MarkdownDescription: "Whether to pin the metafield definition.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *MetafieldDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*shopify.Client)
}

func (r *MetafieldDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MetafieldDefinitionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := shopify.CreateMetafieldDefinitionInput{
		Key:         data.Key.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Namespace:   data.Namespace.ValueString(),
		OwnerType:   data.OwnerType.ValueString(),
		Pin:         data.Pin.ValueBool(),
		Type:        data.Type.ValueString(),
	}
	createdMetafieldDefinition, err := r.client.CreateMetafieldDefinition(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create metafield definition, got error: %s", err))
		return
	}

	createdData := convertMetafieldDefinitionToResourceModel(createdMetafieldDefinition)
	tflog.Trace(ctx, "created a metafield definition", map[string]interface{}{
		"id": createdData.ID,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, createdData)...)
}

func (r *MetafieldDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MetafieldDefinitionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metafieldDefinition, err := r.client.GetMetafieldDefinition(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metafield definition, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, convertMetafieldDefinitionToResourceModel(metafieldDefinition))...)
}

func (r *MetafieldDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MetafieldDefinitionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := shopify.UpdateMetafieldDefinitionInput{
		Key:         data.Key.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Namespace:   data.Namespace.ValueString(),
		OwnerType:   data.OwnerType.ValueString(),
		Pin:         data.Pin.ValueBool(),
	}
	updatedMetafieldDefinition, err := r.client.UpdateMetafieldDefinition(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update metafield definition, got error: %s", err))
		return
	}
	updateData := convertMetafieldDefinitionToResourceModel(updatedMetafieldDefinition)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updateData)...)
}

func (r *MetafieldDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MetafieldDefinitionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMetafieldDefinition(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete metafield definition, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a metafield definition", map[string]interface{}{
		"id": data.ID,
	})
}

func (r *MetafieldDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertMetafieldDefinitionToResourceModel(definition *shopify.MetafieldDefinition) *MetafieldDefinitionResourceModel {
	return &MetafieldDefinitionResourceModel{
		ID:          types.StringValue(definition.ID),
		Name:        types.StringValue(definition.Name),
		Description: types.StringPointerValue(definition.Description),
		OwnerType:   types.StringValue(definition.OwnerType),
		Namespace:   types.StringValue(definition.Namespace),
		Key:         types.StringValue(definition.Key),
		Type:        types.StringValue(definition.Type.Name),
		Pin:         types.BoolValue(definition.PinnedPosition != nil),
	}
}
