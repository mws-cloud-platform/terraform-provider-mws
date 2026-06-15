package mpostgres

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	mpostgressdk "go.mws.cloud/go-sdk/service/mpostgres/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mpostgrestest "go.mws.cloud/terraform-provider-mws/service/resources/mpostgres/acctest"
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

	clusterSDK, err := mpostgressdk.NewPostgresCluster(ctx, s.SDK)
	s.Require().NoError(err)

	tc, err := mpostgrestest.ClusterTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	clusterName := utils.RandResourceName("cluster")
	tc.ResourceConfig = fmt.Sprintf(clusterTF,
		clusterName,
		s.network.GetMetadata().GetId().ID(),
		s.primaryEndpointAddress.GetMetadata().GetId().ID(),
	)
	tc.DataSourceConfig = fmt.Sprintf(clusterDataSourceTF, clusterName)
	tc.ResourceExists = func(ctx context.Context, _ string) error {
		return waitForClusterReady(ctx, clusterSDK, clusterName)
	}

	s.BuildAndRun(ctx, tc)
}
