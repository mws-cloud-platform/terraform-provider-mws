package mkafka

import (
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/resources/vpc"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

type MkafkaTestSuite struct {
	vpc.BaseSubnetSuite

	addressSDK *vpcsdk.Address

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
	s.BaseSubnetSuite.SetupSuite()

	s.subnetBroker1Name = s.NetworkName + "-subnet-broker1"
	s.addressName = s.SubnetName + "-address"
	s.primaryNetworkInterfaceAddressName = s.NetworkName + "-primary-network-interface-address"

	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
	s.Require().NoError(err)

	s.subnetBroker1, err = s.SubnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.NetworkName,
		Subnet:  s.subnetBroker1Name,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("10.243.0.0/18"),
			},
		},
	}, vpcclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetBroker1Name)

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
	s.kafkaName = utils.RandResourceName("kafka")
	s.address1Name = s.kafkaName + "-broker-addr-1"

	subnetBroker1Ref, err := vpcref.ParseSubnetRef(ctx, s.subnetBroker1.GetMetadata().GetId().ID())
	s.Require().NoError(err)
	brokerAddr1, err := s.addressSDK.UpsertAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.NetworkName,
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
		Network: s.NetworkName,
		Address: s.address1Name,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("address %q deletion failed: %v", s.address1Name, err)
	} else {
		s.T().Logf("address %q deleted", s.address1Name)
	}
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

	s.BaseSubnetSuite.TearDownSuite()
}
