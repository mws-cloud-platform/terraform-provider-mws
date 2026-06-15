package mk8s

import (
	_ "embed"

	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

type Mk8sTestSuite struct {
	utils.ResourceSuite
	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet
	addressSDK *vpcsdk.Address

	networkName                        string
	network                            *vpcmodel.NetworkOptionalResponse
	subnetName                         string
	subnet                             *vpcmodel.SubnetOptionalResponse
	subnetServiceName                  string
	subnetService                      *vpcmodel.SubnetOptionalResponse
	subnetPodName                      string
	subnetPod                          *vpcmodel.SubnetOptionalResponse
	addressName                        string
	primaryNetworkInterfaceAddress     *vpcmodel.AddressOptionalResponse
	primaryNetworkInterfaceAddressName string
	primaryNetworkInterfaceAddressID   vpcref.AddressRef
}

func (s *Mk8sTestSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.networkName = utils.RandResourceName("network")
	s.subnetName = s.networkName + "-subnet"
	s.subnetServiceName = s.networkName + "-subnet-service"
	s.subnetPodName = s.networkName + "-subnet-pod"
	s.addressName = s.subnetName + "-address"
	s.primaryNetworkInterfaceAddressName = s.networkName + "-primary-network-interface-address"

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)
	s.subnetSDK, err = vpcsdk.NewSubnet(ctx, s.SDK)
	s.Require().NoError(err)
	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
	s.Require().NoError(err)

	s.network, err = s.networkSDK.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: s.networkName,
	})
	s.Require().NoError(err)
	s.T().Logf("network %q created", s.networkName)

	s.subnetService, err = s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetServiceName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("192.168.0.0/18"),
			},
		},
	})

	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetServiceName)
	s.subnetPod, err = s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetPodName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("10.244.0.0/18"),
			},
		},
	})

	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetPodName)
	s.subnet, err = s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("10.243.0.0/18"),
			},
		},
	})

	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetPodName)
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
	}, vpcclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("primary network interface address %q created", s.primaryNetworkInterfaceAddressName)
	s.primaryNetworkInterfaceAddressID, err = vpcref.ParseAddressRef(ctx, s.primaryNetworkInterfaceAddress.GetMetadata().GetId().ID())
	s.Require().NoError(err)
}

func (s *Mk8sTestSuite) TearDownSuite() {
	ctx := s.T().Context()
	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.networkName,
		Address: s.primaryNetworkInterfaceAddressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"address %q deletion failed: %v",
			s.primaryNetworkInterfaceAddressName,
			err,
		)
	} else {
		s.T().Logf(
			"address %q deleted",
			s.primaryNetworkInterfaceAddressName,
		)
	}

	if err := s.subnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"subnet %q deletion failed: %v",
			s.subnetName,
			err,
		)
	} else {
		s.T().Logf(
			"subnet %q deleted",
			s.subnetName,
		)
	}
	if err := s.subnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetServiceName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"subnet %q deletion failed: %v",
			s.subnetServiceName,
			err,
		)
	} else {
		s.T().Logf(
			"subnet %q deleted",
			s.subnetServiceName,
		)
	}
	if err := s.subnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetPodName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"subnet %q deletion failed: %v",
			s.subnetPodName,
			err,
		)
	} else {
		s.T().Logf(
			"subnet %q deleted",
			s.subnetPodName,
		)
	}

	if err := s.networkSDK.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
		Network: s.networkName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"network %q deletion failed: %v",
			s.networkName,
			err,
		)
	} else {
		s.T().Logf(
			"network %q deleted",
			s.networkName,
		)
	}

	s.ResourceSuite.TearDownSuite()
}
