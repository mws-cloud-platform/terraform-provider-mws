package queue

import (
	_ "embed"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/mws/wait"
	queueclient "go.mws.cloud/go-sdk/service/queue/client"
	queuemodel "go.mws.cloud/go-sdk/service/queue/model"
	queuesdk "go.mws.cloud/go-sdk/service/queue/sdk"

	iamacctest "go.mws.cloud/terraform-provider-mws/acctest/resources/iam"
	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	queuetest "go.mws.cloud/terraform-provider-mws/service/resources/queue/acctest"
)

var (
	//go:embed testdata/topic_role_binding.tf
	topicRoleBindingTF string
	//go:embed testdata/topic_role_binding_updated.tf
	topicRoleBindingUpdatedTF string
	//go:embed testdata/datasource/topic_role_binding.tf
	topicRoleBindingDataSourceTF string
)

func TestTopicRoleBindingSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TopicRoleBindingSuite))
}

type TopicRoleBindingSuite struct {
	iamacctest.BaseServiceAccountSuite

	project   string
	topicName string

	topicSDK *queuesdk.Topic
}

func (s *TopicRoleBindingSuite) SetupSuite() {
	var err error

	s.BaseServiceAccountSuite.SetupSuite()
	ctx := s.T().Context()

	s.project = s.SDK.DefaultProject()
	s.topicName = utils.RandResourceName("topic")

	s.topicSDK, err = queuesdk.NewTopic(ctx, s.SDK)
	s.Require().NoError(err, "init topic sdk")

	partitionCount := int32(3)
	_, err = s.topicSDK.CreateTopic(ctx, queueclient.UpsertTopicRequest{
		Project: s.project,
		Topic:   s.topicName,
		Body: queuemodel.TopicRequest{
			Spec: queuemodel.TopicSpecRequest{
				TopicItemCommonRequest: queuemodel.TopicItemCommonRequest{
					PartitionCount: &partitionCount,
				},
			},
		},
	}, queueclient.WithWait(wait.WithTimeout(time.Hour)))
	s.Require().NoError(err, "create topic")
}

func (s *TopicRoleBindingSuite) TearDownSuite() {
	ctx := s.T().Context()

	err := s.topicSDK.DeleteTopic(ctx, queueclient.DeleteTopicRequest{
		Project: s.project,
		Topic:   s.topicName,
	}, queueclient.WithWait(wait.WithTimeout(time.Hour)))
	s.Assert().NoError(err, "delete topic")

	s.BaseServiceAccountSuite.TearDownSuite()
}

func (s *TopicRoleBindingSuite) TestTopicRoleBinding() {
	ctx := s.T().Context()

	tc, err := queuetest.TopicRoleBindingTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	roleBindingName := utils.RandResourceName("role-binding")
	serviceAccountRef := s.ServiceAccountRef.IDPath()
	tc.ResourceConfig = fmt.Sprintf(topicRoleBindingTF, s.topicName, roleBindingName, serviceAccountRef)
	tc.UpdatedResourceConfig = fmt.Sprintf(topicRoleBindingUpdatedTF, s.topicName, roleBindingName, serviceAccountRef)
	tc.DataSourceConfig = fmt.Sprintf(topicRoleBindingDataSourceTF, s.topicName, roleBindingName)

	s.BuildAndRun(ctx, tc)
}
