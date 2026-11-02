package metastore

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	metastoreclient "go.mws.cloud/go-sdk/service/metastore/client"
	metastoremodel "go.mws.cloud/go-sdk/service/metastore/model"
	metastoresdk "go.mws.cloud/go-sdk/service/metastore/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/resources/iam"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	metastoretest "go.mws.cloud/terraform-provider-mws/service/resources/metastore/acctest"
)

var (
	//go:embed testdata/iceberg_role_binding.tf
	icebergRoleBindingTF string
	//go:embed testdata/datasource/iceberg_role_binding.tf
	icebergRoleBindingDataSourceTF string
)

func TestIcebergRoleBindingSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IcebergRoleBindingSuite))
}

type IcebergRoleBindingSuite struct {
	iam.BaseServiceAccountSuite

	icebergSDK  *metastoresdk.Iceberg
	catalogName string
}

func (s *IcebergRoleBindingSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseServiceAccountSuite.SetupSuite()

	s.catalogName = utils.RandResourceName("catalog")

	s.icebergSDK, err = metastoresdk.NewIceberg(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.icebergSDK.CreateIcebergCatalog(ctx, metastoreclient.UpsertIcebergCatalogRequest{
		CatalogName: s.catalogName,
		Body: metastoremodel.IcebergCatalogRequest{
			Spec: metastoremodel.CatalogSpecRequest{
				Properties: map[string]string{
					"location": fmt.Sprintf("s3://%s-iceberg-warehouse/%s", s.SDK.DefaultProject(), s.catalogName),
				},
				ServiceAccount: &metastoremodel.CatalogSpecServiceAccountRequest{
					Ref: s.ServiceAccountRef,
				},
			},
		},
	}, metastoreclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("iceberg catalog %q created", s.catalogName)
}

func (s *IcebergRoleBindingSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.icebergSDK.DeleteIcebergCatalog(ctx, metastoreclient.DeleteIcebergCatalogRequest{
		Project:     s.SDK.DefaultProject(),
		CatalogName: s.catalogName,
	}, metastoreclient.WithWait()); err != nil {
		s.T().Logf("iceberg catalog %q deletion failed: %v", s.catalogName, err)
	} else {
		s.T().Logf("iceberg catalog %q deleted", s.catalogName)
	}

	s.BaseServiceAccountSuite.TearDownSuite()
}

func (s *IcebergRoleBindingSuite) TestIcebergRoleBinding() {
	ctx := s.T().Context()

	tc, err := metastoretest.IcebergRoleBindingTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	roleBindingName := "rb__" + utils.RandResourceName("iceberg-viewer")

	tc.ResourceConfig = fmt.Sprintf(icebergRoleBindingTF,
		s.catalogName,
		roleBindingName,
		s.ServiceAccountRef.IDPath(),
	)
	tc.DataSourceConfig = fmt.Sprintf(icebergRoleBindingDataSourceTF, roleBindingName)
	s.T().Log(tc.ResourceConfig)

	s.BuildAndRun(ctx, tc)
}
