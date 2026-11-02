package metastore

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/resources/iam"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	metastoretest "go.mws.cloud/terraform-provider-mws/service/resources/metastore/acctest"
)

var (
	//go:embed testdata/iceberg.tf
	icebergTF string
	//go:embed testdata/datasource/iceberg.tf
	icebergDataSourceTF string
)

func TestIcebergSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IcebergSuite))
}

type IcebergSuite struct {
	iam.BaseServiceAccountSuite
}

func (s *IcebergSuite) TestIceberg() {
	ctx := s.T().Context()

	tc, err := metastoretest.IcebergTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	catalogName := utils.RandResourceName("catalog")

	tc.ResourceConfig = fmt.Sprintf(icebergTF,
		catalogName,
		s.ServiceAccountRef.IDPath(),
		fmt.Sprintf("s3://%s-iceberg-warehouse/%s", s.SDK.DefaultProject(), catalogName),
	)
	tc.DataSourceConfig = fmt.Sprintf(icebergDataSourceTF, catalogName)
	s.T().Log(tc.ResourceConfig)

	s.BuildAndRun(ctx, tc)
}
