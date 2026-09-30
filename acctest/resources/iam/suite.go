package iam

import (
	iamclient "go.mws.cloud/go-sdk/service/iam/client"
	iamsdk "go.mws.cloud/go-sdk/service/iam/sdk"
	iamref "go.mws.cloud/go-sdk/service/resources/references/iam"

	"go.mws.cloud/terraform-provider-mws/acctest/utils"
)

// BaseServiceAccountSuite is a suite that creates a service account before the
// tests and deletes it afterwards. Embed it into suites that need a service
// account for their resource under test.
type BaseServiceAccountSuite struct {
	utils.ResourceSuite

	ServiceAccountName string
	ServiceAccountRef  iamref.ServiceAccountRef
	serviceAccountSDK  *iamsdk.ServiceAccount
}

func (s *BaseServiceAccountSuite) SetupSuite() {
	var err error

	ctx := s.T().Context()
	s.ResourceSuite.SetupSuite()

	s.ServiceAccountName = utils.RandResourceName("sa")

	s.serviceAccountSDK, err = iamsdk.NewServiceAccount(ctx, s.SDK)
	s.Require().NoError(err)

	_, err = s.serviceAccountSDK.CreateServiceAccount(ctx, iamclient.UpsertServiceAccountRequest{
		ServiceAccount: s.ServiceAccountName,
	}, iamclient.WithWait())
	s.Require().NoError(err)

	s.ServiceAccountRef, err = iamref.NewServiceAccountRef(s.SDK.DefaultProject(), s.ServiceAccountName)
	s.Require().NoError(err)

	s.T().Logf("service account %q created", s.ServiceAccountName)
}

func (s *BaseServiceAccountSuite) TearDownSuite() {
	ctx := s.T().Context()

	if err := s.serviceAccountSDK.DeleteServiceAccount(ctx, iamclient.DeleteServiceAccountRequest{
		ServiceAccount: s.ServiceAccountName,
	}, iamclient.WithWait()); err != nil {
		s.T().Logf("service account %q deletion failed: %v", s.ServiceAccountName, err)
	} else {
		s.T().Logf("service account %q deleted", s.ServiceAccountName)
	}

	s.ResourceSuite.TearDownSuite()
}
