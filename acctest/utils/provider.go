package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"go.mws.cloud/terraform-provider-mws/service/provider"
)

func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"mws": providerserver.NewProtocol6WithError(provider.NewProvider("test")()),
	}
}
