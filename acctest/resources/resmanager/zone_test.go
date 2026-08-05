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
	//go:embed testdata/datasource/zone.tf
	zoneDataSourceTF string
)

func TestResmanagerZoneSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ResmanagerZoneSuite))
}

type ResmanagerZoneSuite struct {
	utils.Suite
}

func (s *ResmanagerZoneSuite) TestResmanagerZoneDataSource() {
	const zoneName = "ru-central1-a"

	steps := []resource.TestStep{
		{
			Config: fmt.Sprintf(zoneDataSourceTF, zoneName),
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				return resource.TestCheckResourceAttr("data.mws_resmanager_zone.zone_data", "zone", zoneName)(state)
			},
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
