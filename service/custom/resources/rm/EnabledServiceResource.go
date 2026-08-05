package rm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.mws.cloud/go-sdk/mws/wait"
	ctxvalues "go.mws.cloud/go-sdk/pkg/context/values"
	"go.mws.cloud/go-sdk/service/rm/client"
	resourcesdk "go.mws.cloud/go-sdk/service/rm/sdk"

	"go.mws.cloud/terraform-provider-mws/internal/cmp"
	"go.mws.cloud/terraform-provider-mws/internal/diag"
	provider "go.mws.cloud/terraform-provider-mws/internal/provider/public"
	tfmodel "go.mws.cloud/terraform-provider-mws/service/custom/resources/rm/model"
)

var (
	_ resource.ResourceWithConfigure = &EnabledServiceResource{}
)

type EnabledServiceResource struct {
	sdk    *resourcesdk.EnabledService
	config *provider.Config
}

func NewEnabledServiceResource() resource.Resource {
	return &EnabledServiceResource{}
}

func (r *EnabledServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Info(ctx, "EnabledServiceResource.Metadata")
	resp.TypeName = req.ProviderTypeName + "_resmanager_enabled_service"
}

func (r *EnabledServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Info(ctx, "EnabledServiceResource.Schema")
	resp.Schema = schema.Schema{
		MarkdownDescription: "Ресурс для подключения сервисов в проекте",
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				MarkdownDescription: `Путь к проекту`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"service": schema.StringAttribute{
				MarkdownDescription: `Идентификатор сервиса`,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

func (r *EnabledServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	var err error
	tflog.Info(ctx, "EnabledServiceResource.Configure")
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*provider.Data)
	if !ok || providerData == nil {
		resp.Diagnostics.AddError("Internal Provider Error", "Unexpected type of req.ProviderData")
		return
	}

	r.config = providerData.Config

	r.sdk, err = resourcesdk.NewEnabledService(ctx, providerData.SDK)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create SDK client for EnabledService",
			diag.FormatError(err),
		)
		return
	}
}

func (r *EnabledServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "EnabledServiceResource.Create")

	var data tfmodel.EnabledServiceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectParam := cmp.Or(data.ProjectParam, r.config.Project)
	if projectParam.IsNull() || projectParam.IsUnknown() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			`Parameter "project" is null or unknown`,
		)
		return
	}
	data.ProjectParam = projectParam
	ctx = ctxvalues.With(ctx, "project", projectParam.String())

	resourceWaiterTimeout, diags := data.Timeouts.Create(ctx, time.Hour)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "EnabledServiceResource.Timeouts")
		return
	}

	_, err := r.sdk.EnableService(ctx, client.EnableServiceRequest{
		Project: data.ProjectParam.ValueString(),
		Service: data.ServiceParam.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Enable Service",
			diag.FormatError(err),
		)
		return
	}

	_, err = r.sdk.GetEnabledService(
		ctx,
		client.GetEnabledServiceRequest{
			Project: data.ProjectParam.ValueString(),
			Service: data.ServiceParam.ValueString(),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnabledServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "EnabledServiceResource.Read")

	var data tfmodel.EnabledServiceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectParam := cmp.Or(data.ProjectParam, r.config.Project)
	if projectParam.IsNull() || projectParam.IsUnknown() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			`Parameter "project" is null or unknown`,
		)
		return
	}
	data.ProjectParam = projectParam
	ctx = ctxvalues.With(ctx, "project", projectParam.String())

	service := data.ServiceParam.ValueString()

	apiRes, err := r.sdk.GetEnabledService(
		ctx,
		client.GetEnabledServiceRequest{
			Project: data.ProjectParam.ValueString(),
			Service: service,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Get Enabled Service",
			diag.FormatError(err),
		)
		return
	}

	if active, ok := apiRes.GetSpec().Active.Get(); !ok || !active {
		resp.Diagnostics.AddError(
			"Service Disabled",
			fmt.Sprintf("Service %q disabled", service),
		)
		return
	}
}

func (r *EnabledServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "EnabledServiceResource.Update")

	// This method simply syncs the state with the plan data since only timeout
	// attribute can be updated

	var data tfmodel.EnabledServiceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnabledServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "EnabledServiceResource.Delete")
	// This method is no-op because the service can't be disabled
}
