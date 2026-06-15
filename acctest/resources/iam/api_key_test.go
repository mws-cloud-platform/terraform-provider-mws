package iam

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	iamclient "go.mws.cloud/go-sdk/service/iam/client"
	iamsdk "go.mws.cloud/go-sdk/service/iam/sdk"

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
	utils.ResourceSuite

	serviceAccountSDK *iamsdk.ServiceAccount

	serviceAccountName string
	apiKeyName         string
}

func (s *APIKeySuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.serviceAccountName = utils.RandResourceName("sa")
	s.apiKeyName = utils.RandResourceName(s.serviceAccountName + "-api-key")

	s.serviceAccountSDK, err = iamsdk.NewServiceAccount(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.serviceAccountSDK.CreateServiceAccount(ctx, iamclient.UpsertServiceAccountRequest{
		ServiceAccount: s.serviceAccountName,
	}, iamclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("service account %q created", s.serviceAccountName)
}

func (s *APIKeySuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.serviceAccountSDK.DeleteServiceAccount(ctx, iamclient.DeleteServiceAccountRequest{
		ServiceAccount: s.serviceAccountName,
	}, iamclient.WithWait()); err != nil {
		s.T().Logf("service account %q deletion failed: %v", s.serviceAccountName, err)
	} else {
		s.T().Logf("service account %q deleted", s.serviceAccountName)
	}

	s.ResourceSuite.TearDownSuite()
}

func (s *APIKeySuite) TestApiKey() {
	ctx := s.T().Context()

	tc, err := iamtest.ApiKeyTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(apiKeyTF, s.apiKeyName, s.serviceAccountName)
	tc.DataSourceConfig = fmt.Sprintf(apiKeyDataSourceTF, s.apiKeyName, s.serviceAccountName)

	s.BuildAndRun(ctx, tc)
}
