package provider

import (
	tfdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	tfschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"go.mws.cloud/util-toolset/pkg/os/env"
)

const (
	TFEndpointVar                        = "endpoint"
	TFMWSTokenVar                        = "mws_token"
	TFServiceAccountAuthorizedKeyPathVar = "service_account_authorized_key_path"
	TFProjectVar                         = "project"
	TFZoneVar                            = "zone"
)

// Provider represents a base Terraform provider. This provider doesn't implement [tfprovider.Provider] interface,
// so it must be extended. See public/provider.go for an example.
type Provider struct {
	name        string
	version     string
	env         env.Env
	dataSources []func() tfdatasource.DataSource
	resources   []func() tfresource.Resource
}

// NewProvider creates a new base Terraform provider.
func NewProvider(opts ...Option) *Provider {
	p := &Provider{
		name:    "mws",
		version: "devel",
		env:     env.RealEnv{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Provider) Metadata(resp *tfprovider.MetadataResponse) {
	resp.TypeName = p.name
	resp.Version = p.version
}

func (p *Provider) Schema(resp *tfprovider.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Attributes: map[string]tfschema.Attribute{
			TFEndpointVar: tfschema.StringAttribute{
				MarkdownDescription: "Endpoint for the provider API calls, default value is <https://api.mwsapis.ru>. " +
					"This can also be specified using environment variable `MWS_BASE_ENDPOINT`.",
				Optional: true,
			},
			TFMWSTokenVar: tfschema.StringAttribute{
				MarkdownDescription: "IAM token for authentication. " +
					"This can also be specified using environment variable `MWS_TOKEN`. " +
					"Either `mws_token` or `service_account_authorized_key_path` must be specified. " +
					"Alternatively, you can run this provider inside a Compute VM with a linked service account.",
				Optional:  true,
				Sensitive: true,
			},
			TFServiceAccountAuthorizedKeyPathVar: tfschema.StringAttribute{
				MarkdownDescription: "Path to the service account authorization key file. " +
					"This can also be specified using environment variable `MWS_SERVICE_ACCOUNT_AUTHORIZED_KEY_PATH`. " +
					"Either `mws_token` or `service_account_authorized_key_path` must be specified. " +
					"Alternatively, you can run this provider inside a Compute VM with a linked service account.",
				Optional: true,
			},
			TFProjectVar: tfschema.StringAttribute{
				MarkdownDescription: "Project name. This can also be specified using environment variable `MWS_PROJECT`.",
				Optional:            true,
			},
			TFZoneVar: tfschema.StringAttribute{
				MarkdownDescription: "Zone name, default value is `ru-central1-a`." +
					"This can also be specified using environment variable `MWS_ZONE`.",
				Optional: true,
			},
		},
	}
}

func (p *Provider) DataSources() []func() tfdatasource.DataSource {
	return p.dataSources
}

func (p *Provider) Resources() []func() tfresource.Resource {
	return p.resources
}

func (p *Provider) Version() string {
	return p.version
}

func (p *Provider) Env() env.Env {
	return p.env
}
