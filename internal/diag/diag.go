package diag

import (
	"errors"
	"strings"

	mwserrors "go.mws.cloud/go-sdk/mws/errors"
)

// FormatError formats an error into a string.
func FormatError(err error) string {
	var sb strings.Builder

	sb.WriteString("Error: ")
	sb.WriteString(err.Error())

	if apiErr, ok := errors.AsType[*mwserrors.APIError](err); ok {
		if apiErr.Details != nil {
			sb.WriteString("\nDetails: ")
			sb.WriteString(apiErr.Details.String())
		}
	}

	if unexpectedStatusCodeErr, ok := errors.AsType[*mwserrors.UnexpectedStatusCodeError](err); ok {
		if len(unexpectedStatusCodeErr.Data) > 0 {
			sb.WriteString("\nData: ")
			sb.Write(unexpectedStatusCodeErr.Data)
		}
	}

	if clientErr, ok := err.(mwsClientError); ok {
		sb.WriteString("\n\nRequestID: ")
		sb.WriteString(clientErr.GetRequestID())
		sb.WriteString("\nTraceID: ")
		sb.WriteString(clientErr.GetTraceID())
	}

	return sb.String()
}

type mwsClientError interface {
	GetRequestID() string
	GetTraceID() string
}
