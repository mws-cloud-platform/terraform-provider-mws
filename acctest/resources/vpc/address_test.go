package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

var (
	//go:embed testdata/address.tf
	addressTF string
	//go:embed testdata/datasource/address.tf
	addressDataSourceTF string
)

func TestAddressSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AddressSuite))
}

type AddressSuite struct {
	utils.ResourceSuite

	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet

	networkName string
	network     *vpcmodel.NetworkOptionalResponse
	subnetName  string
	subnet      *vpcmodel.SubnetOptionalResponse
	addressName string
}

func (s *AddressSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.networkName = utils.RandResourceName("network")
	s.subnetName = s.networkName + "-subnet"
	s.addressName = s.subnetName + "-address"

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)
	s.subnetSDK, err = vpcsdk.NewSubnet(ctx, s.SDK)
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
	s.T().Logf("subnet %q created", s.subnetName)
}

func (s *AddressSuite) TearDownSuite() {
	ctx := s.T().Context()

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

func (s *AddressSuite) TestAddress() {
	ctx := s.T().Context()

	tc, err := vpctest.AddressTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(addressTF,
		s.networkName,
		s.subnet.GetMetadata().GetId().ID(),
		s.addressName,
	)

	tc.DataSourceConfig = fmt.Sprintf(addressDataSourceTF, s.networkName, s.addressName)

	s.BuildAndRun(ctx, tc)
}
