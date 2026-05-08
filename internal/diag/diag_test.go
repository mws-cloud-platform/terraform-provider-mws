package diag_test

import (
	"net/http"
	"testing"

	"go.mws.cloud/util-toolset/pkg/testing/golden"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"

	mwsclient "go.mws.cloud/terraform-provider-mws/internal/client"
	"go.mws.cloud/terraform-provider-mws/internal/diag"
	mwserrors "go.mws.cloud/terraform-provider-mws/internal/errors"
)

func TestFormatError(t *testing.T) {
	dir := golden.NewDir(t,
		golden.WithPath("testdata/format_error"),
		golden.WithRecreateOnUpdate())

	requestID := "123"
	traceID := "xxx-xxx"
	simpleErr := consterr.Error("error")

	for _, tc := range []struct {
		name string
		err  error
	}{
		{name: "simple", err: simpleErr},
		{
			name: "client_error",
			err: mwsclient.Error{
				Err:       simpleErr,
				RequestID: requestID,
				TraceID:   traceID,
			},
		},
		{
			name: "client_api_error",
			err: mwsclient.Error{
				Err: &mwserrors.APIError{
					Code:        http.StatusBadRequest,
					Status:      mwserrors.InvalidArgument,
					Description: "invalid data",
				},
				RequestID: requestID,
				TraceID:   traceID,
			},
		},
		{
			name: "client_api_error_with_details",
			err: mwsclient.Error{
				Err: &mwserrors.APIError{
					Code:        http.StatusPreconditionFailed,
					Status:      mwserrors.FailedPrecondition,
					Description: "precondition failed",
					Details: mwserrors.Details{
						"foo": "bar",
					},
				},
				RequestID: requestID,
				TraceID:   traceID,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir.String(t, tc.name+".txt", diag.FormatError(tc.err))
		})
	}
}
