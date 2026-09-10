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
	//go:embed testdata/datasource/vm_type.tf
	vmTypeDataSourceTF string
)

func TestVMTypeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(VMTypeSuite))
}

type VMTypeSuite struct {
	utils.Suite
}

func (s *VMTypeSuite) TestVMType() {
	steps := []resource.TestStep{
		{
			Config: vmTypeDataSourceTF,
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				return resource.TestCheckResourceAttr("data.mws_compute_vm_type.base", "vm_type", "base-4-8")(state)
			},
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
