package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

var (
	//go:embed testdata/address_group.tf
	addressGroupTF string
	//go:embed testdata/address_group_updated.tf
	addressGroupUpdatedTF string
	//go:embed testdata/datasource/address_group.tf
	addressGroupDataSourceTF string
)

func TestAddressGroupSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AddressGroupSuite))
}

type AddressGroupSuite struct {
	BaseSubnetSuite

	addressSDK *vpcsdk.Address

	addressGroupName string
	addressName      string
	address          *vpcmodel.AddressOptionalResponse
}

func (s *AddressGroupSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.addressGroupName = utils.RandResourceName("address-group")
	s.addressName = s.addressGroupName + "-address"

	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
	s.Require().NoError(err)

	subnetRef, err := vpcref.ParseSubnetRef(ctx, s.Subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.address, err = s.addressSDK.CreateAddress(ctx, vpcclient.UpsertAddressRequest{
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
}

func (s *AddressGroupSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.NetworkName,
		Address: s.addressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("address %q deletion failed: %v", s.addressName, err)
	} else {
		s.T().Logf("address %q deleted", s.addressName)
	}

	s.BaseSubnetSuite.TearDownSuite()
}

func (s *AddressGroupSuite) TestAddressGroup() {
	ctx := s.T().Context()

	tc, err := vpctest.AddressGroupTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(addressGroupTF,
		s.addressGroupName,
		s.NetworkName,
		s.address.GetMetadata().GetId().ID(),
		s.Subnet.GetMetadata().GetId().ID(),
	)

	tc.UpdatedResourceConfig = fmt.Sprintf(addressGroupUpdatedTF,
		s.addressGroupName,
		s.NetworkName,
		s.address.GetMetadata().GetId().ID(),
	)

	tc.DataSourceConfig = fmt.Sprintf(addressGroupDataSourceTF, s.NetworkName, s.addressGroupName)

	s.BuildAndRun(ctx, tc)
}
