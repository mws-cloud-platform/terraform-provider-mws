package acctest

import (
	"context"
	"fmt"

	"go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/service/kms/client"
	sdkmodel "go.mws.cloud/go-sdk/service/kms/model"
	resourcesdk "go.mws.cloud/go-sdk/service/kms/sdk"
	kmsref "go.mws.cloud/go-sdk/service/resources/references/kms"
)

func GetCryptoKey(ctx context.Context, sdk *resourcesdk.CryptoKey, id string) (*sdkmodel.CryptoKeyOptionalResponse, error) {
	ref, err := kmsref.ParseCryptoKeyRef(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("parse reference: %w", err)
	}
	key, err := sdk.GetCryptoKey(ctx, client.GetCryptoKeyRequest{
		Key:     string(ref.ResourceName()),
		Project: ref.GetProject(),
	})
	if err != nil {
		return nil, err
	}
	if key.GetStatus().Destruction != nil {
		// not show this key
		return nil, errors.NewAPIError(404, errors.NotFound, "key prepared for destruction")
	}
	return key, nil
}
