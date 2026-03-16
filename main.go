package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/henryupton/terraform-provider-cron-utils/internal/provider"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/henryupton/cron-utils",
	}
	if err := providerserver.Serve(context.Background(), provider.New("dev"), opts); err != nil {
		log.Fatal(err)
	}
}
