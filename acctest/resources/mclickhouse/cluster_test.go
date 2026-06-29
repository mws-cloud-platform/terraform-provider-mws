package mclickhouse

import (
	_ "embed"
	"fmt"
	"testing"

	helper "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mclickhousetest "go.mws.cloud/terraform-provider-mws/service/resources/mclickhouse/acctest"
)

var (
	//go:embed testdata/cluster.tf
	clusterTF string
	//go:embed testdata/datasource/cluster.tf
	clusterDataSourceTF string
)

func TestStandaloneClusterSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StandaloneClusterSuite))
}

type StandaloneClusterSuite struct {
	BaseSuite
}

func (s *StandaloneClusterSuite) TestStandaloneCluster() {
	ctx := s.T().Context()

	tc, err := mclickhousetest.ClusterTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	clusterName := utils.RandResourceName("cluster")
	adminPassword := helper.RandString(16)

	tc.ResourceConfig = fmt.Sprintf(clusterTF,
		clusterName,
		s.subnet.GetMetadata().GetId().ID(),
		adminPassword,
	)
	tc.DataSourceConfig = fmt.Sprintf(clusterDataSourceTF, clusterName)

	s.BuildAndRun(ctx, tc)
}
