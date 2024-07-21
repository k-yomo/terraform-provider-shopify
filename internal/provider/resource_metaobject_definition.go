package provider

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/shopify"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MetaobjectDefinitionResource{}
var _ resource.ResourceWithImportState = &MetaobjectDefinitionResource{}

func NewMetaobjectDefinitionResource() resource.Resource {
	return &MetaobjectDefinitionResource{}
}

// MetaobjectDefinitionResource defines the resource implementation.
type MetaobjectDefinitionResource struct {
	client *shopify.Client
}

type MetaobjectFieldDefinitionResourceModel struct {
}

// MetaobjectDefinitionResourceModel describes the resource data model.
type MetaobjectDefinitionResourceModel struct {
	ID                types.String                      `tfsdk:"id"`
	Name              types.String                      `tfsdk:"name"`
	Type              types.String                      `tfsdk:"type"`
	Description       types.String                      `tfsdk:"description"`
	DisplayNameKey    types.String                      `tfsdk:"display_name_key"`
	FieldDefinitions  []*MetaobjectFieldDefinitionModel `tfsdk:"field_definitions"`
	HasThumbnailField types.Bool                        `tfsdk:"has_thumbnail_field"`
}

// MetaobjectFieldDefinitionModel describes the metaobject field definition data model.
type MetaobjectFieldDefinitionModel struct {
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Required    types.Bool   `tfsdk:"required"`
}

func (r *MetaobjectDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metaobject_definition"
}

func (r *MetaobjectDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides the definition of a generic object structure composed of metafields.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique ID of the metaobject.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name for the metaobject definition.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: `The type of the object definition. Defines the namespace of associated metafields.`,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description for the metaobject definition.",
				Optional:            true,
			},
			"display_name_key": schema.StringAttribute{
				MarkdownDescription: "The key of a field to reference as the display name for each object.",
				Optional:            true,
				Computed:            true,
			},
			"field_definitions": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: `The key of the new field definition. This can't be changed.
Must be 3-64 characters long and only contain alphanumeric, hyphen, and underscore characters.
`,
							Required: true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "A human-readable name for the field. This can be changed at any time.",
							Optional:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "An administrative description of the field.",
							Optional:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The metafield type applied to values of the field. If the type is changed, the field will be recreated.",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								utils.LogAttributeChangeModifier(func(ctx context.Context, req planmodifier.StringRequest) diag.Diagnostics {
									return diag.Diagnostics{diag.NewWarningDiagnostic(
										"Changing the type will recreate the field.",
										"Changing the type of the field definition will recreate the field. It will delete the existing data associated with the field.",
									)}
								},
									"Changing the type will recreate the field.",
									"Changing the type will recreate the field.",
								),
							},
						},
						"required": schema.BoolAttribute{
							MarkdownDescription: "Whether metaobjects require a saved value for the field.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
				Required: true,
			},
			"has_thumbnail_field": schema.BoolAttribute{
				MarkdownDescription: "Whether this metaobject definition has field whose type can visually represent a metaobject with the thumbnailField.",
				Computed:            true,
			},
		},
	}
}

func (r *MetaobjectDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.client, _ = req.ProviderData.(*shopify.Client)
}

func (r *MetaobjectDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MetaobjectDefinitionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fieldDefinitions := make([]*shopify.MetaobjectFieldDefinition, 0, len(data.FieldDefinitions))
	for _, fieldDefinitionModel := range data.FieldDefinitions {
		fieldDefinitions = append(fieldDefinitions, &shopify.MetaobjectFieldDefinition{
			Key:         fieldDefinitionModel.Key.ValueString(),
			Name:        fieldDefinitionModel.Name.ValueString(),
			Description: fieldDefinitionModel.Description.ValueStringPointer(),
			Type: &shopify.MetafieldDefinitionType{
				Name: fieldDefinitionModel.Type.ValueString(),
			},
			Required: fieldDefinitionModel.Required.ValueBool(),
		})
	}

	var shopifyFieldDefinitions []*shopify.MetaobjectFieldDefinitionCreateInput
	for _, fieldDefinitionModel := range data.FieldDefinitions {
		shopifyFieldDefinitions = append(shopifyFieldDefinitions, convertMetaobjectFieldDefinitionModelToCreateInput(fieldDefinitionModel))
	}

	var displayNameKey *string
	if data.DisplayNameKey.ValueString() != "" {
		displayNameKey = data.DisplayNameKey.ValueStringPointer()
	}
	input := shopify.MetaobjectDefinitionCreateInput{
		Type:             data.Type.ValueString(),
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueStringPointer(),
		DisplayNameKey:   displayNameKey,
		FieldDefinitions: shopifyFieldDefinitions,
	}
	createdMetaobjectDefinition, err := r.client.CreateMetaobjectDefinition(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create metaobject definition, got error: %s", err))
		return
	}

	createdData := convertMetaobjectDefinitionToResourceModel(createdMetaobjectDefinition, &data)
	tflog.Trace(ctx, "created a metaobject definition", map[string]interface{}{
		"id": createdData.ID,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, createdData)...)
}

