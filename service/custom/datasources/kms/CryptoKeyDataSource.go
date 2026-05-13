package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"go.mws.cloud/go-sdk/service/kms/client"
	"go.mws.cloud/go-sdk/service/kms/model"
	resourcesdk "go.mws.cloud/go-sdk/service/kms/sdk"

	tfmodel "go.mws.cloud/terraform-provider-mws/service/datasources/kms/model"
)

func Read(ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
	sdk *resourcesdk.CryptoKey,
	data tfmodel.CryptoKeyModel) (*model.CryptoKeyOptionalResponse, error) {
	return sdk.GetCryptoKey(
		ctx,
		client.GetCryptoKeyRequest{
			Project: data.ProjectParam.ValueString(),
			Key:     data.KeyParam.ValueString(),
		},
	)
}
