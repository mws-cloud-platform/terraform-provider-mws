package gpt

import (
	_ "embed"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/service/resources/references/gpt"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

var (
	//go:embed testdata/datasource/model.tf
	modelDataSourceTF string
)

func TestModelSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ModelSuite))
}

type ModelSuite struct {
	utils.Suite
}

func (s *ModelSuite) TestModel() {
	project := s.SDK.DefaultProject()
	model := "glm-4.6-357b"
	metadataID := gpt.NewModelID(project, model)

	steps := []resource.TestStep{
		{
			Config: modelDataSourceTF,
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				return resource.TestCheckResourceAttr("data.mws_gpt_model.model_data", "metadata.id", metadataID.String())(state)
			},
		},
	}
	tc := resource.TestCase{
		Steps:                    steps,
		ProtoV6ProviderFactories: utils.ProtoV6ProviderFactories(),
	}
	resource.Test(s.T(), tc)
}
