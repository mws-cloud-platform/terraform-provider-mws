package mk8s

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mk8stest "go.mws.cloud/terraform-provider-mws/service/resources/mk8s/acctest"
)

var (
	//go:embed testdata/cluster.tf
	clusterTF string
	//go:embed testdata/datasource/cluster.tf
	clusterDataSourceTF string
)

type ClusterData struct {
	Name          string
	Zone          string
	IsHA          bool
	Channel       string
	MaintanceDays []string
	CIDR          string
	PodCIDR       string
	Endpoint      string
}

func (c ClusterData) GetDays() string {
	return utils.SliceStringJoin(c.MaintanceDays)
}

func TestStandaloneClusterSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ClusterSuite{IsHA: false})
}

func TestHAClusterSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ClusterSuite{IsHA: true})
}

type ClusterSuite struct {
	Mk8sTestSuite
	IsHA bool
}

func (s *ClusterSuite) TestClusterStandalone() {
	ctx := s.T().Context()
	tc, err := mk8stest.ClusterTestCase(ctx, s.SDK)
	s.Require().NoError(err)
	data := ClusterData{
		Name:          utils.RandResourceName("cluster"),
		Zone:          "ru-central1-a",
		Channel:       "stable",
		IsHA:          s.IsHA,
		MaintanceDays: []string{"MONDAY"},
		CIDR:          s.subnetService.GetSpec().Cidr.String(),
		PodCIDR:       s.subnetPod.GetSpec().Cidr.String(),
		Endpoint:      s.primaryNetworkInterfaceAddress.GetMetadata().GetId().ID(),
	}
	tpl := template.Must(template.New("clusterTF").Parse(clusterTF))
	sb := new(strings.Builder)
	err = tpl.Execute(sb, data)
	s.Require().NoError(err)
	tc.ResourceConfig = sb.String()
	tc.DataSourceConfig = fmt.Sprintf(clusterDataSourceTF, data.Name)
	s.T().Log(tc.ResourceConfig)
	s.BuildAndRun(ctx, tc)
}
