package resmanager

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
	//go:embed testdata/datasource/regions.tf
	regionsDataSourceTF string
)

func TestResmanagerRegionsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ResmanagerRegionsSuite))
}

type ResmanagerRegionsSuite struct {
	utils.Suite
}

func (s *ResmanagerRegionsSuite) TestResmanagerRegionsDataSource() {
	steps := []resource.TestStep{
		{
			Config: regionsDataSourceTF,
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				return resource.TestCheckResourceAttrSet("data.mws_resmanager_regions.all", "items.0.metadata.id")(state)
			},
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
