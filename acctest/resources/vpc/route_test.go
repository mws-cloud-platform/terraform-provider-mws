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
	BaseSubnetSuite

	addressSDK *vpcsdk.Address

	routeName   string
	addressName string
	address     *vpcmodel.AddressOptionalResponse
}

func (s *RouteSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.routeName = utils.RandResourceName("route")
	s.addressName = s.routeName + "-address"

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

func (s *RouteSuite) TearDownSuite() {
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

func (s *RouteSuite) TestRoute() {
	ctx := s.T().Context()

	tc, err := vpctest.RouteTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(routeTF,
		s.routeName,
		s.NetworkName,
		s.address.GetMetadata().GetId().ID(),
	)
	s.T().Log(tc.ResourceConfig)

	tc.DataSourceConfig = fmt.Sprintf(routeDataSourceTF, s.routeName, s.NetworkName)

	s.BuildAndRun(ctx, tc)
}