func (r *MetaobjectDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MetaobjectDefinitionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metaobjectDefinition, err := r.client.GetMetaobjectDefinition(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read metaobject definition, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, convertMetaobjectDefinitionToResourceModel(metaobjectDefinition, &data))...)
}

func (r *MetaobjectDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MetaobjectDefinitionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var oldFieldDefinitions []*MetaobjectFieldDefinitionModel
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("field_definitions"), &oldFieldDefinitions)...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldFieldDefinitionMap := make(map[string]*MetaobjectFieldDefinitionModel, len(oldFieldDefinitions))
	for _, fieldDefinition := range oldFieldDefinitions {
		oldFieldDefinitionMap[fieldDefinition.Key.ValueString()] = fieldDefinition
	}

	var fieldDefinitions1stReq []*shopify.MetaobjectFieldDefinitionOperationInput
	var fieldDefinitions2ndReq []*shopify.MetaobjectFieldDefinitionOperationInput
	var recreateFieldDefinitions []string
	for _, newFieldDef := range data.FieldDefinitions {
		if oldFieldDef, ok := oldFieldDefinitionMap[newFieldDef.Key.ValueString()]; ok {
			delete(oldFieldDefinitionMap, newFieldDef.Key.ValueString())
			if reflect.DeepEqual(oldFieldDef, newFieldDef) {
				continue
			}
			if newFieldDef.Type != oldFieldDef.Type {
				fieldDefinitions1stReq = append(fieldDefinitions1stReq, &shopify.MetaobjectFieldDefinitionOperationInput{
					Delete: &shopify.MetaobjectFieldDefinitionDeleteInput{
						Key: oldFieldDef.Key.ValueString(),
					},
				})
				fieldDefinitions2ndReq = append(fieldDefinitions2ndReq, &shopify.MetaobjectFieldDefinitionOperationInput{
					Create: convertMetaobjectFieldDefinitionModelToCreateInput(newFieldDef),
				})
				recreateFieldDefinitions = append(recreateFieldDefinitions, newFieldDef.Key.ValueString())
			} else {
				fieldDefinitions1stReq = append(fieldDefinitions1stReq, &shopify.MetaobjectFieldDefinitionOperationInput{
					Update: &shopify.MetaobjectFieldDefinitionUpdateInput{
						Key:         newFieldDef.Key.ValueString(),
						Name:        newFieldDef.Name.ValueStringPointer(),
						Description: newFieldDef.Description.ValueStringPointer(),
						Required:    newFieldDef.Required.ValueBool(),
					},
				})
			}
		} else {
			fieldDefinitions1stReq = append(fieldDefinitions1stReq, &shopify.MetaobjectFieldDefinitionOperationInput{
				Create: convertMetaobjectFieldDefinitionModelToCreateInput(newFieldDef),
			})
		}
	}
	if len(recreateFieldDefinitions) > 0 {
		tflog.Warn(ctx, "")
	}

	for _, oldFieldDef := range oldFieldDefinitionMap {
		fieldDefinitions1stReq = append(fieldDefinitions1stReq, &shopify.MetaobjectFieldDefinitionOperationInput{
			Delete: &shopify.MetaobjectFieldDefinitionDeleteInput{
				Key: oldFieldDef.Key.ValueString(),
			},
		})
	}

	var displayNameKey *string
	if data.DisplayNameKey.ValueString() != "" {
		displayNameKey = data.DisplayNameKey.ValueStringPointer()
	}
	input1stReq := shopify.MetaobjectDefinitionUpdateInput{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueStringPointer(),
		DisplayNameKey:   displayNameKey,
		FieldDefinitions: fieldDefinitions1stReq,
	}
	updatedMetaobjectDefinition, err := r.client.UpdateMetaobjectDefinition(ctx, data.ID.ValueString(), &input1stReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update metaobject definition, got error: %s", err))
		return
	}
	updateData := convertMetaobjectDefinitionToResourceModel(updatedMetaobjectDefinition, &data)

	if len(fieldDefinitions2ndReq) > 0 {
		input2ndReq := shopify.MetaobjectDefinitionUpdateInput{
			Name:             data.Name.ValueString(),
			Description:      data.Description.ValueStringPointer(),
			DisplayNameKey:   displayNameKey,
			FieldDefinitions: fieldDefinitions2ndReq,
		}
		updatedMetaobjectDefinition, err := r.client.UpdateMetaobjectDefinition(ctx, data.ID.ValueString(), &input2ndReq)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update metaobject definition, got error: %s", err))
			return
		}
		updateData = convertMetaobjectDefinitionToResourceModel(updatedMetaobjectDefinition, &data)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updateData)...)
}

