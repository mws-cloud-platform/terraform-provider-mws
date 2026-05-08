package public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	mwssdk "go.mws.cloud/go-sdk/mws"
	"go.mws.cloud/go-sdk/mws/credentials"
	"go.mws.cloud/go-sdk/mws/iam"
	"go.mws.cloud/util-toolset/pkg/os/env"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.mws.cloud/terraform-provider-mws/internal/imds"
	base "go.mws.cloud/terraform-provider-mws/internal/provider"
	"go.mws.cloud/terraform-provider-mws/internal/useragent"
)

var errMissingCredentials = fmt.Errorf("one of %q or %q must be specified. "+
	"Alternatively, you can run this provider inside a Compute VM with a linked service account",
	base.TFMWSTokenVar, base.TFServiceAccountAuthorizedKeyPathVar)

// LoadSDKFromConfig loads SDK from the provider configuration.
func LoadSDKFromConfig(ctx context.Context, config *base.Config, env env.Env, version string) (*mwssdk.SDK, error) {
	opts := []mwssdk.LoadSDKOption{
		mwssdk.WithDefaultBaseEndpoint(config.Endpoint.ValueString()),
		mwssdk.WithDefaultProject(config.Project.ValueString()),
		mwssdk.WithDefaultZone(config.Zone.ValueString()),
		mwssdk.WithTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))),
		mwssdk.WithUserAgent(useragent.New(env, "mws-terraform", version)),
	}

	switch {
	case base.IsValueSet(config.ServiceAccountAuthorizedKeyPath):
		path := config.ServiceAccountAuthorizedKeyPath.ValueString()
		saAuthorizedKey, err := loadServiceAccountAuthorizedKeyFile(path)
		if err != nil {
			return nil, fmt.Errorf("load service account authorized key file: %w", err)
		}
		opts = append(opts, mwssdk.WithServiceAccountAuthorizedKey(saAuthorizedKey))
	case base.IsValueSet(config.MWSToken):
		opts = append(opts, mwssdk.WithCredentials(credentials.StaticProvider(credentials.Credentials{
			AccessToken: config.MWSToken.ValueString(),
		})))
	case onComputeVMWithSA(ctx, env):
		// do nothing, since there is no explicit option for setting VM SA credentials provider.
		// but, this provider will automatically be used, because Compute VM environment is detected.
	default:
		return nil, errMissingCredentials
	}

	sdk, err := mwssdk.Load(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("load sdk: %w", err)
	}

	return sdk, nil
}

func loadServiceAccountAuthorizedKeyFile(filePath string) (iam.ServiceAccountAuthorizedKey, error) {
	var serviceAccountAuthorizedKey iam.ServiceAccountAuthorizedKey

	data, err := os.ReadFile(filePath)
	if err != nil {
		return iam.ServiceAccountAuthorizedKey{}, fmt.Errorf("read service account authorized key file: %w", err)
	}

	if err = json.Unmarshal(data, &serviceAccountAuthorizedKey); err != nil {
		return iam.ServiceAccountAuthorizedKey{}, fmt.Errorf("unmarshal service account authorized key: %w", err)
	}

	return serviceAccountAuthorizedKey, nil
}

func onComputeVMWithSA(ctx context.Context, env env.Env) bool {
	serviceAccountRef, err := imds.GetVMServiceAccount(ctx, http.DefaultClient, env)
	if err == nil && serviceAccountRef != "" {
		tflog.Info(ctx, "compute VM with SA detected", map[string]any{"sa": serviceAccountRef})
		return true
	}

	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		tflog.Warn(ctx, "failed to get vm service account", map[string]any{"error": err.Error()})
	}
	return false
}
