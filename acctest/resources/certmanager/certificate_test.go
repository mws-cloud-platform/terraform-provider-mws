package certmanager

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	certmanagertest "go.mws.cloud/terraform-provider-mws/service/resources/certmanager/acctest"
)

var (
	//go:embed testdata/certificate.tf
	certificateTF string
	//go:embed testdata/datasource/certificate.tf
	certificateDataSourceTF string
)

func TestCertificateSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CertificateSuite))
}

func TestCertificateUpdateSelfManagerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CertificateUpdateSelfManagerSuite))
}

type CertificateSuite struct {
	utils.ResourceSuite
}

func (s *CertificateSuite) TestCertificate() {
	ctx := s.T().Context()

	tc, err := certmanagertest.CertificateTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	certificateName := utils.RandResourceName("certificate")
	tc.ResourceConfig = fmt.Sprintf(certificateTF, certificateName, "testdata/certificate", "testdata/private_key", 1)
	tc.DataSourceConfig = fmt.Sprintf(certificateDataSourceTF, certificateName)
	s.BuildAndRun(ctx, tc)
}

type CertificateUpdateSelfManagerSuite struct {
	utils.ResourceSuite
}

func (s *CertificateUpdateSelfManagerSuite) TestCertificateChangeVersion() {
	ctx := s.T().Context()

	tc, err := certmanagertest.CertificateTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	certificateName := utils.RandResourceName("certificate")
	tc.ResourceConfig = fmt.Sprintf(certificateTF, certificateName, "testdata/certificate", "testdata/private_key", 1)
	tc.DataSourceConfig = fmt.Sprintf(certificateDataSourceTF, certificateName)
	tc.UpdatedResourceConfig = fmt.Sprintf(certificateTF, certificateName, "testdata/certificate", "testdata/private_key", 2)

	s.BuildAndRun(ctx, tc)
}
