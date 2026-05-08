package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"go.mws.cloud/terraform-provider-mws/service/provider"
)

var version = "devel"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "terraform.local/mws-cloud/mws",
		Debug:   debug,
	}

	tflog.Info(context.Background(), "MainStarted")

	if err := providerserver.Serve(context.Background(), provider.NewProvider(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
