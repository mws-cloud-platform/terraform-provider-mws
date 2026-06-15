package utils

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/suite"
	"go.mws.cloud/go-sdk/mws"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap/zaptest"
)

const EnvTfAcc = "TF_ACC"

// ResourceTestCaseBuilder represents a builder for resource test cases.
type ResourceTestCaseBuilder interface {
	Build(context.Context) (resource.TestCase, error)
}

// ResourceSuite is a suite for provider resources testing.
type ResourceSuite struct {
	Suite
}

// BuildAndRun builds resource test case, injects the provider factories and
// runs it.
func (s *ResourceSuite) BuildAndRun(ctx context.Context, builder ResourceTestCaseBuilder) {
	tc, err := builder.Build(ctx)
	s.Require().NoError(err)
	s.Run(tc)
}

// Run injects the provider factories into the test case and runs it.
func (s *ResourceSuite) Run(tc resource.TestCase) {
	tc.ProtoV6ProviderFactories = ProtoV6ProviderFactories()
	resource.Test(s.T(), tc)
}

// Suite is a base suite for resources testing.
type Suite struct {
	suite.Suite

	SDK *mws.SDK
}

func (s *Suite) SetupSuite() {
	var err error

	if os.Getenv(EnvTfAcc) == "" {
		s.T().Skipf("Skipped unless env %q is set", EnvTfAcc)
		return
	}

	s.SDK, err = mws.Load(s.T().Context(),
		mws.WithLogger(zaptest.NewLogger(s.T())),
		mws.WithTracerProvider(sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)),
	)
	s.Require().NoError(err, "load sdk")
}

func (s *Suite) TearDownSuite() {
	s.Assert().NoError(s.SDK.Close(context.Background()), "close sdk")
}
