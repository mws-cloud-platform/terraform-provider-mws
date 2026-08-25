package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

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
	BaseNetworkSuite
}

func (s *SubnetSuite) TestSubnet() {
	ctx := s.T().Context()

	subnetName := s.NetworkName + "-subnet"

	tc, err := vpctest.SubnetTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	cidr := "192.168.0.0/16"
	tc.ResourceConfig = fmt.Sprintf(subnetTF,
		s.NetworkName, subnetName, cidr,
	)

	tc.DataSourceConfig = fmt.Sprintf(subnetDataSourceTF, s.NetworkName, subnetName)
	s.BuildAndRun(ctx, tc)
}
