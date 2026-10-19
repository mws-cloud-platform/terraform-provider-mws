package compute

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

var (
	//go:embed testdata/datasource/vm_types.tf
	vmTypesDataSourceTF string
)

func TestVMTypesSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(VMTypesSuite))
}

type VMTypesSuite struct {
	utils.Suite
}

func (s *VMTypesSuite) TestVMTypesDataSource() {
	steps := []resource.TestStep{
		{
			Config: vmTypesDataSourceTF,
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				return resource.TestCheckResourceAttrSet("data.mws_compute_vm_types.all", "items.0.metadata.id")(state)
			},
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
