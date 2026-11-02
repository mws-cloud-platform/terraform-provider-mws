package resmanager

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	iamacctest "go.mws.cloud/terraform-provider-mws/acctest/resources/iam"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	rmtest "go.mws.cloud/terraform-provider-mws/service/resources/rm/acctest"
)

var (
	//go:embed testdata/project_role_binding.tf
	projectRoleBindingTF string
	//go:embed testdata/datasource/project_role_binding.tf
	projectRoleBindingDataSourceTF string
)

func TestProjectRoleBindingSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProjectRoleBindingSuite))
}

type ProjectRoleBindingSuite struct {
	iamacctest.BaseServiceAccountSuite
}

func (s *ProjectRoleBindingSuite) TestProjectRoleBinding() {
	ctx := s.T().Context()

	tc, err := rmtest.ProjectRoleBindingTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	roleBindingName := utils.RandResourceName("role-binding")
	serviceAccountRef := s.ServiceAccountRef.IDPath()
	tc.ResourceConfig = fmt.Sprintf(projectRoleBindingTF, roleBindingName, serviceAccountRef)
	tc.DataSourceConfig = fmt.Sprintf(projectRoleBindingDataSourceTF, roleBindingName)

	s.BuildAndRun(ctx, tc)
}
