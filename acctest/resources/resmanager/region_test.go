package resmanager

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

var (
	//go:embed testdata/datasource/region.tf
	regionDataSourceTF string
)

func TestResmanagerRegionSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ResmanagerRegionSuite))
}

type ResmanagerRegionSuite struct {
	utils.Suite
}

func (s *ResmanagerRegionSuite) TestResmanagerRegionDataSource() {
	regionName := "ru-central1"

	steps := []resource.TestStep{
		{
			Config: fmt.Sprintf(regionDataSourceTF, regionName),
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				return resource.TestCheckResourceAttr("data.mws_resmanager_region.region_data", "region", regionName)(state)
			},
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
