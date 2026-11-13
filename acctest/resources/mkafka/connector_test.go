package mkafka

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	mkafkaclient "go.mws.cloud/go-sdk/service/mkafka/client"
	mkafkamodel "go.mws.cloud/go-sdk/service/mkafka/model"
	mkafkasdk "go.mws.cloud/go-sdk/service/mkafka/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	mkafkatest "go.mws.cloud/terraform-provider-mws/service/resources/mkafka/acctest"
)

var (
	//go:embed testdata/connector.tf
	connectorTF string
	//go:embed testdata/datasource/connector.tf
	connectorDataSourceTF string
)

func TestConnectorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConnectorSuite))
}

type ConnectorSuite struct {
	MkafkaClusterTestSuite
	topicSDK      *mkafkasdk.Topic
	topicName     string
	connectorName string
}

func (s *ConnectorSuite) SetupSuite() {
	var err error

	s.MkafkaClusterTestSuite.SetupSuite()
	ctx := s.T().Context()

	s.topicSDK, err = mkafkasdk.NewTopic(ctx, s.SDK)
	s.Require().NoError(err)
	s.topicName = utils.RandResourceName("topic")
	s.connectorName = utils.RandResourceName("connector")

	for _, topic := range s.topics() {
		_, err = s.topicSDK.CreateKafkaTopic(ctx, mkafkaclient.UpsertKafkaTopicRequest{
			Cluster: s.kafkaName,
			Topic:   topic,
			Body: mkafkamodel.KafkaTopicRequest{
				Spec: mkafkamodel.KafkaTopicSpecRequest{
					Partitions:        new(int32(3)),
					ReplicationFactor: new(int32(1)),
				},
			},
		}, mkafkaclient.WithWait())
		s.Require().NoError(err)
		s.T().Logf("topic %q created", topic)
	}
}

func (s *ConnectorSuite) TearDownSuite() {
	ctx := s.T().Context()
	for _, topic := range s.topics() {
		if err := s.topicSDK.DeleteKafkaTopic(ctx, mkafkaclient.DeleteKafkaTopicRequest{
			Cluster: s.kafkaName,
			Topic:   topic,
		}, mkafkaclient.WithWait()); err != nil {
			s.T().Logf("topic %q deletion failed: %v", topic, err)
		}
	}
	s.MkafkaClusterTestSuite.TearDownSuite()
}

func (s *ConnectorSuite) TestConnector() {
	ctx := s.T().Context()
	tc, err := mkafkatest.ConnectorTestCase(ctx, s.SDK)
	s.Require().NoError(err)
	tc.ResourceConfig = fmt.Sprintf(connectorTF, s.kafkaName, s.connectorName, s.topicName)
	tc.DataSourceConfig = fmt.Sprintf(connectorDataSourceTF, s.kafkaName, s.connectorName)
	s.T().Log(tc.ResourceConfig)
	s.BuildAndRun(ctx, tc)
}

// topics returns the source topic and its DLQ topic, which must exist before the connector is created.
func (s *ConnectorSuite) topics() []string {
	return []string{s.topicName, s.topicName + "-dlq"}
}
