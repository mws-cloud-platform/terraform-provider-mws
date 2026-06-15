package secretmanager

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	kmsclient "go.mws.cloud/go-sdk/service/kms/client"
	kmsmodel "go.mws.cloud/go-sdk/service/kms/model"
	kmssdk "go.mws.cloud/go-sdk/service/kms/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	secretmanagertest "go.mws.cloud/terraform-provider-mws/service/resources/secretmanager/acctest"
)

var (
	//go:embed testdata/secret.tf
	secretTF string
	//go:embed testdata/secret_with_encryption.tf
	secretWithEncryptionTF string
	//go:embed testdata/datasource/secret.tf
	secretDataSourceTF string
)

func TestSecretSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SecretSuite))
}

func TestSecretWithEncryptionSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SecretWithEncryptionSuite))
}

type SecretSuite struct {
	utils.ResourceSuite
}

func (s *SecretSuite) TestSecret() {
	ctx := s.T().Context()

	tc, err := secretmanagertest.SecretTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	secretName := utils.RandResourceName("secret")

	tc.ResourceConfig = fmt.Sprintf(secretTF, secretName)
	tc.DataSourceConfig = fmt.Sprintf(secretDataSourceTF, secretName)
	s.BuildAndRun(ctx, tc)
}

type SecretWithEncryptionSuite struct {
	utils.ResourceSuite

	cryptoKeySDK *kmssdk.CryptoKey

	cryptoKeyName string
	cryptoKey     *kmsmodel.CryptoKeyOptionalResponse
}

func (s *SecretWithEncryptionSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.cryptoKeyName = utils.RandResourceName("crypto-key")

	s.cryptoKeySDK, err = kmssdk.NewCryptoKey(ctx, s.SDK)
	s.Require().NoError(err)

	s.cryptoKey, err = s.cryptoKeySDK.UpsertCryptoKey(ctx, kmsclient.UpsertCryptoKeyRequest{
		Key: s.cryptoKeyName,
		Body: kmsmodel.CryptoKeyRequest{
			Spec: kmsmodel.CryptoKeySpecRequest{
				DefaultAlgorithm: new(kmsmodel.CryptoKeyAlgorithm_AES_256_GCM),
				DestructionPolicy: &kmsmodel.CryptoKeySpecDestructionPolicyRequest{
					DefaultDestructionIntervalDays: new(int32(1)),
				},
				RotationPolicy: &kmsmodel.CryptoKeySpecRotationPolicyRequest{
					Enabled: new(false),
				},
				UsagePolicy: &kmsmodel.CryptoKeySpecUsagePolicyRequest{
					Enabled: new(true),
				},
			},
		},
	}, kmsclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("crypto key %q created", s.cryptoKeyName)
}

func (s *SecretWithEncryptionSuite) TearDownSuite() {
	if _, err := s.cryptoKeySDK.ScheduleDestructionOfCryptoKey(s.T().Context(), kmsclient.ScheduleDestructionOfCryptoKeyRequest{
		Key: s.cryptoKeyName,
	}); err != nil {
		s.T().Logf("failed to schedule crypto key %q destruction: %v", s.cryptoKeyName, err)
	}

	s.ResourceSuite.TearDownSuite()
}

func (s *SecretWithEncryptionSuite) TestSecretWithEncryption() {
	ctx := s.T().Context()

	tc, err := secretmanagertest.SecretTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	secretName := utils.RandResourceName("secret")

	tc.ResourceConfig = fmt.Sprintf(secretWithEncryptionTF, secretName, s.cryptoKey.GetMetadata().GetId().ID())
	tc.DataSourceConfig = fmt.Sprintf(secretDataSourceTF, secretName)
	s.BuildAndRun(ctx, tc)
}
