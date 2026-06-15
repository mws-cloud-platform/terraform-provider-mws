package kms

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	kmstest "go.mws.cloud/terraform-provider-mws/service/resources/kms/acctest"
)

var (
	//go:embed testdata/crypto_key.tf
	cryptoKeyTF string
	//go:embed testdata/datasource/crypto_key.tf
	cryptoKeyDataSourceTF string
)

func TestCryptoKeySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CryptoKeySuite))
}

type CryptoKeySuite struct {
	utils.ResourceSuite
}

func (s *CryptoKeySuite) TestCryptoKey() {
	ctx := s.T().Context()

	tc, err := kmstest.CryptoKeyTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	cryptoKeyName := utils.RandResourceName("crypto-key")
	tc.ResourceConfig = fmt.Sprintf(cryptoKeyTF, cryptoKeyName)
	tc.DataSourceConfig = fmt.Sprintf(cryptoKeyDataSourceTF, cryptoKeyName)

	s.BuildAndRun(ctx, tc)
}
