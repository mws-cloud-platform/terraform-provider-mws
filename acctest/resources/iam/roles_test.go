package iam

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/service/resources/references/iam"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

var (
	//go:embed testdata/datasource/roles.tf
	rolesDataSourceTF string
)

func TestRolesSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RolesSuite))
}

type RolesSuite struct {
	utils.Suite
}

func (s *RolesSuite) TestRolesDataSource() {
	role := "admin"
	metadataID := iam.NewMustRoleID(role)

	steps := []resource.TestStep{
		{
			Config: fmt.Sprintf(rolesDataSourceTF, role),
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.mws_iam_roles.public", "items.#", "1"),
				resource.TestCheckResourceAttr("data.mws_iam_roles.public", "items.0.metadata.id", metadataID.String()),
			),
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
