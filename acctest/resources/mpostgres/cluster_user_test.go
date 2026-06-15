package mpostgres

import (
	_ "embed"
	"fmt"
	"testing"

	helper "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mpostgrestest "go.mws.cloud/terraform-provider-mws/service/resources/mpostgres/acctest"
)

var (
	//go:embed testdata/cluster_user.tf
	clusterUserTF string
	//go:embed testdata/datasource/cluster_user.tf
	clusterUserDataSourceTF string
)

func TestClusterUserSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ClusterUserSuite))
}

type ClusterUserSuite struct {
	BaseClusterSuite
}

func (s *ClusterUserSuite) TestClusterUser() {
	ctx := s.T().Context()

	user := utils.RandResourceName("user")
	password := helper.RandString(16)

	tc, err := mpostgrestest.ClusterUserTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(clusterUserTF,
		s.clusterName,
		user,
		password,
	)
	tc.DataSourceConfig = fmt.Sprintf(clusterUserDataSourceTF, s.clusterName, user)

	s.BuildAndRun(ctx, tc)
}
