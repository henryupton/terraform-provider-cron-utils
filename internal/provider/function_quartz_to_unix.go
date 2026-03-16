package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/reugn/go-quartz/quartz"
)

type QuartzToUnixFunction struct{}

func NewQuartzToUnixFunction() function.Function {
	return &QuartzToUnixFunction{}
}

func (f *QuartzToUnixFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "quartz_to_unix"
}

func (f *QuartzToUnixFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Convert a Quartz cron expression to Unix format",
		Description: "Converts a Quartz 6-7 field cron expression to a Unix 5-field expression by dropping the seconds and year fields. Returns an error if the expression uses Quartz-specific syntax (L, W, #) that has no Unix equivalent.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "expression",
				Description: `A Quartz 6-7 field cron expression (e.g. "0 */5 * * * ?").`,
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *QuartzToUnixFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var expression string
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &expression))
	if resp.Error != nil {
		return
	}

	fields := strings.Fields(expression)
	if len(fields) < 6 || len(fields) > 7 {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("expected a Quartz expression with 6-7 fields, got %d", len(fields)))
		return
	}

	// Validate the expression parses correctly before converting
	expr := expression
	if len(fields) == 7 {
		expr = strings.Join(fields[:6], " ")
	}
	if _, err := quartz.NewCronTrigger(expr); err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("invalid Quartz expression %q: %s", expression, err))
		return
	}

	// Unix fields are: min hr dom mon dow (Quartz fields 1-5, 0-indexed)
	unixFields := fields[1:6]

	for _, field := range unixFields {
		for _, special := range []string{"L", "W", "#"} {
			if strings.ContainsAny(field, special) {
				resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf(
					"cannot convert to Unix cron: field %q uses Quartz-specific syntax %q which has no Unix equivalent",
					field, special,
				))
				return
			}
		}
	}

	// Replace Quartz's "no specific value" marker with Unix's wildcard
	for i, field := range unixFields {
		if field == "?" {
			unixFields[i] = "*"
		}
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, strings.Join(unixFields, " ")))
}
