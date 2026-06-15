package mkafka

import (
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

type MkafkaTestSuite struct {
	utils.ResourceSuite
	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet
	addressSDK *vpcsdk.Address

	networkName                        string
	network                            *vpcmodel.NetworkOptionalResponse
	subnetName                         string
	subnet                             *vpcmodel.SubnetOptionalResponse
	subnetBroker1Name                  string
	subnetBroker1                      *vpcmodel.SubnetOptionalResponse
	addressName                        string
	primaryNetworkInterfaceAddress     *vpcmodel.AddressOptionalResponse
	primaryNetworkInterfaceAddressName string
	primaryNetworkInterfaceAddressID   vpcref.AddressRef
	address1Name                       string
	address1ID                         string
	kafkaName                          string
}

func (s *MkafkaTestSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.networkName = utils.RandResourceName("network")
	s.subnetName = s.networkName + "-subnet"
	s.subnetBroker1Name = s.networkName + "-subnet-broker1"
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

	s.subnetBroker1, err = s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetBroker1Name,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("192.168.0.0/18"),
			},
		},
	}, vpcclient.WithWait())

	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetBroker1Name)

	s.subnet, err = s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("10.243.0.0/18"),
			},
		},
	}, vpcclient.WithWait())

	s.Require().NoError(err)
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
	}, vpcclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("primary network interface address %q created", s.primaryNetworkInterfaceAddressName)
	s.primaryNetworkInterfaceAddressID, err = vpcref.ParseAddressRef(ctx, s.primaryNetworkInterfaceAddress.GetMetadata().GetId().ID())
	s.Require().NoError(err)
	s.kafkaName = utils.RandResourceName("kafka")
	s.address1Name = s.kafkaName + "-broker-addr-1"

	subnetBroker1Ref, err := vpcref.ParseSubnetRef(ctx, s.subnetBroker1.GetMetadata().GetId().ID())
	s.Require().NoError(err)
	brokerAddr1, err := s.addressSDK.UpsertAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.networkName,
		Address: s.address1Name,
		Body: &vpcmodel.AddressRequest{
			Spec: vpcmodel.VpcAddressSpecRequest{
				Subnet: subnetBroker1Ref,
			},
		},
	}, vpcclient.WithWait())
	s.Require().NoError(err)
	s.address1ID = brokerAddr1.GetMetadata().GetId().ID()
	s.T().Logf("broker address 1 %q created", s.address1Name)
}

func (s *MkafkaTestSuite) TearDownSuite() {
	ctx := s.T().Context()
	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.networkName,
		Address: s.address1Name,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("address %q deletion failed: %v", s.address1Name, err)
	} else {
		s.T().Logf("address %q deleted", s.address1Name)
	}
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
		Subnet:  s.subnetBroker1Name,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"subnet %q deletion failed: %v",
			s.subnetBroker1Name,
			err,
		)
	} else {
		s.T().Logf(
			"subnet %q deleted",
			s.subnetBroker1Name,
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
