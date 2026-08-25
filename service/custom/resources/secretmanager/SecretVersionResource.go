package secretmanager

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.mws.cloud/go-sdk/mws/wait"
	ctxvalues "go.mws.cloud/go-sdk/pkg/context/values"
	secretmanagerref "go.mws.cloud/go-sdk/service/resources/references/secretmanager"
	"go.mws.cloud/go-sdk/service/secretmanager/client"
	apimodel "go.mws.cloud/go-sdk/service/secretmanager/model"
	resourcesdk "go.mws.cloud/go-sdk/service/secretmanager/sdk"

	"go.mws.cloud/terraform-provider-mws/internal/cmp"
	"go.mws.cloud/terraform-provider-mws/internal/diag"
	provider "go.mws.cloud/terraform-provider-mws/internal/provider/public"
	conv "go.mws.cloud/terraform-provider-mws/service/resources/secretmanager/converter"
	tfmodel "go.mws.cloud/terraform-provider-mws/service/resources/secretmanager/model"
)

var (
	_ resource.Resource                = &SecretVersionResource{}
	_ resource.ResourceWithImportState = &SecretVersionResource{}
)

type SecretVersionResource struct {
	sdk    *resourcesdk.SecretVersion
	config *provider.Config
}

func NewSecretVersionResource() resource.Resource {
	return &SecretVersionResource{}
}

func (m *SecretVersionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Info(ctx, "SecretVersionResource.Metadata")
	resp.TypeName = req.ProviderTypeName + "_secretmanager_secret_version"
}

