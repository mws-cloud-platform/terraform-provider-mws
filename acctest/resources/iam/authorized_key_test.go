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
	//go:embed testdata/authorized_key.tf
	authorizedKeyTF string
	//go:embed testdata/datasource/authorized_key.tf
	authorizedKeyDatasourceTF string
)

func TestAuthorizedKeySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AuthorizedKeySuite))
}

type AuthorizedKeySuite struct {
	baseServiceAccountSuite
}

func (s *AuthorizedKeySuite) TestAuthorizedKey() {
	ctx := s.T().Context()

	authorizedKeyName := utils.RandResourceName(s.serviceAccountName + "-auth-key")

	tc, err := iamtest.AuthorizedKeyTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(authorizedKeyTF, authorizedKeyName, s.serviceAccountName)
	tc.DataSourceConfig = fmt.Sprintf(authorizedKeyDatasourceTF, authorizedKeyName, s.serviceAccountName)

	s.BuildAndRun(ctx, tc)
}