func (r *MetaobjectDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MetaobjectDefinitionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMetaobjectDefinition(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete metaobject definition, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a metaobject definition", map[string]interface{}{
		"id": data.ID,
	})
}

func (r *MetaobjectDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertMetaobjectDefinitionToResourceModel(definition *shopify.MetaobjectDefinition, data *MetaobjectDefinitionResourceModel) *MetaobjectDefinitionResourceModel {
	fieldDefinitionModels := make([]*MetaobjectFieldDefinitionModel, 0, len(definition.FieldDefinitions))
	for _, fieldDefinition := range definition.FieldDefinitions {
		fieldDefinitionModels = append(fieldDefinitionModels, convertMetaobjectFieldDefinitionToModel(fieldDefinition))
	}

	// Sort field definitions by order in the original data not to produce unnecessary diffs
	fieldDefinitionOrderMap := make(map[string]int, len(data.FieldDefinitions))
	for i, fieldDefinition := range data.FieldDefinitions {
		fieldDefinitionOrderMap[fieldDefinition.Key.ValueString()] = i
	}
	sort.Slice(fieldDefinitionModels, func(i, j int) bool {
		return fieldDefinitionOrderMap[fieldDefinitionModels[i].Key.ValueString()] < fieldDefinitionOrderMap[fieldDefinitionModels[j].Key.ValueString()]
	})

	return &MetaobjectDefinitionResourceModel{
		ID:                types.StringValue(definition.ID),
		Name:              types.StringValue(definition.Name),
		Type:              types.StringValue(definition.Type),
		Description:       types.StringPointerValue(definition.Description),
		DisplayNameKey:    types.StringPointerValue(definition.DisplayNameKey),
		FieldDefinitions:  fieldDefinitionModels,
		HasThumbnailField: types.BoolValue(definition.HasThumbnailField),
	}
}

func convertMetaobjectFieldDefinitionToModel(definition *shopify.MetaobjectFieldDefinition) *MetaobjectFieldDefinitionModel {
	return &MetaobjectFieldDefinitionModel{
		Key:         types.StringValue(definition.Key),
		Name:        types.StringValue(definition.Name),
		Description: types.StringPointerValue(definition.Description),
		Type:        types.StringValue(definition.Type.Name),
		Required:    types.BoolValue(definition.Required),
	}
}

func convertMetaobjectFieldDefinitionModelToCreateInput(model *MetaobjectFieldDefinitionModel) *shopify.MetaobjectFieldDefinitionCreateInput {
	return &shopify.MetaobjectFieldDefinitionCreateInput{
		Key:         model.Key.ValueString(),
		Name:        model.Name.ValueStringPointer(),
		Description: model.Description.ValueStringPointer(),
		Type:        model.Type.ValueString(),
		Required:    model.Required.ValueBool(),
	}
}
