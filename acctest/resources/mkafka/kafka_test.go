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
	//go:embed testdata/kafka.tf
	kafkaTF string
	//go:embed testdata/datasource/kafka.tf
	kafkaDataSourceTF string
)

type KafkaData struct {
	Name              string
	NetworkID         string
	Broker1AddressRef string
}

func TestKafkaSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(KafkaSuite))
}

type KafkaSuite struct {
	MkafkaTestSuite
}

func (s *KafkaSuite) SetupSuite() {
	s.MkafkaTestSuite.SetupSuite()
}

func (s *KafkaSuite) TearDownSuite() {
	s.MkafkaTestSuite.TearDownSuite()
}

func (s *KafkaSuite) TestKafka() {
	ctx := s.T().Context()
	tc, err := mkafkatest.ClusterTestCase(ctx, s.SDK)
	s.Require().NoError(err)
	data := KafkaData{
		Name:              s.kafkaName,
		NetworkID:         s.network.GetMetadata().GetId().ID(),
		Broker1AddressRef: s.address1ID,
	}
	tpl := template.Must(template.New("kafkaTF").Parse(kafkaTF))
	sb := new(strings.Builder)
	err = tpl.Execute(sb, data)
	s.Require().NoError(err)
	tc.ResourceConfig = sb.String()
	tc.DataSourceConfig = fmt.Sprintf(kafkaDataSourceTF, data.Name)
	s.T().Log(tc.ResourceConfig)
	s.BuildAndRun(ctx, tc)
}
