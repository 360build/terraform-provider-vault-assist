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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &BootstrapSecretResource{}
	_ resource.ResourceWithConfigure = &BootstrapSecretResource{}
)

// BootstrapSecretResource is the resource implementation.
type BootstrapSecretResource struct {
	client *vaultclient.VaultClient
}

// NewBootstrapSecretResource returns a new instance of BootstrapSecretResource.
func NewBootstrapSecretResource() resource.Resource {
	return &BootstrapSecretResource{}
}

// Metadata returns the resource type name.
func (r *BootstrapSecretResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vaultassist_bootstrap_secret"
}

// Schema defines the schema for the resource.
func (r *BootstrapSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created": schema.BoolAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Required: true,
			},
			"mount": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *BootstrapSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
type BootstrapSecretResourceModel struct {
	Created types.Bool   `tfsdk:"created"`
	Path    types.String `tfsdk:"path"`
	Mount   types.String `tfsdk:"mount"`
}

func (r *BootstrapSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from the plan into the state model
	var plan BootstrapSecretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construct the path for the secret
	path := fmt.Sprintf("%s/data/%s", plan.Mount.ValueString(), plan.Path.ValueString())

	// Check if the secret already exists
	secret, err := r.client.ReadSecret(path)
	fmt.Print(err)
	fmt.Println(secret)

	emptySecret := map[string]interface{}{
		"options": map[string]interface{}{
			"cas": 0, // Ensure write only if secret does not exist
		},
		"data": map[string]interface{}{},
	}
	r.client.WriteSecret(path, emptySecret)

	// Set the new state using the struct model
	state := BootstrapSecretResourceModel{
		Created: types.BoolValue(true),
		Path:    plan.Path,
		Mount:   plan.Mount,
	}
	resp.State.Set(ctx, &state)
}

// Read reads the resource state.
func (r *BootstrapSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// No-op
}

// Update updates the resource state.
func (r *BootstrapSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

// Delete deletes the resource state.
func (r *BootstrapSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
