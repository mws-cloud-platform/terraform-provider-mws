package compute

import (
	_ "embed"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	computetest "go.mws.cloud/terraform-provider-mws/service/resources/compute/acctest"
)

var (
	//go:embed testdata/disk.tf
	diskTF string
	//go:embed testdata/datasource/disk.tf
	diskDataSourceTF string
	//go:embed testdata/disk_with_image.tf
	diskWithImageTF string
)

func TestDiskSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiskSuite))
}

func TestDiskWithImageSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiskWithImageSuite))
}

func TestDiskIncreaseSizeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiskIncreaseSizeSuite))
}

func TestDiskDecreaseSizeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiskDecreaseSizeSuite))
}

type DiskSuite struct {
	utils.ResourceSuite
}

func (s *DiskSuite) TestDisk() {
	ctx := s.T().Context()

	tc, err := computetest.DiskTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	diskName := utils.RandResourceName("disk")
	tc.ResourceConfig = fmt.Sprintf(diskTF, diskName, "1GB")
	tc.DataSourceConfig = fmt.Sprintf(diskDataSourceTF, diskName)

	s.BuildAndRun(ctx, tc)
}

type DiskWithImageSuite struct {
	utils.ResourceSuite
}

func (s *DiskWithImageSuite) TestDiskWithImage() {
	ctx := s.T().Context()

	tc, err := computetest.DiskTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	diskName := utils.RandResourceName("disk")
	tc.ResourceConfig = fmt.Sprintf(diskWithImageTF, diskName, "4GB", ubuntuImageID)

	s.BuildAndRun(ctx, tc)
}

type DiskIncreaseSizeSuite struct {
	utils.ResourceSuite
}

func (s *DiskIncreaseSizeSuite) TestDiskIncreaseSize() {
	ctx := s.T().Context()

	tc, err := computetest.DiskTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	diskName := utils.RandResourceName("disk")
	tc.ResourceConfig = fmt.Sprintf(diskTF, diskName, "1GB")
	tc.DataSourceConfig = fmt.Sprintf(diskDataSourceTF, diskName)
	tc.UpdatedResourceConfig = fmt.Sprintf(diskTF, diskName, "2GB")

	s.BuildAndRun(ctx, tc)
}

type DiskDecreaseSizeSuite struct {
	utils.ResourceSuite
}

func (s *DiskDecreaseSizeSuite) TestDiskDecreaseSize() {
	ctx := s.T().Context()

	t, err := computetest.DiskTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	diskName := utils.RandResourceName("disk")
	t.ResourceConfig = fmt.Sprintf(diskTF, diskName, "2GB")
	t.DataSourceConfig = fmt.Sprintf(diskDataSourceTF, diskName)
	t.UpdatedResourceConfig = fmt.Sprintf(diskTF, diskName, "1GB")

	tc, err := t.Build(ctx)
	s.Require().NoError(err)

	// TODO: disk should be recreated on size decrease
	tc.Steps[2].ExpectError = regexp.MustCompile(`(?i)less\s+than\s+old\s+size`)

	s.Run(tc)
}
