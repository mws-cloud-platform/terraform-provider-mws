package iam

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/service/resources/references/iam"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

var (
	//go:embed testdata/datasource/role.tf
	roleDataSourceTF string
)

func TestRoleSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RoleSuite))
}

type RoleSuite struct {
	utils.Suite
}

func (s *RoleSuite) TestRole() {
	role := "admin"
	metadataID := iam.NewMustRoleID(role)

	steps := []resource.TestStep{
		{
			Config: fmt.Sprintf(roleDataSourceTF, role),
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				return resource.TestCheckResourceAttr("data.mws_iam_role.role_data", "metadata.id", metadataID.String())(state)
			},
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
