package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

var (
	//go:embed testdata/route.tf
	routeTF string
	//go:embed testdata/datasource/route.tf
	routeDataSourceTF string
)

func TestRouteSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RouteSuite))
}

type RouteSuite struct {
	utils.ResourceSuite

	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet
	addressSDK *vpcsdk.Address

	networkName string
	network     *vpcmodel.NetworkOptionalResponse
	subnetName  string
	subnet      *vpcmodel.SubnetOptionalResponse
	addressName string
	address     *vpcmodel.AddressOptionalResponse
	routeName   string
}

func (s *RouteSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.routeName = utils.RandResourceName("route")
	s.subnetName = s.routeName + "-subnet"
	s.addressName = s.routeName + "-address"
	s.networkName = s.routeName + "-net"

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

	subnetRef, err := vpcref.ParseSubnetRef(s.T().Context(), s.subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetName)
	s.address, err = s.addressSDK.CreateAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.networkName,
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

func (s *RouteSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.networkName,
		Address: s.addressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("address %q deletion failed: %v", s.addressName, err)
	} else {
		s.T().Logf("address %q deleted", s.addressName)
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

func (s *RouteSuite) TestRoute() {
	ctx := s.T().Context()

	tc, err := vpctest.RouteTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(routeTF,
		s.routeName,
		s.networkName,
		s.address.GetMetadata().GetId().ID(),
	)
	s.T().Log(tc.ResourceConfig)

	tc.DataSourceConfig = fmt.Sprintf(routeDataSourceTF, s.routeName, s.networkName)

	s.BuildAndRun(ctx, tc)
}
