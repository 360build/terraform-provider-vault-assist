package provider

import (
	"context"
	// "crypto/sha256"
	// "encoding/json"
	"fmt"
	// "os"

	"github.com/360-build/terraform-provider-vaultassist/internal/provider/vaultclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &PatchSecretResource{}
	_ resource.ResourceWithConfigure = &PatchSecretResource{}
)

// PatchSecretResource is the resource implementation.
type PatchSecretResource struct {
	client *vaultclient.VaultClient
}

// NewPatchSecretResource returns a new instance of PatchSecretResource.
func NewPatchSecretResource() resource.Resource {
	return &PatchSecretResource{}
}

// Metadata returns the resource type name.
func (r *PatchSecretResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vaultassist_patch_secret"
}

// Schema defines the schema for the resource.
func (r *PatchSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created": schema.BoolAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mount": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *PatchSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*vaultclient.VaultClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *vaultclient.VaultClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Define the resource state model
type PatchSecretResourceModel struct {
	Created types.Bool   `tfsdk:"created"`
	Path    types.String `tfsdk:"path"`
	Mount   types.String `tfsdk:"mount"`
	Key     types.String `tfsdk:"key"`
	Value   types.String `tfsdk:"value"`
}

func (r *PatchSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from the plan into the state model
	var plan PatchSecretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construct the path for the secret
	path := fmt.Sprintf("%s/data/%s", plan.Mount.ValueString(), plan.Path.ValueString())

	patchSecret := map[string]interface{}{
		"options": map[string]interface{}{},
		"data": map[string]interface{}{
			plan.Key.ValueString(): plan.Value.ValueString(),
		},
	}
	r.client.PatchSecret(path, patchSecret)

	// Set the new state using the struct model
	state := PatchSecretResourceModel{
		Created: types.BoolValue(true),
		Path:    plan.Path,
		Mount:   plan.Mount,
		Key:     plan.Key,
		Value:   plan.Value,
	}
	resp.State.Set(ctx, &state)
}

// Read reads the resource state.
func (r *PatchSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// No-op
}

// Update updates the resource state.
func (r *PatchSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the current state
	var state PatchSecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the planned changes
	var plan PatchSecretResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the value has changed
	if state.Value.ValueString() == plan.Value.ValueString() {
		// No changes to the value, no update needed
		return
	}

	// Construct the path for the secret
	path := fmt.Sprintf("%s/data/%s", plan.Mount.ValueString(), plan.Path.ValueString())

	// Prepare the patch payload
	patchSecret := map[string]interface{}{
		"options": map[string]interface{}{},
		"data": map[string]interface{}{
			plan.Key.ValueString(): plan.Value.ValueString(),
		},
	}

	// Attempt to update the secret
	if err := r.client.PatchSecret(path, patchSecret); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Secret",
			fmt.Sprintf("Could not update secret at path '%s', unexpected error: %s", path, err.Error()),
		)
		return
	}

	// Update the state with the new value
	state.Value = plan.Value
	resp.State.Set(ctx, &state)
}

func (r *PatchSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from the state into the state model
	var state PatchSecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure the key is not empty
	if state.Key.IsNull() || state.Key.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid Key",
			"The key to delete is empty or null.",
		)
		return
	}

	// Construct the path for the secret
	path := fmt.Sprintf("%s/data/%s", state.Mount.ValueString(), state.Path.ValueString())

	// Prepare the patch payload
	patchSecret := map[string]interface{}{
		"options": map[string]interface{}{},
		"data": map[string]interface{}{
			state.Key.ValueString(): nil,
		},
	}

	// Attempt to delete the secret
	if err := r.client.PatchSecret(path, patchSecret); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Secret",
			fmt.Sprintf("Could not delete secret at path '%s', unexpected error: %s", path, err.Error()),
		)
		return
	}

	// Remove the resource from the state
	resp.State.RemoveResource(ctx)
}
