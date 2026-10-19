package mk8s

import (
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/resources/vpc"
)

type Mk8sTestSuite struct {
	vpc.BaseSubnetSuite

	addressSDK *vpcsdk.Address

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
	s.BaseSubnetSuite.SetupSuite()

	s.subnetServiceName = s.NetworkName + "-subnet-service"
	s.subnetPodName = s.NetworkName + "-subnet-pod"
	s.addressName = s.SubnetName + "-address"
	s.primaryNetworkInterfaceAddressName = s.NetworkName + "-primary-network-interface-address"

	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
	s.Require().NoError(err)

	s.subnetService, err = s.SubnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.NetworkName,
		Subnet:  s.subnetServiceName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("10.243.0.0/18"),
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetServiceName)

	s.subnetPod, err = s.SubnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.NetworkName,
		Subnet:  s.subnetPodName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("10.244.0.0/18"),
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetPodName)

	subnetRef, err := vpcref.ParseSubnetRef(s.T().Context(), s.Subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)
	s.primaryNetworkInterfaceAddress, err = s.addressSDK.CreateAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.NetworkName,
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
		Network: s.NetworkName,
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

	if err := s.SubnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.NetworkName,
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
	if err := s.SubnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.NetworkName,
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

	s.BaseSubnetSuite.TearDownSuite()
}
