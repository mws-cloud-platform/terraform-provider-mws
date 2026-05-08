package client_test

import (
	"context"

	"go.mws.cloud/terraform-provider-mws/internal/client"
)

func noopInvoke(context.Context, any, client.APIResp) error {
	return nil
}
