package iam

import (
	iamclient "go.mws.cloud/go-sdk/service/iam/client"
	iamsdk "go.mws.cloud/go-sdk/service/iam/sdk"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

type baseServiceAccountSuite struct {
	utils.ResourceSuite
	serviceAccountSDK *iamsdk.ServiceAccount

	serviceAccountName string
}

func (s *baseServiceAccountSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.serviceAccountName = utils.RandResourceName("sa")

	s.serviceAccountSDK, err = iamsdk.NewServiceAccount(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.serviceAccountSDK.CreateServiceAccount(ctx, iamclient.UpsertServiceAccountRequest{
		ServiceAccount: s.serviceAccountName,
	}, iamclient.WithWait())
	s.Require().NoError(err)
	s.T().Logf("service account %q created", s.serviceAccountName)
}

func (s *baseServiceAccountSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.serviceAccountSDK.DeleteServiceAccount(ctx, iamclient.DeleteServiceAccountRequest{
		ServiceAccount: s.serviceAccountName,
	}, iamclient.WithWait()); err != nil {
		s.T().Logf("service account %q deletion failed: %v", s.serviceAccountName, err)
	} else {
		s.T().Logf("service account %q deleted", s.serviceAccountName)
	}

	s.ResourceSuite.TearDownSuite()
}
