package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/robfig/cron/v3"
)

type UnixToQuartzFunction struct{}

func NewUnixToQuartzFunction() function.Function {
	return &UnixToQuartzFunction{}
}

func (f *UnixToQuartzFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "unix_to_quartz"
}

func (f *UnixToQuartzFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Convert a Unix cron expression to Quartz format",
		Description: "Converts a Unix 5-field cron expression to a Quartz 6-field expression by prepending a seconds field of 0.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "expression",
				Description: `A Unix 5-field cron expression (e.g. "*/5 * * * *").`,
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *UnixToQuartzFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var expression string
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &expression))
	if resp.Error != nil {
		return
	}

	if len(strings.Fields(expression)) != 5 {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("expected a Unix expression with 5 fields, got %d", len(strings.Fields(expression))))
		return
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if _, err := parser.Parse(expression); err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("invalid Unix expression %q: %s", expression, err))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, "0 "+expression))
}
