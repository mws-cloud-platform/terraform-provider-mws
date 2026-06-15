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
	//go:embed testdata/snapshot_from_disk.tf
	snapshotFromDiskTF string
	//go:embed testdata/datasource/snapshot.tf
	snapshotDataSourceTF string
)

func TestSnapshoFromDiskSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SnapshotFromDiskSuite))
}

type SnapshotFromDiskSuite struct {
	baseSnapshotSuite
}

func (s *SnapshotFromDiskSuite) TestSnapshot() {
	ctx := s.T().Context()
	tc, err := computetest.SnapshotTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	snapshotName := utils.RandResourceName("snapshot")
	tc.ResourceConfig = fmt.Sprintf(snapshotFromDiskTF,
		snapshotName,
		s.disk.GetMetadata().GetId().ID(),
	)
	tc.DataSourceConfig = fmt.Sprintf(snapshotDataSourceTF, snapshotName)

	s.BuildAndRun(ctx, tc)
}

type baseSnapshotSuite struct {
	utils.ResourceSuite

	imageSDK *computesdk.Image
	diskSDK  *computesdk.Disk

	imageName string
	diskName  string
	disk      *computemodel.DiskOptionalResponse
}

func (s *baseSnapshotSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.diskName = utils.RandResourceName("disk")
	s.imageName = utils.RandResourceName("image")

	s.imageSDK, err = computesdk.NewImage(ctx, s.SDK)
	s.Require().NoError(err)
	s.diskSDK, err = computesdk.NewDisk(ctx, s.SDK)
	s.Require().NoError(err)

	image, err := s.imageSDK.GetImage(ctx, computeclient.GetImageRequest{
		Project: ubuntuProjectName,
		Image:   ubuntuImageName,
	})
	s.Require().NoError(err)
	s.Require().NotNil(image.GetMetadata().GetId())

	imageRef, err := computeref.ParseImageRef(s.T().Context(), image.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.disk, err = s.diskSDK.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Disk: s.diskName,
		Body: computemodel.DiskRequest{
			Spec: computemodel.DiskSpecRequest{
				Zone:     "ru-central1-a",
				DiskType: new(computeref.NewDiskTypeRef("nbs-pl2")),
				Size:     new(bytesize.MustNewFromInt64(10, bytesize.GB)),
				Source: &computemodel.DiskSpecSourceRequest{
					Image: new(imageRef),
				},
				Iops:      new(computemodel.Iops(1000)),
				BlockSize: new(bytesize.MustNewFromInt64(4, bytesize.KB)),
			},
		},
	}, computeclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("disk %q created", s.diskName)
}

func (s *baseSnapshotSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.diskSDK.DeleteDisk(ctx, computeclient.DeleteDiskRequest{
		Disk: s.diskName,
	}, computeclient.WithWait()); err != nil {
		s.T().Logf("disk %q deletion failed: %v", s.diskName, err)
	} else {
		s.T().Logf("disk %q deleted", s.diskName)
	}

	if err := s.imageSDK.DeleteImage(ctx, computeclient.DeleteImageRequest{
		Image: s.imageName,
	}, computeclient.WithWait()); err != nil {
		s.T().Logf("image %q deletion failed: %v", s.imageName, err)
	} else {
		s.T().Logf("image %q deleted", s.imageName)
	}

	s.ResourceSuite.TearDownSuite()
}
