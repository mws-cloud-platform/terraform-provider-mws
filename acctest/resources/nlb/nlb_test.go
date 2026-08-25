package nlb

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	commonmodel "go.mws.cloud/go-sdk/service/common/model"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/resources/vpc"
	nlbtest "go.mws.cloud/terraform-provider-mws/service/resources/nlb/acctest"
)

var (
	//go:embed testdata/nlb.tf
	nlbTF string
	//go:embed testdata/datasource/nlb.tf
	nlbDataSourceTF string
)

func TestNLBSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NLBSuite))
}

type NLBSuite struct {
	vpc.BaseSubnetSuite

	addressSDK      *vpcsdk.Address
	addressGroupSDK *vpcsdk.AddressGroup

	addressName      string
	addressID        string
	addressGroupName string
	addressGroupID   string
}

func (s *NLBSuite) TestNLB() {
	ctx := s.T().Context()

	nlbName := s.SubnetName + "-nlb"

	tc, err := nlbtest.NlbTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(nlbTF,
		s.NetworkName,
		nlbName,
		s.addressID,
		s.addressGroupID,
	)

	tc.DataSourceConfig = fmt.Sprintf(nlbDataSourceTF, s.NetworkName, nlbName)

	s.BuildAndRun(ctx, tc)
}

func (s *NLBSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
	s.Require().NoError(err)

	s.addressGroupSDK, err = vpcsdk.NewAddressGroup(ctx, s.SDK)
	s.Require().NoError(err)

	s.addressName = s.SubnetName + "-internal-address"
	s.addressGroupName = s.SubnetName + "-address-group"

	subnetRef, err := vpcref.ParseSubnetRef(ctx, s.Subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	address, err := s.addressSDK.CreateAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.NetworkName,
		Address: s.addressName,
		Body: &vpcmodel.AddressRequest{
			Spec: vpcmodel.VpcAddressSpecRequest{
				Subnet: subnetRef,
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("address %q created", s.addressName)
	s.addressID = address.GetMetadata().GetId().ID()

	addressGroup, err := s.addressGroupSDK.CreateAddressGroup(ctx, vpcclient.UpsertAddressGroupRequest{
		Network:      s.NetworkName,
		AddressGroup: s.addressGroupName,
		Body: &vpcmodel.VpcAddressGroupRequest{
			Spec: commonmodel.VpcAddressGroupSpecRequest{
				Addresses: []commonmodel.ResourceAddressSpecOrRefRequest{
					{
						Spec: &commonmodel.ResourceAddressSpecRequest{
							Subnet: subnetRef,
						},
					},
					{
						Spec: &commonmodel.ResourceAddressSpecRequest{
							Subnet: subnetRef,
						},
					},
				},
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("address group %q created", s.addressGroupName)
	s.addressGroupID = addressGroup.GetMetadata().GetId().ID()
}

func (s *NLBSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.NetworkName,
		Address: s.addressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("address %q deletion failed: %v", s.addressName, err)
	} else {
		s.T().Logf("address %q deleted", s.addressName)
	}

	if err := s.addressGroupSDK.DeleteAddressGroup(ctx, vpcclient.DeleteAddressGroupRequest{
		Network:      s.NetworkName,
		AddressGroup: s.addressGroupName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("address group %q deletion failed: %v", s.addressGroupName, err)
	} else {
		s.T().Logf("address group %q deleted", s.addressGroupName)
	}

	s.BaseSubnetSuite.TearDownSuite()
}
