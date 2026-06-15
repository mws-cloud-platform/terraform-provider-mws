package mpostgres

import (
	"context"
	"time"

	"go.mws.cloud/go-sdk/mws/wait"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	"go.mws.cloud/go-sdk/service/common/model"
	mpostgresclient "go.mws.cloud/go-sdk/service/mpostgres/client"
	mpostgresmodel "go.mws.cloud/go-sdk/service/mpostgres/model"
	mpostgressdk "go.mws.cloud/go-sdk/service/mpostgres/sdk"
	"go.mws.cloud/go-sdk/service/resources/references/compute"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"
	vpcclient "go.mws.cloud/go-sdk/service/vpc/client"
	vpcmodel "go.mws.cloud/go-sdk/service/vpc/model"
	vpcsdk "go.mws.cloud/go-sdk/service/vpc/sdk"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

const clusterCreateTimeout = 15 * time.Minute

type BaseSuite struct {
	utils.ResourceSuite
	networkSDK *vpcsdk.Network
	subnetSDK  *vpcsdk.Subnet
	addressSDK *vpcsdk.Address

	networkName                string
	network                    *vpcmodel.NetworkOptionalResponse
	subnetName                 string
	subnet                     *vpcmodel.SubnetOptionalResponse
	primaryEndpointAddressName string
	primaryEndpointAddress     *vpcmodel.AddressOptionalResponse
}

func (s *BaseSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.networkName = utils.RandResourceName("mpostgres-network")
	s.subnetName = s.networkName + "-subnet"
	s.primaryEndpointAddressName = utils.RandResourceName("mpostgres-primary-endpoint-address")

	s.networkSDK, err = vpcsdk.NewNetwork(ctx, s.SDK)
	s.Require().NoError(err)
	s.subnetSDK, err = vpcsdk.NewSubnet(ctx, s.SDK)
	s.Require().NoError(err)
	s.addressSDK, err = vpcsdk.NewAddress(ctx, s.SDK)
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

	subnetRef, err := vpcref.ParseSubnetRef(ctx, s.subnet.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.primaryEndpointAddress, err = s.addressSDK.CreateAddress(ctx, vpcclient.UpsertAddressRequest{
		Network: s.networkName,
		Address: s.primaryEndpointAddressName,
		Body: &vpcmodel.AddressRequest{
			Spec: vpcmodel.VpcAddressSpecRequest{
				Subnet: subnetRef,
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("primary endpoint address %q created", s.primaryEndpointAddressName)
}

func (s *BaseSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.addressSDK.DeleteAddress(ctx, vpcclient.DeleteAddressRequest{
		Network: s.networkName,
		Address: s.primaryEndpointAddressName,
	}, vpcclient.WithWait()); err != nil {
		s.T().Logf(
			"primary endpoint address %q deletion failed: %v",
			s.primaryEndpointAddressName,
			err,
		)
	} else {
		s.T().Logf(
			"primary endpoint address %q deleted",
			s.primaryEndpointAddressName,
		)
	}

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
	clusterSDK *mpostgressdk.PostgresCluster

	clusterName string
	cluster     *mpostgresmodel.PostgresClusterResponse
}

func (s *BaseClusterSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.BaseSuite.SetupSuite()

	s.clusterName = utils.RandResourceName("cluster")

	networkRef, err := vpcref.ParseNetworkRef(ctx, s.network.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	primaryEndpointAddressRef, err := vpcref.ParseAddressRef(ctx, s.primaryEndpointAddress.GetMetadata().GetId().ID())
	s.Require().NoError(err)

	s.clusterSDK, err = mpostgressdk.NewPostgresCluster(ctx, s.SDK)
	s.Require().NoError(err)

	s.cluster, err = s.clusterSDK.CreatePostgresCluster(ctx, mpostgresclient.UpsertPostgresClusterRequest{
		Cluster: s.clusterName,
		Body: mpostgresmodel.PostgresClusterRequest{
			Spec: mpostgresmodel.PostgresClusterSpecRequest{
				Version: "17",
				Active:  true,
				Endpoints: []mpostgresmodel.PostgresEndpointRequest{{
					Name:    "primary-endpoint",
					Network: networkRef,
					PrimaryAddresses: []mpostgresmodel.PostgresNetworkAddressRequest{{
						Ref: new(primaryEndpointAddressRef),
					}},
				}},
				InstanceTemplate: mpostgresmodel.PostgresInstanceTemplateRequest{
					VmType: compute.NewVmTypeRef("gen-2-8"),
					Disk: mpostgresmodel.DataDiskSpecRequest{
						Size: bytesize.MustNewFromInt64(20, bytesize.GB),
						Type: mpostgresmodel.DataDiskType_NETWORK_STANDARD_SSD,
					},
				},
				Instances: []mpostgresmodel.PostgresInstanceRequest{{
					Count: 1,
					Zone:  new("ru-central1-a"),
				}},
			},
		},
	})
	s.Require().NoError(err)
	s.Require().NoError(waitForClusterReady(ctx, s.clusterSDK, s.clusterName))
	s.T().Logf("cluster %q created", s.clusterName)
}

func (s *BaseClusterSuite) TearDownSuite() {
	ctx := s.T().Context()

	_, err := s.clusterSDK.GetPostgresCluster(ctx, mpostgresclient.GetPostgresClusterRequest{
		Cluster: s.clusterName,
	}, mpostgresclient.WithWait())
	s.Assert().NoError(err)

	if err := s.clusterSDK.DeletePostgresCluster(ctx, mpostgresclient.DeletePostgresClusterRequest{
		Cluster: s.clusterName,
	}, mpostgresclient.WithWait()); err != nil {
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

func waitForClusterReady(ctx context.Context, clusterSDK *mpostgressdk.PostgresCluster, cluster string) error {
	callback := func(ctx context.Context) (*mpostgresmodel.PostgresClusterResponse, bool, error) {
		response, err := clusterSDK.GetPostgresCluster(ctx, mpostgresclient.GetPostgresClusterRequest{Cluster: cluster})

		status := new(response.GetStatus().GetReady()).GetState()
		state := ptr.Value(response.GetStatus().GetState())
		health := ptr.Value(response.GetStatus().GetHealth())

		// Sometimes cluster can go from OK to PROCESSING, it happens because of
		// internal mpostgres logic. As a workaround we can wait for RUNNING and
		// ALIVE state, in addition to OK reconciliation state.
		ok := status == model.ResourceStatusState_OK &&
			state == mpostgresmodel.ClusterState_RUNNING &&
			health == mpostgresmodel.ClusterHealth_ALIVE
		failed := status == model.ResourceStatusState_FAILED
		return response, ok || failed, err
	}

	waiter := wait.NewWaiter(callback, wait.WithTimeout(clusterCreateTimeout))
	_, err := waiter.Wait(ctx)
	return err
}
