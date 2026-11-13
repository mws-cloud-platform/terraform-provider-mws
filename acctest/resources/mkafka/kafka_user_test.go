package mkafka

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/suite"

	mkafkatest "go.mws.cloud/terraform-provider-mws/service/resources/mkafka/acctest"
)

var (
	//go:embed testdata/kafka_user.tf
	kafkaUserTF string
	//go:embed testdata/datasource/kafka_user.tf
	kafkaUserDataSourceTF string
)

type KafkaUserData struct {
	Name          string
	KafkaName     string
	AdminPassword string
}

func TestKafkaUserSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(KafkaUserSuite))
}

type KafkaUserSuite struct {
	MkafkaClusterTestSuite
}

func (s *KafkaUserSuite) TestKafkaUser() {
	ctx := s.T().Context()
	tc, err := mkafkatest.ClusterUserTestCase(ctx, s.SDK)
	s.Require().NoError(err)
	data := KafkaUserData{
		Name:          "test-user",
		KafkaName:     s.kafkaName,
		AdminPassword: "UserPassword123!",
	}
	tpl := template.Must(template.New("kafkaUserTF").Parse(kafkaUserTF))
	sb := new(strings.Builder)
	err = tpl.Execute(sb, data)
	s.Require().NoError(err)
	tc.ResourceConfig = sb.String()
	tc.DataSourceConfig = fmt.Sprintf(kafkaUserDataSourceTF, data.KafkaName, data.Name)
	s.T().Log(tc.ResourceConfig)
	s.BuildAndRun(ctx, tc)
}
