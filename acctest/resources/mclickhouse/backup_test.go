package mclickhouse

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mclickhousetest "go.mws.cloud/terraform-provider-mws/service/resources/mclickhouse/acctest"
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

	tc, err := mclickhousetest.BackupTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(backupTF, s.clusterName, backup)
	tc.DataSourceConfig = fmt.Sprintf(backupDataSourceTF, s.clusterName, backup)
	s.BuildAndRun(ctx, tc)
}
