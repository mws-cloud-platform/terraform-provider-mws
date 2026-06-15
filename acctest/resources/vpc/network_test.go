package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

var (
	//go:embed testdata/network_defaults.tf
	networkDefaultsTF string
	//go:embed testdata/network.tf
	networkTF string
	//go:embed testdata/datasource/network.tf
	networkDataSource string
)

func TestNetworkDefaultsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NetworkDefaultsSuite))
}

func TestNetworkSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NetworkSuite))
}

func TestNetworkEnableInternetAccessSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NetworkEnableInternetAccessSuite))
}

func TestNetworkRemoveFieldsWithDefaults(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NetworkRemoveFieldsWithDefaultsSuite))
}

type NetworkDefaultsSuite struct {
	utils.ResourceSuite
}

func (s *NetworkDefaultsSuite) TestNetworkDefaults() {
	ctx := s.T().Context()

	tc, err := vpctest.NetworkTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	networkName := utils.RandResourceName("network")
	tc.ResourceConfig = fmt.Sprintf(networkDefaultsTF, networkName)
	tc.DataSourceConfig = fmt.Sprintf(networkDataSource, networkName)

	s.BuildAndRun(ctx, tc)
}

type NetworkSuite struct {
	utils.ResourceSuite
}

func (s *NetworkSuite) TestNetwork() {
	ctx := s.T().Context()

	tc, err := vpctest.NetworkTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	networkName := utils.RandResourceName("network")
	mtu := 1500
	internetAccessDisabled := false
	tc.ResourceConfig = fmt.Sprintf(networkTF,
		networkName, mtu, internetAccessDisabled,
	)
	tc.DataSourceConfig = fmt.Sprintf(networkDataSource, networkName)

	s.BuildAndRun(ctx, tc)
}

type NetworkEnableInternetAccessSuite struct {
	utils.ResourceSuite
}

func (s *NetworkEnableInternetAccessSuite) TestNetworkEnableInternetAccess() {
	ctx := s.T().Context()

	tc, err := vpctest.NetworkTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	networkName := utils.RandResourceName("network")
	mtu := 1500
	internetAccessDisabled := false
	internetAccessEnabled := true
	tc.ResourceConfig = fmt.Sprintf(networkTF,
		networkName, mtu, internetAccessDisabled,
	)
	tc.DataSourceConfig = fmt.Sprintf(networkDataSource, networkName)

	tc.UpdatedResourceConfig = fmt.Sprintf(networkTF,
		networkName, mtu, internetAccessEnabled,
	)

	s.BuildAndRun(ctx, tc)
}

type NetworkRemoveFieldsWithDefaultsSuite struct {
	utils.ResourceSuite
}

func (s *NetworkRemoveFieldsWithDefaultsSuite) TestNetworkRemoveFieldsWithDefaults() {
	ctx := s.T().Context()

	tc, err := vpctest.NetworkTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	networkName := utils.RandResourceName("network")
	mtu := 1500
	internetAccessDisabled := false
	tc.ResourceConfig = fmt.Sprintf(networkTF,
		networkName, mtu, internetAccessDisabled,
	)
	tc.UpdatedResourceConfig = fmt.Sprintf(networkDefaultsTF, networkName)
	tc.RecreateOnUpdate = true
	tc.DataSourceConfig = fmt.Sprintf(networkDataSource, networkName)

	s.BuildAndRun(ctx, tc)
}
