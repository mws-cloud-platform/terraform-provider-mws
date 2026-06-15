package examples

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	"go.mws.cloud/terraform-provider-mws/internal/acctest"
)

func TestDataSourcesExamples(t *testing.T) {
	t.Parallel()

	providerConfig := string(acctest.NewEmptyProviderConfig("mws").Bytes())

	root, err := filepath.Abs("../../examples/data-sources")
	require.NoError(t, err)

	names, err := listExampleResources(root)
	require.NoError(t, err)

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(filepath.Join(root, name, "data-source.tf"))
			require.NoError(t, err)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config:      providerConfig + "\n" + string(data),
						ExpectError: regexp.MustCompile("(?i).*not.found.*"),
					},
				},
			})
		})
	}
}
