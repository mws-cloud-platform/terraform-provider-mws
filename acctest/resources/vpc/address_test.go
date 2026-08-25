package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

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
	BaseSubnetSuite
}

func (s *AddressSuite) TestAddress() {
	ctx := s.T().Context()

	addressName := s.SubnetName + "-address"

	tc, err := vpctest.AddressTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(addressTF,
		s.NetworkName,
		s.Subnet.GetMetadata().GetId().ID(),
		addressName,
	)

	tc.DataSourceConfig = fmt.Sprintf(addressDataSourceTF, s.NetworkName, addressName)

	s.BuildAndRun(ctx, tc)
}
