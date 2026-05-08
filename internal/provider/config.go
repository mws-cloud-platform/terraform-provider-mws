package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mws.cloud/util-toolset/pkg/os/env"
)

const (
	DefaultEndpoint = "https://api.mwsapis.ru"
	DefaultZone     = "ru-central1-a"

	EndpointEnv                        = "MWS_ENDPOINT"
	MWSTokenEnv                        = "MWS_TOKEN"
	ProjectEnv                         = "MWS_PROJECT"
	ZoneEnv                            = "MWS_ZONE"
	ServiceAccountAuthorizedKeyPathEnv = "MWS_SERVICE_ACCOUNT_AUTHORIZED_KEY_PATH"
)

// Config is a terraform provider configuration.
//
// See more information about fields in the provider [Schema].
type Config struct {
	Endpoint                        types.String `tfsdk:"endpoint"`
	MWSToken                        types.String `tfsdk:"mws_token"`
	ServiceAccountAuthorizedKeyPath types.String `tfsdk:"service_account_authorized_key_path"`
	Project                         types.String `tfsdk:"project"`
	Zone                            types.String `tfsdk:"zone"`
}

// Load loads environment variables and sets defaults for the unset fields.
func (c *Config) Load(e env.Env) {
	c.Endpoint = c.loadString(e, c.Endpoint, EndpointEnv, types.StringValue(DefaultEndpoint))
	c.MWSToken = c.loadString(e, c.MWSToken, MWSTokenEnv, types.StringNull())
	c.ServiceAccountAuthorizedKeyPath = c.loadString(e, c.ServiceAccountAuthorizedKeyPath, ServiceAccountAuthorizedKeyPathEnv, types.StringNull())
	c.Project = c.loadString(e, c.Project, ProjectEnv, types.StringNull())
	c.Zone = c.loadString(e, c.Zone, ZoneEnv, types.StringValue(DefaultZone))
}

func (*Config) loadString(e env.Env, field types.String, envKey string, defaultValue types.String) types.String {
	if IsValueSet(field) {
		return field
	}
	if value, ok := e.LookupEnv(envKey); ok {
		return types.StringValue(value)
	}
	return defaultValue
}
