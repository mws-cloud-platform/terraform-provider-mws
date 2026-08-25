package vpc

import (
	"bytes"
	_ "embed"
	"testing"
	"text/template"

	"github.com/stretchr/testify/suite"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

//go:embed testdata/egress_nat.tf
var egressNatTF string

func TestEgressNatSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(EgressNatSuite))
}

type egressNatTemplateData struct {
	Name            string
	Network         string
	Subnet          string
	ExternalAddress string
}

type EgressNatSuite struct {
	BaseSubnetSuite

	addressSDK *vpcsdk.ExternalAddress

	externalAddressName string
	externalAddressID   string
	egressNatName       string
}

func (s *EgressNatSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.addressSDK, err = vpcsdk.NewExternalAddress(ctx, s.SDK)
	s.Require().NoError(err)

	s.externalAddressName = utils.RandResourceName("egress-nat-addr")
	s.egressNatName = utils.RandResourceName("egress-nat")

	addr, err := s.addressSDK.CreateExternalAddress(ctx, vpcclient.UpsertExternalAddressRequest{
		ExternalAddress: s.externalAddressName,
		Body: &vpcmodel.ExternalAddressRequest{
			Spec: vpcmodel.VpcExternalAddressSpecRequest{},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("external address %q created", s.externalAddressName)
	s.externalAddressID = addr.GetMetadata().GetId().ID()
}

func (s *EgressNatSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.addressSDK.DeleteExternalAddress(ctx, vpcclient.DeleteExternalAddressRequest{
		ExternalAddress: s.externalAddressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("external address %q deletion failed: %v", s.externalAddressName, err)
	} else {
		s.T().Logf("external address %q deleted", s.externalAddressName)
	}

	s.BaseSubnetSuite.TearDownSuite()
}

func (s *EgressNatSuite) TestEgressNat() {
	ctx := s.T().Context()

	tc, err := vpctest.EgressNatTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	templateData := egressNatTemplateData{
		Network:         s.NetworkName,
		Subnet:          s.Subnet.GetMetadata().GetId().ID(),
		ExternalAddress: s.externalAddressID,
		Name:            s.egressNatName,
	}

	tmpl, err := template.New("egress_nat").Parse(egressNatTF)
	s.Require().NoError(err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	s.Require().NoError(err)

	tc.ResourceConfig = buf.String()
	s.T().Log(tc.ResourceConfig)
	s.BuildAndRun(ctx, tc)
}
