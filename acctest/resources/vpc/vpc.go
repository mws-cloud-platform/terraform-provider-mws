package vpc

import (
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

type BaseNetworkSuite struct {
	utils.ResourceSuite

	NetworkSDK *vpcsdk.Network

	NetworkName string
	Network     *vpcmodel.NetworkOptionalResponse
}

func (s *BaseNetworkSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.NetworkName = utils.RandResourceName("network")

	s.NetworkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)

	s.Network, err = s.NetworkSDK.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: s.NetworkName,
	})
	s.Require().NoError(err)
	s.T().Logf("network %q created", s.NetworkName)
}

func (s *BaseNetworkSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.NetworkSDK.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
		Network: s.NetworkName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("network %q deletion failed: %v", s.NetworkName, err)
	} else {
		s.T().Logf("network %q deleted", s.NetworkName)
	}

	s.ResourceSuite.TearDownSuite()
}

type BaseSubnetSuite struct {
	BaseNetworkSuite

	SubnetSDK *vpcsdk.Subnet

	SubnetName string
	Subnet     *vpcmodel.SubnetOptionalResponse
}

func (s *BaseSubnetSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseNetworkSuite.SetupSuite()

	s.SubnetName = s.NetworkName + "-subnet"

	s.SubnetSDK, err = vpcsdk.NewSubnet(ctx, s.SDK)
	s.Require().NoError(err)

	s.Subnet, err = s.SubnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.NetworkName,
		Subnet:  s.SubnetName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("192.168.0.0/16"),
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.SubnetName)
}

func (s *BaseSubnetSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.SubnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.NetworkName,
		Subnet:  s.SubnetName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("subnet %q deletion failed: %v", s.SubnetName, err)
	} else {
		s.T().Logf("subnet %q deleted", s.SubnetName)
	}

	s.BaseNetworkSuite.TearDownSuite()
}
