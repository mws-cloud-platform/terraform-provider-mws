package mclickhouse

import (
	helper "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	commonmodel "go.mws.cloud/go-sdk/service/common/model"
	mclickhouseclient "go.mws.cloud/go-sdk/service/mclickhouse/client"
	mclickhousemodel "go.mws.cloud/go-sdk/service/mclickhouse/model"
	mclickhousesdk "go.mws.cloud/go-sdk/service/mclickhouse/sdk"
	computeref "go.mws.cloud/go-sdk/service/resources/references/compute"
	mclickhouseref "go.mws.cloud/go-sdk/service/resources/references/mclickhouse"
	rmref "go.mws.cloud/go-sdk/service/resources/references/rm"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

type BaseSuite struct {
	utils.ResourceSuite
	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet

	networkName string
	network     *vpcmodel.NetworkOptionalResponse
	subnetName  string
	subnet      *vpcmodel.SubnetOptionalResponse
}

func (s *BaseSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.networkName = utils.RandResourceName("mclickhouse-network")
	s.subnetName = s.networkName + "-subnet"

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)
	s.subnetSDK, err = vpcsdk.NewSubnet(ctx, s.SDK)
	s.Require().NoError(err)

	s.network, err = s.networkSDK.CreateNetwork(ctx, vpcclient.UpsertNetworkRequest{
		Network: s.networkName,
	})
	s.Require().NoError(err)
	s.T().Logf("network %q created", s.networkName)

	s.subnet, err = s.subnetSDK.CreateSubnet(ctx, vpcclient.UpsertSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
		Body: vpcmodel.SubnetRequest{
			Spec: vpcmodel.SubnetSpecRequest{
				Cidr: cidraddress.MustParseCIDR4AddressString("10.243.0.0/18"),
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("subnet %q created", s.subnetName)
}

func (s *BaseSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.subnetSDK.DeleteSubnet(ctx, vpcclient.DeleteSubnetRequest{
		Network: s.networkName,
		Subnet:  s.subnetName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"subnet %q deletion failed: %v",
			s.subnetName,
			err,
		)
	} else {
		s.T().Logf(
			"subnet %q deleted",
			s.subnetName,
		)
	}

	if err := s.networkSDK.DeleteNetwork(ctx, vpcclient.DeleteNetworkRequest{
		Network: s.networkName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"network %q deletion failed: %v",
			s.networkName,
			err,
		)
	} else {
		s.T().Logf(
			"network %q deleted",
			s.networkName,
		)
	}

	s.ResourceSuite.TearDownSuite()
}

type BaseClusterSuite struct {
	BaseSuite
	clusterSDK *mclickhousesdk.ClickhouseCluster

	clusterName string
	cluster     *mclickhousemodel.ClickhouseClusterOptionalResponse
}

func (s *BaseClusterSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSuite.SetupSuite()

	s.clusterName = utils.RandResourceName("cluster")

	subnetRef, err := vpcref.ParseSubnetRef(ctx, s.subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.clusterSDK, err = mclickhousesdk.NewClickhouseCluster(ctx, s.SDK)
	s.Require().NoError(err)

	zoneRef, err := rmref.ParseZoneRef(s.T().Context(), "ru-central1-a")
	s.Require().NoError(err)

	s.cluster, err = s.clusterSDK.CreateClickhouseCluster(ctx, mclickhouseclient.UpsertClickhouseClusterRequest{
		Cluster: s.clusterName,
		Body: &mclickhousemodel.ClickhouseClusterRequest{
			Spec: mclickhousemodel.ClickhouseClusterSpecRequest{
				Version: "25.3",
				Active:  new(true),
				Shards: []mclickhousemodel.ClickhouseClusterShardRequest{{
					Name: "shard",
					Resources: mclickhousemodel.ClickhouseInstanceHWResourcesRequest{
						VmType: mclickhouseref.NewClickhouseVmTypeRef("gen-4-8"),
						Disk: mclickhousemodel.ClickhouseInstanceDiskSpecRequest{
							Type: computeref.NewDiskTypeRef("NETWORK_STANDARD_SSD"),
							Size: bytesize.MustParseString("10GB"),
						},
					},
					Instances: []mclickhousemodel.ClickhouseClusterInstanceRequest{{
						Name:  new("instance-1"),
						Zone:  zoneRef,
						Count: new(1),
						Endpoints: []mclickhousemodel.ClickhouseEndpointRequest{{
							Address: mclickhousemodel.ClickhouseEndpointAddressSpecOrRefRequest{
								Spec: &mclickhousemodel.ClickhouseEndpointAddressSpecRequest{
									Subnet: subnetRef,
								},
							},
							ExternalAddress: &mclickhousemodel.ClickhouseEndpointExternalAddressSpecOrRefRequest{
								Spec: &mclickhousemodel.ClickhouseEndpointExternalAddressSpecOrRefSpecRequest{},
							},
						}},
					}},
				}},
				BootstrapAdmin: mclickhousemodel.ClickhouseClusterBootstrapAdminSpecRequest{
					Username: "admin",
					Password: helper.RandString(16),
				},
				MaintenanceWindow: &commonmodel.MaintenanceWindowRequest{
					Weekly: commonmodel.WeeklyMaintenanceWindowRequest{
						Days: []commonmodel.DayOfWeek{"MONDAY"},
						Hour: 3,
					},
				},
				Backup: &mclickhousemodel.ClickhouseClusterBackupRequest{
					Hour:             new(2),
					RetainPeriodDays: new(7),
				},
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("cluster %q created", s.clusterName)
}

func (s *BaseClusterSuite) TearDownSuite() {
	ctx := s.T().Context()

	_, err := s.clusterSDK.GetClickhouseCluster(ctx, mclickhouseclient.GetClickhouseClusterRequest{
		Cluster: s.clusterName,
	}, mclickhouseclient.WithWait())
	s.Assert().NoError(err)

	if err := s.clusterSDK.DeleteClickhouseCluster(ctx, mclickhouseclient.DeleteClickhouseClusterRequest{
		Cluster: s.clusterName,
	}, mclickhouseclient.WithWait()); err != nil {
		s.T().Logf(
			"cluster %q deletion failed: %v",
			s.clusterName,
			err,
		)
	} else {
		s.T().Logf(
			"cluster %q deleted",
			s.clusterName,
		)
	}

	s.BaseSuite.TearDownSuite()
}