func (m *SecretVersionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "SecretVersionResource.Schema")
	resp.Schema = new(tfmodel.SecretVersion).GetSchema()
	resp.Schema.Attributes["name"] = schema.StringAttribute{
		MarkdownDescription: `Имя секрета.`,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplaceIfConfigured(),
		},
	}
	resp.Schema.Attributes["project"] = schema.StringAttribute{
		MarkdownDescription: `Путь к проекту.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			stringplanmodifier.RequiresReplaceIfConfigured(),
		},
	}
	resp.Schema.Attributes["version"] = schema.StringAttribute{
		MarkdownDescription: `Версия секрета.`,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	resp.Schema.Attributes["timeouts"] = timeouts.Attributes(ctx, timeouts.Opts{
		Create: true,
		Update: true,
		Delete: true,
	})
	resp.Schema.Attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func (m *SecretVersionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	var err error
	tflog.Info(ctx, "SecretVersionResource.Configure")
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*provider.Data)
	if !ok || providerData == nil {
		resp.Diagnostics.AddError("Internal Provider Error", "Unexpected type of req.ProviderData")
		return
	}

	m.config = providerData.Config

	m.sdk, err = resourcesdk.NewSecretVersion(ctx, providerData.SDK)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create SDK client for SecretVersion",
			diag.FormatError(err),
		)
		return
	}
}

func (m *SecretVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "SecretVersionResource.Create")

	var data tfmodel.SecretVersionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectParam := cmp.Or(data.ProjectParam, m.config.Project)
	if projectParam.IsNull() || projectParam.IsUnknown() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			`Parameter "project" is null or unknown`,
		)
		return
	}
	data.ProjectParam = projectParam
	ctx = ctxvalues.With(ctx, "project", projectParam.String())

	resourceWaiterTimeout, diags := data.Timeouts.Create(ctx, 3600*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "SecretVersionResource.Timeouts")
		return
	}

	var config tfmodel.SecretVersionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "SecretVersionResource.Config")
		return
	}
	data.Data = config.Data

	body, diags := conv.SecretVersionTFToAPIRequestModel(ctx, &data.SecretVersion)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "SecretVersionResource.TFToAPI")
		return
	}

	apiRes, err := m.sdk.AddSecretVersion(
		ctx,
		client.AddSecretVersionRequest{
			Project: data.ProjectParam.ValueString(),
			Name:    data.NameParam.ValueString(),
			Body: apimodel.AddSecretVersionRequest{
				Metadata: body.Metadata,
				Spec:     &body.Spec,
			},
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Add SecretVersion",
			diag.FormatError(err),
		)
		return
	}

	version := string(apiRes.Metadata.Id.ResourceName())

	_, err = m.sdk.GetSecretVersion(
		ctx,
		client.GetSecretVersionRequest{
			Project: data.ProjectParam.ValueString(),
			Name:    data.NameParam.ValueString(),
			Version: version,
		},
		client.WithWait(wait.WithTimeout(resourceWaiterTimeout)),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Get Enabled Service",
			diag.FormatError(err),
		)
		return
	}

	data.ID = types.StringValue(apiRes.Metadata.Id.ID())
	data.VersionParam = types.StringValue(version)

	tfRes, diags := conv.SecretVersionAPIOptionalResponseToTFModel(ctx, apiRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || tfRes == nil {
		tflog.Debug(ctx, "SecretVersionResource.ApiToTfConvert.tfRes", map[string]any{"res": tfRes})
		return
	}

	data.SecretVersion = *tfRes

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *SecretVersionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "SecretVersionResource.Read")

	var data tfmodel.SecretVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectParam := cmp.Or(data.ProjectParam, m.config.Project)
	if projectParam.IsNull() || projectParam.IsUnknown() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			`Parameter "project" is null or unknown`,
		)
		return
	}
	data.ProjectParam = projectParam
	ctx = ctxvalues.With(ctx, "project", projectParam.String())

	apiRes, err := m.sdk.GetSecretVersion(
		ctx,
		client.GetSecretVersionRequest{
			Project: data.ProjectParam.ValueString(),
			Name:    data.NameParam.ValueString(),
			Version: data.VersionParam.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Get SecretVersion",
			diag.FormatError(err),
		)
		return
	}

	data.ID = types.StringValue(apiRes.Metadata.Id.ID())

	tfRes, diags := conv.SecretVersionAPIOptionalResponseToTFModel(ctx, apiRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || tfRes == nil {
		tflog.Debug(ctx, "SecretVersionResource.ApiToTfConvert.tfRes", map[string]any{"res": tfRes})
		return
	}

	data.SecretVersion = *tfRes

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *SecretVersionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "SecretVersionResource.Update")

	var data tfmodel.SecretVersionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectParam := cmp.Or(data.ProjectParam, m.config.Project)
	if projectParam.IsNull() || projectParam.IsUnknown() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			`Parameter "project" is null or unknown`,
		)
		return
	}
	data.ProjectParam = projectParam
	ctx = ctxvalues.With(ctx, "project", projectParam.String())

	resourceWaiterTimeout, diags := data.Timeouts.Update(ctx, 3600*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "SecretVersionResource.Timeouts")
		return
	}

	body, diags := conv.SecretVersionTFToAPIRequestModel(ctx, &data.SecretVersion)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "SecretVersionResource.TFToAPI")
		return
	}

	apiRes, err := m.sdk.UpdateSecretVersion(
		ctx,
		client.UpdateSecretVersionRequest{
			Project: data.ProjectParam.ValueString(),
			Name:    data.NameParam.ValueString(),
			Version: data.VersionParam.ValueString(),
			Body:    body.AsUpdateModel(),
		},
		client.WithWait(wait.WithTimeout(resourceWaiterTimeout)),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update SecretVersion",
			diag.FormatError(err),
		)
		return
	}

	data.ID = types.StringValue(apiRes.Metadata.Id.ID())

	tfRes, diags := conv.SecretVersionAPIOptionalResponseToTFModel(ctx, apiRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || tfRes == nil {
		tflog.Debug(ctx, "SecretVersionResource.ApiToTfConvert.tfRes", map[string]any{"res": tfRes})
		return
	}

	data.SecretVersion = *tfRes

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *SecretVersionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "SecretVersionResource.Delete")

	var data tfmodel.SecretVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectParam := cmp.Or(data.ProjectParam, m.config.Project)
	if projectParam.IsNull() || projectParam.IsUnknown() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			`Parameter "project" is null or unknown`,
		)
		return
	}
	data.ProjectParam = projectParam
	ctx = ctxvalues.With(ctx, "project", projectParam.String())

	resourceWaiterTimeout, diags := data.Timeouts.Delete(ctx, 3600*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "SecretVersionResource.Timeouts")
		return
	}

	err := m.sdk.DeleteSecretVersion(
		ctx,
		client.DeleteSecretVersionRequest{
			Project: data.ProjectParam.ValueString(),
			Name:    data.NameParam.ValueString(),
			Version: data.VersionParam.ValueString(),
		},
		client.WithWait(wait.WithTimeout(resourceWaiterTimeout)),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete SecretVersion",
			diag.FormatError(err),
		)
		return
	}
}

func (m *SecretVersionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "SecretVersionResource.ImportState")

	var data tfmodel.SecretVersionModel

	ref, err := secretmanagerref.ParseSecretVersionRef(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Parse reference for SecretVersion",
			diag.FormatError(err),
		)
		return
	}

	apiRes, err := m.sdk.GetSecretVersion(
		ctx,
		client.GetSecretVersionRequest{
			Project: ref.GetProject(),
			Name:    ref.GetSecretName(),
			Version: ref.GetVersion(),
		})
	if err != nil {
		resp.Diagnostics.AddError(
			"Get SecretVersion",
			diag.FormatError(err),
		)
		return
	}

	data.ID = types.StringValue(apiRes.Metadata.Id.ID())

	tfRes, diags := conv.SecretVersionAPIOptionalResponseToTFModel(ctx, apiRes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || tfRes == nil {
		tflog.Debug(ctx, "SecretVersionResource.ApiToTfConvert.tfRes", map[string]any{"res": tfRes})
		return
	}

	data.SecretVersion = *tfRes

	data.ProjectParam = types.StringValue(ref.GetProject())
	data.NameParam = types.StringValue(ref.GetSecretName())
	data.VersionParam = types.StringValue(ref.GetVersion())

	var rwTimeouts timeouts.Value
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, tfpath.Root("timeouts"), &rwTimeouts)...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "SecretVersionResource.timeouts.GetAttribute")
		return
	}
	data.Timeouts = rwTimeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
