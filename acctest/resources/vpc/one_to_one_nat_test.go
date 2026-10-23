package vpc

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	vpctest "go.mws.cloud/terraform-provider-mws/service/resources/vpc/acctest"
)

var (
	//go:embed testdata/one_to_one_nat.tf
	oneToOneNatTF string
	//go:embed testdata/datasource/one_to_one_nat.tf
	oneToOneNatDataSourceTF string
)

func TestOneToOneNatSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(OneToOneNatSuite))
}

type OneToOneNatSuite struct {
	BaseSubnetSuite

	addressSDK         *vpcsdk.Address
	externalAddressSDK *vpcsdk.ExternalAddress

	oneToOneNatName     string
	internalAddressName string
	internalAddressID   string
	externalAddressName string
	externalAddressID   string
}

func (s *OneToOneNatSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.oneToOneNatName = utils.RandResourceName("one-to-one-nat")
	s.internalAddressName = s.oneToOneNatName + "-address"
	s.externalAddressName = s.oneToOneNatName + "-external-address"

	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
	s.Require().NoError(err)

	s.externalAddressSDK, err = vpcsdk.NewExternalAddress(ctx, s.SDK)
	s.Require().NoError(err)

	subnetRef, err := vpcref.ParseSubnetRef(ctx, s.Subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	internalAddress, err := s.addressSDK.CreateAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.NetworkName,
		Address: s.internalAddressName,
		Body: &vpcmodel.AddressRequest{
			Spec: vpcmodel.VpcAddressSpecRequest{
				Subnet: subnetRef,
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("address %q created", s.internalAddressName)
	s.internalAddressID = internalAddress.GetMetadata().GetId().ID()

	externalAddress, err := s.externalAddressSDK.CreateExternalAddress(ctx, vpcclient.UpsertExternalAddressRequest{
		ExternalAddress: s.externalAddressName,
		Body: &vpcmodel.ExternalAddressRequest{
			Spec: vpcmodel.VpcExternalAddressSpecRequest{},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("external address %q created", s.externalAddressName)
	s.externalAddressID = externalAddress.GetMetadata().GetId().ID()
}

func (s *OneToOneNatSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.NetworkName,
		Address: s.internalAddressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("address %q deletion failed: %v", s.internalAddressName, err)
	} else {
		s.T().Logf("address %q deleted", s.internalAddressName)
	}

	if err := s.externalAddressSDK.DeleteExternalAddress(ctx, vpcclient.DeleteExternalAddressRequest{
		ExternalAddress: s.externalAddressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf("external address %q deletion failed: %v", s.externalAddressName, err)
	} else {
		s.T().Logf("external address %q deleted", s.externalAddressName)
	}

	s.BaseSubnetSuite.TearDownSuite()
}

func (s *OneToOneNatSuite) TestOneToOneNat() {
	ctx := s.T().Context()

	tc, err := vpctest.OneToOneNatTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	tc.ResourceConfig = fmt.Sprintf(oneToOneNatTF,
		s.oneToOneNatName,
		s.NetworkName,
		s.externalAddressID,
		s.internalAddressID,
	)
	s.T().Log(tc.ResourceConfig)

	tc.DataSourceConfig = fmt.Sprintf(oneToOneNatDataSourceTF, s.oneToOneNatName, s.NetworkName)

	s.BuildAndRun(ctx, tc)
}
