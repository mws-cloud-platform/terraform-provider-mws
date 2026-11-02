package compute

import (
	_ "embed"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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
	//go:embed testdata/disk.tf
	diskTF string
	//go:embed testdata/datasource/disk.tf
	diskDataSourceTF string
	//go:embed testdata/disk_with_image.tf
	diskWithImageTF string
	//go:embed testdata/disk_import_by_identity.tf
	diskImportByIdentityTF string
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

func TestDiskImportByIdentitySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiskImportByIdentitySuite))
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

type DiskImportByIdentitySuite struct {
	utils.ResourceSuite

	diskSDK  *computesdk.Disk
	diskName string
}

func (s *DiskImportByIdentitySuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.diskName = utils.RandResourceName("disk")

	s.diskSDK, err = computesdk.NewDisk(ctx, s.SDK)
	s.Require().NoError(err)

	diskType := computeref.NewMustDiskTypeRef("nbs-pl2")
	size := bytesize.MustParseString("1GB")
	_, err = s.diskSDK.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Project: s.SDK.DefaultProject(),
		Disk:    s.diskName,
		Body: computemodel.DiskRequest{
			Spec: computemodel.DiskSpecRequest{
				Zone:     s.SDK.DefaultZone(),
				Size:     &size,
				DiskType: &diskType,
				Iops:     new(computemodel.Iops(1000)),
			},
		},
	}, computeclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("disk %q created", s.diskName)
}

func (s *DiskImportByIdentitySuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.diskSDK.DeleteDisk(ctx, computeclient.DeleteDiskRequest{
		Project: s.SDK.DefaultProject(),
		Disk:    s.diskName,
	}, computeclient.WithWait()); err != nil {
		s.T().Logf("disk %q deletion failed: %v", s.diskName, err)
	} else {
		s.T().Logf("disk %q deleted", s.diskName)
	}

	s.ResourceSuite.TearDownSuite()
}

func (s *DiskImportByIdentitySuite) TestImportDiskByIdentity() {
	const resourceName = "mws_compute_disk.imported"

	s.Run(resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(diskImportByIdentityTF, s.diskName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity(resourceName, map[string]knownvalue.Check{
						"project": knownvalue.StringExact(s.SDK.DefaultProject()),
						"disk":    knownvalue.StringExact(s.diskName),
					}),
				},
			},
		},
	})
}
