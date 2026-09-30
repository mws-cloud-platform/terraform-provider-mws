package queue

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	queuetest "go.mws.cloud/terraform-provider-mws/service/resources/queue/acctest"
)

var (
	//go:embed testdata/topic.tf
	topicTF string
	//go:embed testdata/topic_updated.tf
	topicUpdatedTF string
	//go:embed testdata/datasource/topic.tf
	topicDataSourceTF string
)

func TestTopicSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TopicSuite))
}

type TopicSuite struct {
	utils.ResourceSuite
}

func (s *TopicSuite) TestTopic() {
	ctx := s.T().Context()

	tc, err := queuetest.TopicTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	topicName := utils.RandResourceName("topic")
	tc.ResourceConfig = fmt.Sprintf(topicTF, topicName)
	tc.UpdatedResourceConfig = fmt.Sprintf(topicUpdatedTF, topicName)
	tc.DataSourceConfig = fmt.Sprintf(topicDataSourceTF, topicName)

	s.BuildAndRun(ctx, tc)
}
