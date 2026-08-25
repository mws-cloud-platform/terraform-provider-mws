package secretmanager

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	secretmanagerclient "go.mws.cloud/go-sdk/service/secretmanager/client"
	secretmanagermodel "go.mws.cloud/go-sdk/service/secretmanager/model"
	secretmanagersdk "go.mws.cloud/go-sdk/service/secretmanager/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	secretmanagertest "go.mws.cloud/terraform-provider-mws/service/resources/secretmanager/acctest"
)

var (
	//go:embed testdata/secret_version.tf
	secretVersionTF string
	//go:embed testdata/datasource/secret_version.tf
	secretVersionDataSourceTF string
)

func TestSecretVersionSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SecretVersionSuite))
}

type SecretVersionSuite struct {
	utils.ResourceSuite

	secretSDK *secretmanagersdk.Secret

	secretName string
	secret     *secretmanagermodel.SecretOptionalResponse
}

func (s *SecretVersionSuite) TestSecretVersion() {
	ctx := s.T().Context()

	tc, err := secretmanagertest.SecretVersionTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(secretVersionTF, s.secretName)
	tc.DataSourceConfig = fmt.Sprintf(secretVersionDataSourceTF, s.secretName)
	s.BuildAndRun(ctx, tc)
}

func (s *SecretVersionSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.secretName = utils.RandResourceName("secret")

	s.secretSDK, err = secretmanagersdk.NewSecret(ctx, s.SDK)
	s.Require().NoError(err)

	active := true
	s.secret, err = s.secretSDK.UpsertSecret(ctx, secretmanagerclient.UpsertSecretRequest{
		Name: s.secretName,
		Body: secretmanagermodel.SecretRequest{
			Spec: secretmanagermodel.SecretSpecRequest{
				Active: &active,
			},
		},
	}, secretmanagerclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("secret %q created", s.secretName)
}

func (s *SecretVersionSuite) TearDownSuite() {
	if err := s.secretSDK.DeleteSecret(s.T().Context(), secretmanagerclient.DeleteSecretRequest{
		Name: s.secretName,
	}, secretmanagerclient.WithWait()); err != nil {
		s.T().Logf("failed to delete secret %q: %v", s.secretName, err)
	}

	s.ResourceSuite.TearDownSuite()
}
