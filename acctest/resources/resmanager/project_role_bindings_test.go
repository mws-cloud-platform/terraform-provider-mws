package resmanager

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/suite"
	commonmodel "go.mws.cloud/go-sdk/service/common/model"
	iamref "go.mws.cloud/go-sdk/service/resources/references/iam"
	resmanagerclient "go.mws.cloud/go-sdk/service/rm/client"
	resourcesdk "go.mws.cloud/go-sdk/service/rm/sdk"

	iamacctest "go.mws.cloud/terraform-provider-mws/acctest/resources/iam"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

var (
	//go:embed testdata/datasource/project_role_bindings.tf
	projectRoleBindingsDataSourceTF string
)

var iamRoleAdminRef = iamref.NewMustRoleRef("admin")

func TestProjectRoleBindingsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProjectRoleBindingsSuite))
}

type ProjectRoleBindingsSuite struct {
	iamacctest.BaseServiceAccountSuite

	project        string
	roleBinding    string
	roleBindingSDK *resourcesdk.ProjectRoleBinding
}

func (s *ProjectRoleBindingsSuite) SetupSuite() {
	var err error

	s.BaseServiceAccountSuite.SetupSuite()
	ctx := s.T().Context()

	s.project = s.SDK.DefaultProject()
	s.roleBinding = utils.RandResourceName("role-binding")

	s.roleBindingSDK, err = resourcesdk.NewProjectRoleBinding(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.roleBindingSDK.CreateProjectRoleBinding(ctx, resmanagerclient.UpsertProjectRoleBindingRequest{
		Project:     s.project,
		RoleBinding: s.roleBinding,
		Body: commonmodel.CommonRoleBindingRequest{
			Spec: commonmodel.CommonRoleBindingSpecRequest{
				Subject: commonmodel.CommonRoleBindingSpecSubjectRequest{
					ServiceAccount: &s.ServiceAccountRef,
				},
				Role: iamRoleAdminRef,
			},
		},
	}, resmanagerclient.WithWait())
	s.Require().NoError(err, "create project role binding")

	s.T().Logf("project role binding %q created", s.roleBinding)
}

func (s *ProjectRoleBindingsSuite) TearDownSuite() {
	ctx := s.T().Context()

	err := s.roleBindingSDK.DeleteProjectRoleBinding(ctx, resmanagerclient.DeleteProjectRoleBindingRequest{
		Project:     s.project,
		RoleBinding: s.roleBinding,
	}, resmanagerclient.WithWait())
	if err != nil {
		s.T().Logf("project role binding %q deletion failed: %v", s.roleBinding, err)
	} else {
		s.T().Logf("project role binding %q deleted", s.roleBinding)
	}

	s.BaseServiceAccountSuite.TearDownSuite()
}

func (s *ProjectRoleBindingsSuite) TestProjectRoleBindingsDataSource() {
	config := fmt.Sprintf(projectRoleBindingsDataSourceTF, s.roleBinding)

	steps := []resource.TestStep{
		{
			Config: config,
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.mws_resmanager_project_role_bindings.filtered", "items.#", "1"),
				resource.TestCheckResourceAttr("data.mws_resmanager_project_role_bindings.filtered", "items.0.role", iamRoleAdminRef.String()),
			),
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
