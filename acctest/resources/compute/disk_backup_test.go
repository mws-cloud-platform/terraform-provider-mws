package compute

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computemodel "go.mws.cloud/go-sdk/service/compute/model"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	computetest "go.mws.cloud/terraform-provider-mws/service/resources/compute/acctest"
)

var (
	//go:embed testdata/disk_backup.tf
	diskBackupTF string
	//go:embed testdata/datasource/disk_backup.tf
	diskBackupDataSourceTF string
)

func TestDiskBackupSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiskBackupSuite))
}

type DiskBackupSuite struct {
	utils.ResourceSuite
	diskSDK  *computesdk.Disk
	diskName string
	disk     *computemodel.DiskOptionalResponse
}

func (s *DiskBackupSuite) SetupSuite() {
	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.diskName = utils.RandResourceName("backup-disk")
	var err error

	s.diskSDK, err = computesdk.NewDisk(ctx, s.SDK)
	s.Require().NoError(err)

	s.disk, err = s.diskSDK.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Disk: s.diskName,
		Body: computemodel.DiskRequest{
			Spec: computemodel.DiskSpecRequest{
				Zone:      "ru-central1-a",
				DiskType:  new(computeref.NewDiskTypeRef("nbs-pl2")),
				Size:      new(bytesize.MustNewFromInt64(10, bytesize.GB)),
				Iops:      new(computemodel.Iops(1000)),
				BlockSize: new(bytesize.MustNewFromInt64(4, bytesize.KB)),
			},
		},
	}, computeclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("disk %q created", s.diskName)
}

func (s *DiskBackupSuite) TearDownSuite() {
	ctx := s.T().Context()
	err := s.diskSDK.DeleteDisk(ctx, computeclient.DeleteDiskRequest{
		Disk: s.diskName,
	}, computeclient.WithWait())
	s.Require().NoError(err)

	s.ResourceSuite.TearDownSuite()
}

func (s *DiskBackupSuite) TestDiskBackup() {
	ctx := s.T().Context()
	tc, err := computetest.DiskBackupTestCase(ctx, s.SDK)
	s.Require().NoError(err)
	diskRef, err := computeref.ParseDiskRef(s.T().Context(), s.disk.GetMetadata().GetId().ID())
	s.Require().NoError(err)
	tc.ResourceConfig = fmt.Sprintf(diskBackupTF, s.diskName, diskRef.IDPath())
	tc.DataSourceConfig = fmt.Sprintf(diskBackupDataSourceTF, s.diskName)

	s.BuildAndRun(ctx, tc)
}
