package examples

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	"go.mws.cloud/terraform-provider-mws/internal/acctest"
)

func TestResourcesExamples(t *testing.T) {
	t.Parallel()

	providerConfig := string(acctest.NewEmptyProviderConfig("mws").Bytes())

	root, err := filepath.Abs("../../examples/resources")
	require.NoError(t, err)

	names, err := listExampleResources(root)
	require.NoError(t, err)

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resourceConfig, err := os.ReadFile(filepath.Join(root, name, "resource.tf"))
			require.NoError(t, err)

			configVariables, err := randomizeResourceNames(resourceConfig)
			require.NoError(t, err)

			configVariables["private_key_file_path"] = tfconfig.StringVariable("testdata/private_key")
			configVariables["certificate_file_path"] = tfconfig.StringVariable("testdata/certificate")

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config:          providerConfig + "\n" + string(resourceConfig),
						ConfigVariables: configVariables,
					},
				},
			})
		})
	}
}

// randomizeResourceNames adds random suffix to all variables with a _name
// suffix and default value.
func randomizeResourceNames(config []byte) (tfconfig.Variables, error) {
	f, diags := hclwrite.ParseConfig(config, "resource.tf", hcl.InitialPos)
	if diags.HasErrors() {
		return nil, diags
	}

	suffix := utils.RandResourceName("-ex")

	variables := make(tfconfig.Variables, 0)

	for _, block := range f.Body().Blocks() {
		if block.Type() != "variable" {
			continue
		}

		labels := block.Labels()
		if len(labels) != 1 || !strings.HasSuffix(labels[0], "_name") {
			continue
		}

		defaultAttr := block.Body().GetAttribute("default")
		if defaultAttr == nil {
			continue
		}

		v := bytes.Trim(defaultAttr.Expr().BuildTokens(nil).Bytes(), `" `)

		variables[labels[0]] = tfconfig.StringVariable(string(v) + suffix)
	}

	return variables, nil
}
