package compute

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	computeclient "go.mws.cloud/go-sdk/service/compute/client"
	computemodel "go.mws.cloud/go-sdk/service/compute/model"
	computesdk "go.mws.cloud/go-sdk/service/compute/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	computetest "go.mws.cloud/terraform-provider-mws/service/resources/compute/acctest"
)

var (
	//go:embed testdata/virtual_machine.tf
	virtualMachineTF string
	//go:embed testdata/datasource/virtual_machine.tf
	virtualMachineDataSourceTF string
)

func TestVirtualMachineSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(VirtualMachineSuite))
}

type VirtualMachineSuite struct {
	utils.ResourceSuite

	virtualMachineName string

	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet
	addressSDK *vpcsdk.Address
	imageSDK   *computesdk.Image
	diskSDK    *computesdk.Disk

	networkName                        string
	network                            *vpcmodel.NetworkOptionalResponse
	subnetName                         string
	subnet                             *vpcmodel.SubnetOptionalResponse
	primaryNetworkInterfaceAddressName string
	primaryNetworkInterfaceAddress     *vpcmodel.AddressOptionalResponse
	bootDiskName                       string
	bootDisk                           *computemodel.DiskOptionalResponse
	dataDiskName                       string
	dataDisk                           *computemodel.DiskOptionalResponse
}

func (s *VirtualMachineSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.virtualMachineName = utils.RandResourceName("vm")
	s.networkName = s.virtualMachineName + "-network"
	s.subnetName = s.virtualMachineName + "-subnet"
	s.primaryNetworkInterfaceAddressName = s.virtualMachineName + "-primary-network-interface-address"
	s.bootDiskName = s.virtualMachineName + "-boot-disk"
	s.dataDiskName = s.virtualMachineName + "-data-disk"

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)
	s.subnetSDK, err = vpcsdk.NewSubnet(ctx, s.SDK)
	s.Require().NoError(err)
	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
	s.Require().NoError(err)
	s.imageSDK, err = computesdk.NewImage(ctx, s.SDK)
	s.Require().NoError(err)
	s.diskSDK, err = computesdk.NewDisk(ctx, s.SDK)
	s.Require().NoError(err)

	s.network, err = s.networkSDK.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: s.networkName,
	})
	s.Require().NoError(err)
	s.T().Logf("network %q created", s.networkName)

	s.subnet, err = s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("192.168.0.0/16"),
			},
		},
	})
	s.Require().NoError(err)
	s.Require().NotNil(s.subnet.GetMetadata().GetId())
	s.T().Logf("subnet %q created", s.subnetName)

	subnetRef, err := vpcref.ParseSubnetRef(s.T().Context(), s.subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.primaryNetworkInterfaceAddress, err = s.addressSDK.CreateAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.networkName,
		Address: s.primaryNetworkInterfaceAddressName,
		Body: &vpcmodel.AddressRequest{
			Spec: vpcmodel.VpcAddressSpecRequest{
				Subnet: subnetRef,
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("primary network interface address %q created", s.primaryNetworkInterfaceAddressName)

	image, err := s.imageSDK.GetImage(ctx, computeclient.GetImageRequest{
		Project: "mws-ubuntu",
		Image:   "mws-ubuntu-2204-lts-v20250529",
	})
	s.Require().NoError(err)
	s.Require().NotNil(image.GetMetadata().GetId())

	imageRef, err := computeref.ParseImageRef(s.T().Context(), image.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.bootDisk, err = s.diskSDK.CreateDisk(ctx, computeclient.UpsertDiskRequest{
		Disk: s.bootDiskName,
		Body: computemodel.DiskRequest{
			Spec: s.diskSpec(&computemodel.DiskSpecSourceRequest{
				Image: new(imageRef),
			}),
		},
	})
	s.Require().NoError(err)
	s.T().Logf("disk %q created", s.bootDiskName)

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

	for _, name := range []string{s.bootDiskName, s.dataDiskName} {
		err := s.diskSDK.DeleteDisk(ctx, computeclient.DeleteDiskRequest{
			Disk: name,
		}, computeclient.WithWait())
		if err != nil {
			s.T().Logf("disk %q deletion failed: %v", name, err)
		} else {
			s.T().Logf("disk %q deleted", name)
		}
	}

	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.networkName,
		Address: s.primaryNetworkInterfaceAddressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("primary network interface address %q deletion failed: %v", s.primaryNetworkInterfaceAddressName, err)
	} else {
		s.T().Logf("primary network interface address %q deleted", s.primaryNetworkInterfaceAddressName)
	}

	if err := s.subnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("subnet %q deletion failed: %v", s.subnetName, err)
	} else {
		s.T().Logf("subnet %q deleted", s.subnetName)
	}

	if err := s.networkSDK.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
		Network: s.networkName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("network %q deletion failed: %v", s.networkName, err)
	} else {
		s.T().Logf("network %q deleted", s.networkName)
	}

	s.ResourceSuite.TearDownSuite()
}

func (s *VirtualMachineSuite) TestVirtualMachine() {
	ctx := s.T().Context()

	tc, err := computetest.VirtualMachineTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(virtualMachineTF,
		s.virtualMachineName,
		s.bootDisk.GetMetadata().GetId().ID(),
		s.dataDisk.GetMetadata().GetId().ID(),
		s.primaryNetworkInterfaceAddress.GetMetadata().GetId().ID(),
	)

	tc.DataSourceConfig = fmt.Sprintf(virtualMachineDataSourceTF, s.virtualMachineName)
	s.BuildAndRun(ctx, tc)
}

func (s *VirtualMachineSuite) diskSpec(source *computemodel.DiskSpecSourceRequest) computemodel.DiskSpecRequest {
	return computemodel.DiskSpecRequest{
		Zone:      "ru-central1-a",
		DiskType:  new(computeref.NewDiskTypeRef("nbs-pl2")),
		Size:      new(bytesize.MustNewFromInt64(10, bytesize.GB)),
		Iops:      new(computemodel.Iops(1000)),
		Source:    source,
		BlockSize: new(bytesize.MustNewFromInt64(4, bytesize.KB)),
	}
}
