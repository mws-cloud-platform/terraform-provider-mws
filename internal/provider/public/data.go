package public

import mwssdk "go.mws.cloud/go-sdk/mws"

// Data is a provider-defined data that is passed to resource and data
// source configuration methods.
type Data struct {
	Config *Config
	SDK    *mwssdk.SDK
}
