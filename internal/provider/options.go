package provider

import (
	tfdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"go.mws.cloud/util-toolset/pkg/os/env"
)

// Option is a function that configures a terraform provider.
type Option func(*Provider)

// WithName sets the name of the terraform provider.
func WithName(name string) Option {
	return func(p *Provider) {
		p.name = name
	}
}

// WithVersion sets the version of the terraform provider.
func WithVersion(version string) Option {
	return func(p *Provider) {
		p.version = version
	}
}

// WithEnv sets the environment of the terraform provider.
func WithEnv(e env.Env) Option {
	return func(p *Provider) {
		p.env = e
	}
}

// WithDataSources sets the data sources of the terraform provider.
func WithDataSources(dataSources []func() tfdatasource.DataSource) Option {
	return func(p *Provider) {
		p.dataSources = dataSources
	}
}

// WithResources sets the resources of the terraform provider.
func WithResources(resources []func() tfresource.Resource) Option {
	return func(p *Provider) {
		p.resources = resources
	}
}
