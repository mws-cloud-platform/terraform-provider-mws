package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

var (
	//go:embed testdata/firewall_rule.tf
	firewallRuleTF string
	//go:embed testdata/datasource/firewall_rule.tf
	firewallRuleDataSourceTF string
)

func TestFirewallRuleAllowSSHSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FirewallRuleAllowSSHSuite))
}

type FirewallRuleAllowSSHSuite struct {
	utils.ResourceSuite

	networkSDK *vpcsdk.Network

	networkName      string
	firewallRuleName string
}

func (s *FirewallRuleAllowSSHSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.networkName = utils.RandResourceName("network")
	s.firewallRuleName = utils.RandResourceName("firewall-rule")

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.networkSDK.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: s.networkName,
	})
	s.Require().NoError(err)
	s.T().Logf("network %q created", s.networkName)
}

func (s *FirewallRuleAllowSSHSuite) TearDownSuite() {
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

func (s *FirewallRuleAllowSSHSuite) TestFirewallRuleAllowSSH() {
	ctx := s.T().Context()

	tc, err := vpctest.FirewallRuleTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(firewallRuleTF,
		s.networkName, s.firewallRuleName,
	)

	tc.DataSourceConfig = fmt.Sprintf(firewallRuleDataSourceTF, s.networkName, s.firewallRuleName)
	s.BuildAndRun(ctx, tc)
}
