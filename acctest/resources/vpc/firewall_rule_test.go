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
	BaseNetworkSuite
}

func (s *FirewallRuleAllowSSHSuite) TestFirewallRuleAllowSSH() {
	ctx := s.T().Context()

	firewallRuleName := utils.RandResourceName("firewall-rule")

	tc, err := vpctest.FirewallRuleTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(firewallRuleTF,
		s.NetworkName, firewallRuleName,
	)

	tc.DataSourceConfig = fmt.Sprintf(firewallRuleDataSourceTF, s.NetworkName, firewallRuleName)
	s.BuildAndRun(ctx, tc)
}
