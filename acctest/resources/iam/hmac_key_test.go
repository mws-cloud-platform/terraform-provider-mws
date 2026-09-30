package iam

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	iamtest "go.mws.cloud/terraform-provider-mws/service/resources/iam/acctest"
)

var (
	//go:embed testdata/hmac_key.tf
	hmacKeyTF string
	//go:embed testdata/datasource/hmac_key.tf
	hmacKeyDataSourceTF string
)

func TestHmacKeySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(hmacKeySuite))
}

type hmacKeySuite struct {
	BaseServiceAccountSuite
}

func (s *hmacKeySuite) TestHmacKey() {
	ctx := s.T().Context()

	hmacKeyName := utils.RandResourceName(s.ServiceAccountName + "-hmac-key")

	tc, err := iamtest.HmacKeyTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(hmacKeyTF, hmacKeyName, s.ServiceAccountName)
	tc.DataSourceConfig = fmt.Sprintf(hmacKeyDataSourceTF, hmacKeyName, s.ServiceAccountName)

	s.BuildAndRun(ctx, tc)
}
