package compute

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
	computetest "go.mws.cloud/terraform-provider-mws/service/resources/compute/acctest"
)

var (
	//go:embed testdata/image.tf
	imageTF string
	//go:embed testdata/datasource/image.tf
	imageDataSourceTF string
)

func TestImageSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImageSuite))
}

type ImageSuite struct {
	utils.ResourceSuite
}

func (s *ImageSuite) TestImage() {
	ctx := s.T().Context()

	tc, err := computetest.ImageTestCase(ctx, s.SDK)
	s.Require().NoError(err)

	imageName := utils.RandResourceName("image")
	tc.ResourceConfig = fmt.Sprintf(imageTF, imageName)
	tc.DataSourceConfig = fmt.Sprintf(imageDataSourceTF, imageName)

	s.BuildAndRun(ctx, tc)
}
