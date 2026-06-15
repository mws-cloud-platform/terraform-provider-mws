package mk8s

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/suite"
	mk8sclient "go.mws.cloud/go-sdk/service/mk8s/client"
	mk8smodel "go.mws.cloud/go-sdk/service/mk8s/model"
	mk8ssdk "go.mws.cloud/go-sdk/service/mk8s/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mk8stest "go.mws.cloud/terraform-provider-mws/service/resources/mk8s/acctest"
)

var (
	//go:embed testdata/node_group.tf
	nodeGroupTF string
	//go:embed testdata/datasource/node_group.tf
	nodeGroupDataSourceTF string
)

type NodeGroupData struct {
	Name          string
	ClusterName   string
	Zone          string
	Channel       string
	MaintanceDays []string
	SubnetID      string
	Autoscalling  struct {
		Min int
		Max int
	}
	Rollout struct {
		Surge       int
		Unavailable int
	}
	StorageSize string
	Version     string
}

func (c NodeGroupData) GetDays() string {
	return utils.SliceStringJoin(c.MaintanceDays)
}

func TestNodeGroupSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NodeGroupSuite))
}

type NodeGroupSuite struct {
	Mk8sTestSuite
	clusterSDK  *mk8ssdk.Mk8sCluster
	clusterName string
}

func (s *NodeGroupSuite) SetupSuite() {
	var version = "v1.34.1-mws.1"
	var err error
	ctx := s.T().Context()
	s.Mk8sTestSuite.SetupSuite()
	s.clusterName = utils.RandResourceName("cluster")

	s.clusterSDK, err = mk8ssdk.NewMk8sCluster(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.clusterSDK.CreateMk8sCluster(ctx, mk8sclient.UpsertMk8sClusterRequest{
		ClusterName: s.clusterName,
		Body: mk8smodel.ClusterRequest{
			Spec: mk8smodel.ClusterSpecRequest{
				Availability: mk8smodel.ClusterAvailabilitySpecRequest{
					Standalone: &mk8smodel.ClusterAvailabilitySpecStandaloneRequest{
						Zone: "ru-central1-a",
					},
				},
				Network: mk8smodel.ClusterSpecNetworkRequest{
					PrimaryEndpoint: mk8smodel.ClusterPrimaryEndpointSpecOrRefRequest{
						Ref: &s.primaryNetworkInterfaceAddressID,
					},
					PodsCidr:     s.subnetPod.Spec.Cidr,
					ServicesCidr: s.subnetService.Spec.Cidr,
				},
				VersionControl: mk8smodel.ClusterVersionControlSpecRequest{
					ReleaseChannel: "stable",
					Version:        &version,
				},
			},
		},
	})
	s.Require().NoError(err)
	s.T().Logf("mk8s cluster %q created", s.clusterName)
}

func (s *NodeGroupSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.clusterSDK.DeleteMk8sCluster(ctx, mk8sclient.DeleteMk8sClusterRequest{
		ClusterName: s.clusterName,
	}, mk8sclient.WithWait()); err != nil {
		s.T().Logf("cluster %q deletion failed: %v", s.clusterName, err)
	} else {
		s.T().Logf("cluster %q deleted", s.clusterName)
	}

	s.Mk8sTestSuite.TearDownSuite()
}

func (s *NodeGroupSuite) TestCluster() {
	ctx := s.T().Context()
	tc, err := mk8stest.NodeGroupTestCase(ctx, s.SDK)
	s.Require().NoError(err)
	data := NodeGroupData{
		ClusterName:   s.clusterName,
		Name:          utils.RandResourceName("node-group"),
		Zone:          "ru-central1-a",
		Channel:       "stable",
		MaintanceDays: []string{"MONDAY"},
		Autoscalling: struct {
			Min int
			Max int
		}{1, 3},
		Rollout: struct {
			Surge       int
			Unavailable int
		}{0, 1},
		StorageSize: "20Gb",
		SubnetID:    s.subnetPod.GetMetadata().GetId().ID(),
		Version:     "v1.34.1-mws.1",
	}
	tpl := template.Must(template.New("nodeGroupTF").Parse(nodeGroupTF))
	sb := new(strings.Builder)
	err = tpl.Execute(sb, data)
	s.Require().NoError(err)
	tc.ResourceConfig = sb.String()
	tc.DataSourceConfig = fmt.Sprintf(nodeGroupDataSourceTF, data.ClusterName, data.Name)
	s.T().Log(tc.ResourceConfig)
	s.BuildAndRun(ctx, tc)
}
