package mpostgres

import (
	_ "embed"
	"fmt"
	"testing"

	helper "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/stretchr/testify/suite"
	mpostgresclient "go.mws.cloud/go-sdk/service/mpostgres/client"
	mpostgresmodel "go.mws.cloud/go-sdk/service/mpostgres/model"
	mpostgressdk "go.mws.cloud/go-sdk/service/mpostgres/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mpostgrestest "go.mws.cloud/terraform-provider-mws/service/resources/mpostgres/acctest"
)

var (
	//go:embed testdata/cluster_database.tf
	clusterDatabaseTF string
	//go:embed testdata/datasource/cluster_database.tf
	clusterDatabaseDataSourceTF string
)

func TestClusterDatabaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ClusterDatabaseSuite))
}

type ClusterDatabaseSuite struct {
	BaseClusterSuite
	userSDK *mpostgressdk.PostgresClusterUser

	userName string
	user     *mpostgresmodel.PostgresClusterUserResponse
}

func (s *ClusterDatabaseSuite) TestClusterDatabase() {
	ctx := s.T().Context()

	database := utils.RandResourceName("db")

	tc, err := mpostgrestest.ClusterDatabaseTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(clusterDatabaseTF,
		s.clusterName,
		database,
		s.user.GetMetadata().GetId().ID(),
	)
	tc.DataSourceConfig = fmt.Sprintf(clusterDatabaseDataSourceTF, s.clusterName, database)

	s.BuildAndRun(ctx, tc)
}

func (s *ClusterDatabaseSuite) SetupSuite() {
	ctx := s.T().Context()
	s.BaseClusterSuite.SetupSuite()

	var err error
	s.userSDK, err = mpostgressdk.NewPostgresClusterUser(ctx, s.SDK)
	s.Require().NoError(err)

	s.userName = utils.RandResourceName("user")
	s.user, err = s.userSDK.CreatePostgresClusterUser(ctx, mpostgresclient.UpsertPostgresClusterUserRequest{
		Cluster: s.clusterName,
		User:    s.userName,
		Body: mpostgresmodel.PostgresClusterUserRequest{
			Spec: mpostgresmodel.PostgresClusterUserSpecRequest{
				Password: helper.RandString(16),
				Role:     new(mpostgresmodel.PostgresUserRole_DB_OWNER_USER),
			},
		},
	})
	s.Require().NoError(err)
}

func (s *ClusterDatabaseSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.userSDK.DeletePostgresClusterUser(ctx, mpostgresclient.DeletePostgresClusterUserRequest{
		Cluster: s.clusterName,
		User:    s.userName,
	}, mpostgresclient.WithWait()); err != nil {
		s.T().Logf(
			"user %q deletion failed: %v",
			s.userName,
			err,
		)
	} else {
		s.T().Logf(
			"user %q deleted",
			s.userName,
		)
	}

	s.BaseClusterSuite.TearDownSuite()
}
