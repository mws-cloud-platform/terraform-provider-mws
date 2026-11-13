package mkafka

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/suite"

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
	MkafkaClusterTestSuite
	topicName string
}

func (s *TopicSuite) SetupSuite() {
	s.MkafkaClusterTestSuite.SetupSuite()
	s.topicName = utils.RandResourceName("topic")
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
