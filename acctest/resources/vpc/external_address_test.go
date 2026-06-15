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
	//go:embed testdata/external_address.tf
	externalAddressTF string
	//go:embed testdata/datasource/external_address.tf
	externalAddressDataSourceTF string
)

func TestExternalAddressSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExternalAddressSuite))
}

type ExternalAddressSuite struct {
	utils.ResourceSuite
}

func (s *ExternalAddressSuite) TestExternalAddress() {
	ctx := s.T().Context()

	tc, err := vpctest.ExternalAddressTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	externalAddressName := utils.RandResourceName("external-address")
	tc.ResourceConfig = fmt.Sprintf(externalAddressTF, externalAddressName)
	tc.DataSourceConfig = fmt.Sprintf(externalAddressDataSourceTF, externalAddressName)

	s.BuildAndRun(ctx, tc)
}
