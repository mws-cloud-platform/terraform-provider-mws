package mpostgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mpostgrestest "go.mws.cloud/terraform-provider-mws/service/resources/mpostgres/acctest"
)

var (
	//go:embed testdata/backup.tf
	backupTF string
	//go:embed testdata/datasource/backup.tf
	backupDataSourceTF string
)

func TestBackupSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BackupSuite))
}

type BackupSuite struct {
	BaseClusterSuite
}

func (s *BackupSuite) TestBackup() {
	ctx := s.T().Context()

	backup := utils.RandResourceName("backup")

	tc, err := mpostgrestest.BackupTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	backupExists := tc.ResourceExists
	clusterOK := func(ctx context.Context) error {
		return waitForClusterReady(ctx, s.clusterSDK, s.clusterName)
	}

	tc.ResourceConfig = fmt.Sprintf(backupTF, s.clusterName, backup)
	tc.DataSourceConfig = fmt.Sprintf(backupDataSourceTF, s.clusterName, backup)
	// we must ensure cluster state is OK, before deleting a backup
	tc.ResourceExists = func(ctx context.Context, id string) error {
		return errors.Join(backupExists(ctx, id), clusterOK(ctx))
	}
	// we must ensure cluster state is OK, before creating a backup
	s.Require().NoError(clusterOK(ctx))
	s.BuildAndRun(ctx, tc)
}
