package compute

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computemodel "go.mws.cloud/go-sdk/service/compute/model"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	iamclient "go.mws.cloud/go-sdk/service/iam/client"
	iammodel "go.mws.cloud/go-sdk/service/iam/model"
	iamsdk "go.mws.cloud/go-sdk/service/iam/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"

	"go.mws.cloud/terraform-provider-mws/acctest/resources/vpc"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	computetest "go.mws.cloud/terraform-provider-mws/service/resources/compute/acctest"
)

var (
	//go:embed testdata/virtual_machine.tf
	virtualMachineTF string
	//go:embed testdata/virtual_machine_attach_disk.tf
	virtualMachineAttachDiskTF string
	//go:embed testdata/virtual_machine_with_sa.tf
	virtualMachineWithSATF string
	//go:embed testdata/datasource/virtual_machine.tf
	virtualMachineDataSourceTF string
)

func TestVirtualMachineSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(VirtualMachineSuite))
}

func TestVirtualMachineWithSASuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(VirtualMachineWithSASuite))
}

type VirtualMachineSuite struct {
	vpc.BaseSubnetSuite

	virtualMachineName string

	diskSDK *computesdk.Disk

	dataDiskName string
	dataDisk     *computemodel.DiskOptionalResponse
}

func (s *VirtualMachineSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.virtualMachineName = utils.RandResourceName("vm")
	s.dataDiskName = s.virtualMachineName + "-data-disk"

	s.diskSDK, err = computesdk.NewDisk(ctx, s.SDK)
	s.Require().NoError(err)

	s.dataDisk, err = s.diskSDK.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Disk: s.dataDiskName,
		Body: computemodel.DiskRequest{
			Spec: s.diskSpec(nil),
		},
	})
	s.Require().NoError(err)
	s.T().Logf("disk %q created", s.dataDiskName)
}

func (s *VirtualMachineSuite) TearDownSuite() {
	ctx := s.T().Context()

	err := s.diskSDK.DeleteDisk(ctx, computeclient.DeleteDiskRequest{
		Disk: s.dataDiskName,
	}, computeclient.WithWait())
	if err != nil {
		s.T().Logf("disk %q deletion failed: %v", s.dataDiskName, err)
	} else {
		s.T().Logf("disk %q deleted", s.dataDiskName)
	}

	s.BaseSubnetSuite.TearDownSuite()
}

func (s *VirtualMachineSuite) TestVirtualMachine() {
	ctx := s.T().Context()

	tc, err := computetest.VirtualMachineTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	image, err := getBaseImage(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(virtualMachineTF,
		s.virtualMachineName,
		image.GetMetadata().GetId().ID(),
		s.dataDisk.GetMetadata().GetId().ID(),
		s.Subnet.GetMetadata().GetId().ID(),
	)
	tc.UpdatedResourceConfig = fmt.Sprintf(virtualMachineAttachDiskTF,
		s.virtualMachineName,
		image.GetMetadata().GetId().ID(),
		s.dataDisk.GetMetadata().GetId().ID(),
		s.Subnet.GetMetadata().GetId().ID(),
	)

	tc.DataSourceConfig = fmt.Sprintf(virtualMachineDataSourceTF, s.virtualMachineName)
	s.BuildAndRun(ctx, tc)
}

func (s *VirtualMachineSuite) diskSpec(source *computemodel.DiskSpecSourceRequest) computemodel.DiskSpecRequest {
	return computemodel.DiskSpecRequest{
		Zone:      "ru-central1-a",
		DiskType:  new(computeref.NewMustDiskTypeRef("nbs-pl2")),
		Size:      new(bytesize.MustNewFromInt64(10, bytesize.GB)),
		Iops:      new(computemodel.Iops(1000)),
		Source:    source,
		BlockSize: new(bytesize.MustNewFromInt64(4, bytesize.KB)),
	}
}

type VirtualMachineWithSASuite struct {
	vpc.BaseSubnetSuite

	virtualMachineName string

	serviceAccountSDK *iamsdk.ServiceAccount

	serviceAccountName string
	serviceAccount     *iammodel.ServiceAccountResponse
}

func (s *VirtualMachineWithSASuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.virtualMachineName = utils.RandResourceName("vm")
	s.serviceAccountName = s.virtualMachineName + "-service-account"

	s.serviceAccountSDK, err = iamsdk.NewServiceAccount(ctx, s.SDK)
	s.Require().NoError(err)

	s.serviceAccount, err = s.serviceAccountSDK.CreateServiceAccount(ctx, iamclient.UpsertServiceAccountRequest{
		ServiceAccount: s.serviceAccountName,
		Body: iammodel.ServiceAccountRequest{
			Spec: iammodel.ServiceAccountSpecRequest{},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("service account %q created", s.serviceAccountName)
}

func (s *VirtualMachineWithSASuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.serviceAccountSDK.DeleteServiceAccount(ctx, iamclient.DeleteServiceAccountRequest{
		ServiceAccount: s.serviceAccountName,
	}, iamclient.WithWait()); err != nil {
		s.T().Logf("service account %q deletion failed: %v", s.serviceAccountName, err)
	} else {
		s.T().Logf("service account %q deleted", s.serviceAccountName)
	}

	s.BaseSubnetSuite.TearDownSuite()
}

func (s *VirtualMachineWithSASuite) TestVirtualMachine() {
	ctx := s.T().Context()

	tc, err := computetest.VirtualMachineTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	image, err := getBaseImage(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(virtualMachineWithSATF,
		s.virtualMachineName,
		image.GetMetadata().GetId().ID(),
		s.Subnet.GetMetadata().GetId().ID(),
		strconv.Quote(s.serviceAccount.GetMetadata().GetId().ID()),
	)
	tc.UpdatedResourceConfig = fmt.Sprintf(virtualMachineWithSATF,
		s.virtualMachineName,
		image.GetMetadata().GetId().ID(),
		s.Subnet.GetMetadata().GetId().ID(),
		"null",
	)

	s.BuildAndRun(ctx, tc)
}

func getBaseImage(ctx context.Context, sdk *mws.SDK) (*computemodel.ImageOptionalResponse, error) {
	imageSDK, err := computesdk.NewImage(ctx, sdk)
	if err != nil {
		return nil, err
	}

	return imageSDK.GetImage(ctx, computeclient.GetImageRequest{
		Project: "mws-ubuntu",
		Image:   "mws-ubuntu-2204-lts-v20250529",
	})
}
