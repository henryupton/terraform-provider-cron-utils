package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &CronProvider{}
var _ provider.ProviderWithFunctions = &CronProvider{}

type CronProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CronProvider{version: version}
	}
}

func (p *CronProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cron"
	resp.Version = p.version
}

func (p *CronProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for parsing cron expressions.",
	}
}

func (p *CronProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *CronProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *CronProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *CronProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewParseFunction,
		NewQuartzToUnixFunction,
		NewUnixToQuartzFunction,
	}
}
