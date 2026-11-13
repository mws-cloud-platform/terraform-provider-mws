package auditlogs

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/duration"
	auditlogsclient "go.mws.cloud/go-sdk/service/auditlogs/client"
	auditlogsmodel "go.mws.cloud/go-sdk/service/auditlogs/model"
	auditlogssdk "go.mws.cloud/go-sdk/service/auditlogs/sdk"
	auditlogsref "go.mws.cloud/go-sdk/service/resources/references/auditlogs"

	iamacctest "go.mws.cloud/terraform-provider-mws/acctest/resources/iam"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	auditlogstest "go.mws.cloud/terraform-provider-mws/service/resources/auditlogs/acctest"
)

var (
	//go:embed testdata/collector.tf
	collectorTF string
	//go:embed testdata/collector_updated.tf
	collectorUpdatedTF string
	//go:embed testdata/datasource/collector.tf
	collectorDataSourceTF string
)

func TestCollectorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CollectorSuite))
}

type CollectorSuite struct {
	iamacctest.BaseServiceAccountSuite

	storageSDK  *auditlogssdk.Storage
	storageName string
	storageRef  auditlogsref.StorageRef
}

func (s *CollectorSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseServiceAccountSuite.SetupSuite()

	s.storageSDK, err = auditlogssdk.NewStorage(ctx, s.SDK)
	s.Require().NoError(err)

	s.storageName = utils.RandResourceName("storage")

	_, err = s.storageSDK.CreateStorage(ctx, auditlogsclient.UpsertStorageRequest{
		StorageName: s.storageName,
		Body: auditlogsmodel.StorageRequest{
			Spec: auditlogsmodel.AuditLogsStorageSpecRequest{
				Retention: auditlogsmodel.AuditLogsStorageSpecRetentionRequest{
					Size:   new(bytesize.MustParseString("10GB")),
					Period: duration.MustParseString("30d"),
				},
			},
		},
	}, auditlogsclient.WithWait())
	s.Require().NoError(err)

	s.storageRef, err = auditlogsref.NewStorageRef(s.SDK.DefaultProject(), s.storageName)
	s.Require().NoError(err)

	s.T().Logf("storage %q created", s.storageName)
}

func (s *CollectorSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.storageSDK.DeleteStorage(ctx, auditlogsclient.DeleteStorageRequest{
		StorageName: s.storageName,
	}, auditlogsclient.WithWait()); err != nil {
		s.T().Logf("storage %q deletion failed: %v", s.storageName, err)
	} else {
		s.T().Logf("storage %q deleted", s.storageName)
	}

	s.BaseServiceAccountSuite.TearDownSuite()
}

func (s *CollectorSuite) TestCollector() {
	ctx := s.T().Context()

	collectorName := utils.RandResourceName("collector")

	tc, err := auditlogstest.CollectorTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	serviceAccountRef := s.ServiceAccountRef.IDPath()
	storageRef := s.storageRef.IDPath()
	project := s.SDK.DefaultProject()

	tc.ResourceConfig = fmt.Sprintf(collectorTF, collectorName, storageRef, project, serviceAccountRef)
	tc.UpdatedResourceConfig = fmt.Sprintf(collectorUpdatedTF, collectorName, storageRef, project, serviceAccountRef)
	tc.DataSourceConfig = fmt.Sprintf(collectorDataSourceTF, collectorName)

	s.BuildAndRun(ctx, tc)
}
