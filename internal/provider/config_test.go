package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/os/env"

	"go.mws.cloud/terraform-provider-mws/internal/provider"
)

func TestConfig_Load(t *testing.T) {
	for _, tc := range []struct {
		name     string
		config   provider.Config
		env      env.Env
		expected provider.Config
	}{
		{
			name:   "empty",
			config: provider.Config{},
			env:    env.MapEnv{},
			expected: provider.Config{
				Endpoint: types.StringValue(provider.DefaultEndpoint),
				MWSToken: types.StringNull(),
				Project:  types.StringNull(),
				Zone:     types.StringValue(provider.DefaultZone),
			},
		},
		{
			name: "config",
			config: provider.Config{
				Endpoint:                        types.StringValue("https://example.com"),
				MWSToken:                        types.StringValue("token"),
				Project:                         types.StringValue("project"),
				Zone:                            types.StringValue("zone"),
				ServiceAccountAuthorizedKeyPath: types.StringValue("ServiceAccountAuthorizedKeyPath"),
			},
			env: env.MapEnv{ // check priority
				provider.EndpointEnv:                        "https://foo.com",
				provider.MWSTokenEnv:                        "bar",
				provider.ProjectEnv:                         "baz",
				provider.ZoneEnv:                            "qux",
				provider.ServiceAccountAuthorizedKeyPathEnv: "saakp",
			},
			expected: provider.Config{
				Endpoint:                        types.StringValue("https://example.com"),
				MWSToken:                        types.StringValue("token"),
				Project:                         types.StringValue("project"),
				Zone:                            types.StringValue("zone"),
				ServiceAccountAuthorizedKeyPath: types.StringValue("ServiceAccountAuthorizedKeyPath"),
			},
		},
		{
			name:   "envs",
			config: provider.Config{},
			env: env.MapEnv{
				provider.EndpointEnv:                        "https://example.com",
				provider.MWSTokenEnv:                        "token",
				provider.ProjectEnv:                         "project",
				provider.ZoneEnv:                            "zone",
				provider.ServiceAccountAuthorizedKeyPathEnv: "test/path",
			},
			expected: provider.Config{
				Endpoint:                        types.StringValue("https://example.com"),
				MWSToken:                        types.StringValue("token"),
				Project:                         types.StringValue("project"),
				Zone:                            types.StringValue("zone"),
				ServiceAccountAuthorizedKeyPath: types.StringValue("test/path"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.config.Load(tc.env)
			require.Equal(t, tc.expected, tc.config)
		})
	}
}
