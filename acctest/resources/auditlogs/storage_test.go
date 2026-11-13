package auditlogs

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	auditlogstest "go.mws.cloud/terraform-provider-mws/service/resources/auditlogs/acctest"
)

var (
	//go:embed testdata/storage.tf
	storageTF string
	//go:embed testdata/storage_updated.tf
	storageUpdatedTF string
	//go:embed testdata/datasource/storage.tf
	storageDataSourceTF string
)

func TestStorageSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StorageSuite))
}

type StorageSuite struct {
	utils.ResourceSuite
}

func (s *StorageSuite) TestStorage() {
	ctx := s.T().Context()

	tc, err := auditlogstest.StorageTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	storageName := utils.RandResourceName("storage")
	tc.ResourceConfig = fmt.Sprintf(storageTF, storageName)
	tc.UpdatedResourceConfig = fmt.Sprintf(storageUpdatedTF, storageName)
	tc.DataSourceConfig = fmt.Sprintf(storageDataSourceTF, storageName)

	s.BuildAndRun(ctx, tc)
}
