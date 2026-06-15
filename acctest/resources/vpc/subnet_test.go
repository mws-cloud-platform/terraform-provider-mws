package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

var (
	//go:embed testdata/subnet.tf
	subnetTF string
	//go:embed testdata/datasource/subnet.tf
	subnetDataSourceTF string
)

func TestSubnetSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SubnetSuite))
}

type SubnetSuite struct {
	utils.ResourceSuite

	networkSDK *vpcsdk.Network

	networkName string
	network     *vpcmodel.NetworkOptionalResponse
	subnetName  string
}

func (s *SubnetSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.networkName = utils.RandResourceName("network")
	s.subnetName = s.networkName + "-subnet"

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)

	s.network, err = s.networkSDK.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: s.networkName,
	})
	s.Require().NoError(err)
	s.T().Logf("network %q created", s.networkName)
}

func (s *SubnetSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.networkSDK.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
		Network: s.networkName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("network %q deletion failed: %v", s.networkName, err)
	} else {
		s.T().Logf("network %q deleted", s.networkName)
	}

	s.ResourceSuite.TearDownSuite()
}

func (s *SubnetSuite) TestSubnet() {
	ctx := s.T().Context()

	tc, err := vpctest.SubnetTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	cidr := "192.168.0.0/16"
	tc.ResourceConfig = fmt.Sprintf(subnetTF,
		s.networkName, s.subnetName, cidr,
	)

	tc.DataSourceConfig = fmt.Sprintf(subnetDataSourceTF, s.networkName, s.subnetName)
	s.BuildAndRun(ctx, tc)
}
