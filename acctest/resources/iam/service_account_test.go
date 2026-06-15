package iam

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	iamtest "go.mws.cloud/terraform-provider-mws/service/resources/iam/acctest"
)

var (
	//go:embed testdata/service_account.tf
	serviceAccountTF string
	//go:embed testdata/datasource/service_account.tf
	serviceAccountDataSourceTF string
)

func TestServiceAccountSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceAccountSuite))
}

type ServiceAccountSuite struct {
	utils.ResourceSuite
}

func (s *ServiceAccountSuite) TestServiceAccount() {
	ctx := s.T().Context()

	tc, err := iamtest.ServiceAccountTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	serviceAccountName := utils.RandResourceName("sa")
	tc.ResourceConfig = fmt.Sprintf(serviceAccountTF, serviceAccountName)
	tc.DataSourceConfig = fmt.Sprintf(serviceAccountDataSourceTF, serviceAccountName)

	s.BuildAndRun(ctx, tc)
}
