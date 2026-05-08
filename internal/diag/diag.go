package diag

import (
	"errors"
	"strings"

	mwsclient "go.mws.cloud/terraform-provider-mws/internal/client"
	mwserrors "go.mws.cloud/terraform-provider-mws/internal/errors"
)

// FormatError formats an error into a string.
func FormatError(err error) string {
	var sb strings.Builder

	sb.WriteString("Error: " + err.Error())

	var clientErr mwsclient.Error
	if errors.As(err, &clientErr) {
		var apiErr *mwserrors.APIError
		if errors.As(clientErr.Err, &apiErr) {
			if apiErr.Details != nil {
				sb.WriteString("\nDetails: ")
				sb.WriteString(apiErr.Details.String())
			}
		}

		sb.WriteString("\n\nRequestID: ")
		sb.WriteString(clientErr.RequestID)
		sb.WriteString("\nTraceID: ")
		sb.WriteString(clientErr.TraceID)
	}

	return sb.String()
}
