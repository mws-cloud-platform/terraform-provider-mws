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
	utils.ResourceSuite

	serviceAccountSDK *iamsdk.ServiceAccount

	serviceAccountName string
	hmacKeyName        string
}

func (s *hmacKeySuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.serviceAccountName = utils.RandResourceName("sa")
	s.hmacKeyName = utils.RandResourceName(s.serviceAccountName + "-hmac-key")

	s.serviceAccountSDK, err = iamsdk.NewServiceAccount(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.serviceAccountSDK.CreateServiceAccount(ctx, iamclient.UpsertServiceAccountRequest{
		ServiceAccount: s.serviceAccountName,
	}, iamclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("service account %q created", s.serviceAccountName)
}

func (s *hmacKeySuite) TearDownSuite() {
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

func (s *hmacKeySuite) TestHmacKey() {
	ctx := s.T().Context()

	tc, err := iamtest.HmacKeyTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(hmacKeyTF, s.hmacKeyName, s.serviceAccountName)
	tc.DataSourceConfig = fmt.Sprintf(hmacKeyDataSourceTF, s.hmacKeyName, s.serviceAccountName)

	s.BuildAndRun(ctx, tc)
}
