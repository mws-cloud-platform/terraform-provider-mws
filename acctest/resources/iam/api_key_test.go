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
	//go:embed testdata/api_key.tf
	apiKeyTF string
	//go:embed testdata/datasource/api_key.tf
	apiKeyDataSourceTF string
)

func TestApiKeySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(APIKeySuite))
}

type APIKeySuite struct {
	BaseServiceAccountSuite
}

func (s *APIKeySuite) TestApiKey() {
	ctx := s.T().Context()

	apiKeyName := utils.RandResourceName(s.ServiceAccountName + "-api-key")

	tc, err := iamtest.ApiKeyTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(apiKeyTF, apiKeyName, s.ServiceAccountName)
	tc.DataSourceConfig = fmt.Sprintf(apiKeyDataSourceTF, apiKeyName, s.ServiceAccountName)

	s.BuildAndRun(ctx, tc)
}
