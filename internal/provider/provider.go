// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ShopifyProvider satisfies various provider interfaces.
var _ provider.Provider = &ShopifyProvider{}
var _ provider.ProviderWithFunctions = &ShopifyProvider{}

// ShopifyProvider defines the provider implementation.
type ShopifyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ShopifyProviderModel describes the provider data model.
type ShopifyProviderModel struct {
	APIKey              types.String `tfsdk:"api_key"`
	APISecretKey        types.String `tfsdk:"api_secret_key"`
	AdminAPIAccessToken types.String `tfsdk:"admin_api_access_token"`
}

func (p *ShopifyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "shopify"
	resp.Version = p.version
}

func (p *ShopifyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Shopify app API key.",
				Required:            true,
			},
			"api_secret_key": schema.StringAttribute{
				MarkdownDescription: "Shopify app API secret key.",
				Required:            true,
				Sensitive:           true,
			},
			"admin_api_access_token": schema.StringAttribute{
				MarkdownDescription: "Shopify Admin API access token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ShopifyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ShopifyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ShopifyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *ShopifyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func (p *ShopifyProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ShopifyProvider{
			version: version,
		}
	}
}
