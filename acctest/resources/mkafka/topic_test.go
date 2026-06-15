package mkafka

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/mws/wait"
	"go.mws.cloud/go-sdk/pkg/apimodels/units/bytesize"
	mkafkaclient "go.mws.cloud/go-sdk/service/mkafka/client"
	mkafkamodel "go.mws.cloud/go-sdk/service/mkafka/model"
	mkafkasdk "go.mws.cloud/go-sdk/service/mkafka/sdk"
	"go.mws.cloud/go-sdk/service/resources/references/compute"
	vpcref "go.mws.cloud/go-sdk/service/resources/references/vpc"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mkafkatest "go.mws.cloud/terraform-provider-mws/service/resources/mkafka/acctest"
)

var (
	//go:embed testdata/topic.tf
	topicTF string
	//go:embed testdata/datasource/topic.tf
	topicDataSourceTF string
)

type TopicData struct {
	KafkaName string
	Name      string
}

func TestTopicSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TopicSuite))
}

type TopicSuite struct {
	MkafkaTestSuite
	kafkaSDK     *mkafkasdk.Kafka
	kafkaCluster *mkafkamodel.KafkaClusterResponse
	topicName    string
}

func (s *TopicSuite) SetupSuite() {
	var err error

	s.MkafkaTestSuite.SetupSuite()
	ctx := s.T().Context()

	s.kafkaSDK, err = mkafkasdk.NewKafka(ctx, s.SDK)
	s.Require().NoError(err)
	s.topicName = utils.RandResourceName("topic")
	networkRef, err := vpcref.ParseNetworkRef(ctx, s.network.GetMetadata().GetId().ID())
	s.Require().NoError(err)
	addressRef, err := vpcref.ParseAddressRef(ctx, s.address1ID)
	s.Require().NoError(err)
	gen24 := compute.NewVmTypeRef("gen-2-4")
	disk10GB := bytesize.MustParseString("10GB")
	isActive := true
	s.kafkaCluster, err = s.kafkaSDK.CreateKafkaCluster(ctx, mkafkaclient.UpsertKafkaClusterRequest{
		Cluster: s.kafkaName,
		Body: mkafkamodel.KafkaClusterRequest{
			Spec: mkafkamodel.KafkaClusterSpecRequest{
				Version: "4.0",
				Active:  &isActive,
				Endpoints: []mkafkamodel.KafkaEndpointRequest{{
					Name:    "vpc-endpoints",
					Network: networkRef,
					BrokerAddresses: []mkafkamodel.KafkaEndpointBrokerAddressRequest{{
						Ref: &addressRef,
					}},
				}},
				Instances: mkafkamodel.KafkaInstanceRequest{
					Broker: mkafkamodel.KafkaInstanceSpecRequest{
						VmType: gen24,
						Disk:   mkafkamodel.KafkaDataDiskSpecRequest{Size: disk10GB},
						Allocation: []mkafkamodel.KafkaAllocationRequest{{
							Zone:  "ru-central1-a",
							Count: 1,
						}},
					},
					Controller: mkafkamodel.KafkaControllerInstanceSpecRequest{
						CombinedWithBroker: new(true),
					},
				},
			},
		},
	}, mkafkaclient.WithWait(wait.WithTimeout(time.Hour)))

	s.Require().NoError(err)
}

func (s *TopicSuite) TearDownSuite() {
	err := s.kafkaSDK.DeleteKafkaCluster(s.T().Context(), mkafkaclient.DeleteKafkaClusterRequest{Cluster: s.kafkaName}, mkafkaclient.WithWait(wait.WithTimeout(time.Hour)))
	s.Require().NoError(err)
	s.MkafkaTestSuite.TearDownSuite()
}

func (s *TopicSuite) TestTopic() {
	ctx := s.T().Context()
	tc, err := mkafkatest.TopicTestCase(ctx, s.SDK)
	s.Require().NoError(err)
	data := TopicData{
		KafkaName: s.kafkaName,
		Name:      s.topicName,
	}
	tpl := template.Must(template.New("topicTF").Parse(topicTF))
	sb := new(strings.Builder)
	err = tpl.Execute(sb, data)
	s.Require().NoError(err)
	tc.ResourceConfig = sb.String()
	tc.DataSourceConfig = fmt.Sprintf(topicDataSourceTF, data.KafkaName, data.Name)
	s.T().Log(tc.ResourceConfig)
	s.BuildAndRun(ctx, tc)
}
