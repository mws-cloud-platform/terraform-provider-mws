package public

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	base "go.mws.cloud/terraform-provider-mws/internal/provider"
)

type Provider struct {
	base *base.Provider
}

var _ tfprovider.Provider = (*Provider)(nil)

func NewProvider(opts ...base.Option) func() tfprovider.Provider {
	return func() tfprovider.Provider {
		return &Provider{
			base: base.NewProvider(opts...),
		}
	}
}

func (p Provider) Metadata(_ context.Context, _ tfprovider.MetadataRequest, resp *tfprovider.MetadataResponse) {
	p.base.Metadata(resp)
}

func (p Provider) Schema(_ context.Context, _ tfprovider.SchemaRequest, resp *tfprovider.SchemaResponse) {
	p.base.Schema(resp)
}

func (p Provider) Configure(ctx context.Context, req tfprovider.ConfigureRequest, resp *tfprovider.ConfigureResponse) {
	config := &Config{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Load(p.base.Env())

	sdk, err := LoadSDKFromConfig(ctx, config, p.base.Env(), p.base.Version())
	if err != nil {
		resp.Diagnostics.AddError("Failed to load SDK", err.Error())
		return
	}

	data := &Data{
		Config: config,
		SDK:    sdk,
	}

	resp.ResourceData = data
	resp.DataSourceData = data
}

func (p Provider) DataSources(context.Context) []func() datasource.DataSource {
	return p.base.DataSources()
}

func (p Provider) Resources(context.Context) []func() resource.Resource {
	return p.base.Resources()
}
