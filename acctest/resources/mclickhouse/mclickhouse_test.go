package mclickhouse

import (
	helper "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"go.mws.cloud/go-sdk/pkg/apimodels/sensitive"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	commonmodel "go.mws.cloud/go-sdk/service/common/model"
	mclickhouseclient "go.mws.cloud/go-sdk/service/mclickhouse/client"
	mclickhousemodel "go.mws.cloud/go-sdk/service/mclickhouse/model"
	mclickhousesdk "go.mws.cloud/go-sdk/service/mclickhouse/sdk"
	mclickhouseref "go.mws.cloud/go-sdk/service/resources/references/mclickhouse"
	rmref "go.mws.cloud/go-sdk/service/resources/references/rm"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"

	"go.mws.cloud/terraform-provider-mws/acctest/resources/vpc"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

type BaseClusterSuite struct {
	vpc.BaseSubnetSuite
	clusterSDK *mclickhousesdk.ClickhouseCluster

	clusterName string
	cluster     *mclickhousemodel.ClickhouseClusterOptionalResponse
}

func (s *BaseClusterSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSubnetSuite.SetupSuite()

	s.clusterName = utils.RandResourceName("cluster")

	subnetRef, err := vpcref.ParseSubnetRef(ctx, s.Subnet.GetMetadata().GetId().ID())
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
						VmType: mclickhouseref.NewMustClickhouseVmTypeRef("gen-4-8"),
						Disk: mclickhousemodel.ClickhouseInstanceDiskSpecRequest{
							Type: mclickhousemodel.ClickhouseDataDiskType_NETWORK_STANDARD_SSD,
							Size: bytesize.MustParseString("10GB"),
						},
					},
					Instances: []mclickhousemodel.ClickhouseClusterInstanceRequest{{
						Name:  "instance-1",
						Zone:  zoneRef,
						Count: new(1),
						Endpoints: []mclickhousemodel.ClickhouseEndpointRequest{{
							Address: mclickhousemodel.ClickhouseEndpointAddressSpecOrRefRequest{
								Spec: &mclickhousemodel.ClickhouseEndpointAddressSpecRequest{
									Subnet: subnetRef,
								},
							},
							ExternalAddress: &mclickhousemodel.ClickhouseEndpointExternalAddressSpecOrRefRequest{
								Spec: &mclickhousemodel.ClickhouseEndpointExternalAddressSpecRequest{},
							},
						}},
					}},
				}},
				BootstrapAdmin: mclickhousemodel.ClickhouseClusterBootstrapAdminSpecRequest{
					Username: "admin",
					Password: sensitive.New(helper.RandString(16)),
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

	s.BaseSubnetSuite.TearDownSuite()
}
