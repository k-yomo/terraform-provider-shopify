// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"shopify": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	mustEnv(t, "SHOPIFY_SHOP")
	mustEnv(t, "SHOPIFY_API_VERSION")
	mustEnv(t, "SHOPIFY_API_KEY")
	mustEnv(t, "SHOPIFY_API_SECRET_KEY")
	mustEnv(t, "SHOPIFY_ADMIN_API_ACCESS_TOKEN")
}

func mustEnv(t *testing.T, name string) {
	t.Helper()
	if os.Getenv(name) == "" {
		t.Fatalf("%s environment variable must be set for acceptance tests", name)
	}
}
