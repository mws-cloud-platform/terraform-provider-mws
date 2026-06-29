package gpt

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	gpttest "go.mws.cloud/terraform-provider-mws/service/resources/gpt/acctest"
)

var (
	//go:embed testdata/deployment.tf
	deploymentTF string
	//go:embed testdata/datasource/deployment.tf
	deploymentDataSourceTF string
)

func TestDeploymentSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DeploymentSuite))
}

type DeploymentSuite struct {
	utils.ResourceSuite
}

func (s *DeploymentSuite) TestDeployment() {
	ctx := s.T().Context()

	tc, err := gpttest.DeploymentTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	deploymentName := utils.RandResourceName("deployment")
	tc.ResourceConfig = fmt.Sprintf(deploymentTF, deploymentName)
	tc.DataSourceConfig = fmt.Sprintf(deploymentDataSourceTF, deploymentName)

	s.BuildAndRun(ctx, tc)
}
