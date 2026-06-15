package vpc

import (
	"bytes"
	_ "embed"
	"testing"
	"text/template"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
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
	utils.ResourceSuite

	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet
	addressSDK *vpcsdk.ExternalAddress

	networkName         string
	subnetName          string
	subnetID            string
	externalAddressName string
	externalAddressID   string
	egressNatName       string
}

func (s *EgressNatSuite) SetupSuite() {
	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	var err error

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)

	s.subnetSDK, err = vpcsdk.NewSubnet(ctx, s.SDK)
	s.Require().NoError(err)

	s.addressSDK, err = vpcsdk.NewExternalAddress(ctx, s.SDK)
	s.Require().NoError(err)

	s.networkName = utils.RandResourceName("egress-nat-network")
	s.subnetName = utils.RandResourceName("egress-nat-subnet")
	s.externalAddressName = utils.RandResourceName("egress-nat-addr")
	s.egressNatName = utils.RandResourceName("egress-nat")

	_, err = s.networkSDK.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: s.networkName,
	})
	s.Require().NoError(err)
	s.T().Logf("network %q created", s.networkName)

	subnet, err := s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("192.168.0.0/17"),
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetName)
	s.subnetID = subnet.GetMetadata().GetId().ID()
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

	if err := s.subnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("subnet %q deletion failed: %v", s.subnetName, err)
	} else {
		s.T().Logf("subnet %q deleted", s.subnetName)
	}

	if err := s.networkSDK.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
		Network: s.networkName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("network %q deletion failed: %v", s.networkName, err)
	} else {
		s.T().Logf("network %q deleted", s.networkName)
	}

	s.ResourceSuite.TearDownSuite()
}

func (s *EgressNatSuite) TestEgressNat() {
	ctx := s.T().Context()

	tc, err := vpctest.EgressNatTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	templateData := egressNatTemplateData{
		Network:         s.networkName,
		Subnet:          s.subnetID,
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
