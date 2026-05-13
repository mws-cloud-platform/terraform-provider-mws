package kms

import (
	"context"
	"errors"
	"time"

	tfdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"go.mws.cloud/go-sdk/service/kms/client"
	"go.mws.cloud/go-sdk/service/kms/model"
	resourcesdk "go.mws.cloud/go-sdk/service/kms/sdk"

	tfmodel "go.mws.cloud/terraform-provider-mws/service/resources/kms/model"
)

func Read(ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
	sdk *resourcesdk.CryptoKey,
	data tfmodel.CryptoKeyModel) (*model.CryptoKeyOptionalResponse, error) {
	apiRes, err := sdk.GetCryptoKey(
		ctx,
		client.GetCryptoKeyRequest{
			Project: data.ProjectParam.ValueString(),
			Key:     data.KeyParam.ValueString(),
		},
	)
	if err != nil {
		return nil, err
	}
	if apiRes.GetStatus().Destruction != nil {
		// not show this key
		err := errors.New("key is pending destruction")
		resp.Diagnostics.AddError("Key is pending destruction", "Either delete the key from the state or recover it via console or CLI.")
		return nil, err
	}
	return apiRes, nil
}

func Delete(ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
	sdk *resourcesdk.CryptoKey,
	data tfmodel.CryptoKeyModel,
	diags *tfdiag.Diagnostics,
	resourceWaiterTimeout time.Duration) error {
	_, err := sdk.ScheduleDestructionOfCryptoKey(ctx,
		client.ScheduleDestructionOfCryptoKeyRequest{
			Project: data.ProjectParam.ValueString(),
			Key:     data.KeyParam.ValueString(),
		},
	)
	return err
}
