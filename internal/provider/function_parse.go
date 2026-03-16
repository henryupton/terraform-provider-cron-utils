package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/reugn/go-quartz/quartz"
	"github.com/robfig/cron/v3"
)

const nextExecutionCount = 5

// regularity check uses extra executions beyond what we return to the caller
const regularityCheckCount = nextExecutionCount + 8

var parseReturnAttrTypes = map[string]attr.Type{
	"expression_type":     types.StringType,
	"next_execution":      types.StringType,
	"next_execution_unix": types.Int64Type,
	"is_regular":          types.BoolType,
	"interval_seconds":    types.Int64Type,
	"next_executions":     types.ListType{ElemType: types.StringType},
}

// cronSchedule is a unified interface over robfig and go-quartz schedules.
type cronSchedule interface {
	Next(t time.Time) (time.Time, error)
}

type robfigAdapter struct{ s cron.Schedule }

func (a *robfigAdapter) Next(t time.Time) (time.Time, error) {
	return a.s.Next(t), nil
}

type quartzAdapter struct{ t *quartz.CronTrigger }

func (a *quartzAdapter) Next(t time.Time) (time.Time, error) {
	nanos, err := a.t.NextFireTime(t.UnixNano())
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, nanos), nil
}

type ParseFunction struct{}

func NewParseFunction() function.Function {
	return &ParseFunction{}
}

func (f *ParseFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parse"
}

func (f *ParseFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Parse a cron expression",
		Description: "Parses a Unix (5-field) or Quartz (6-7 field) cron expression and returns schedule metadata. Format is auto-detected by field count.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "expression",
				Description: `A Unix 5-field (e.g. "*/5 * * * *") or Quartz 6-7 field cron expression (e.g. "0 */5 * * * ?").`,
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: parseReturnAttrTypes,
		},
	}
}

func (f *ParseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var expression string
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &expression))
	if resp.Error != nil {
		return
	}

	schedule, exprType, err := detectAndParse(expression)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("invalid cron expression %q: %s", expression, err))
		return
	}

	now := time.Now()
	allTimes, err := computeNextExecutions(schedule, now, regularityCheckCount)
	if err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("error computing next executions: %s", err))
		return
	}

	isRegular, intervalSeconds := checkRegularity(allTimes)

	nextExecutionVals := make([]attr.Value, nextExecutionCount)
	for i := 0; i < nextExecutionCount; i++ {
		nextExecutionVals[i] = types.StringValue(allTimes[i].Format(time.RFC3339))
	}

	nextExecutionsList, diags := types.ListValue(types.StringType, nextExecutionVals)
	if diags.HasError() {
		resp.Error = function.NewFuncError("internal error constructing next_executions list")
		return
	}

	result, diags := types.ObjectValue(parseReturnAttrTypes, map[string]attr.Value{
		"expression_type":     types.StringValue(exprType),
		"next_execution":      types.StringValue(allTimes[0].Format(time.RFC3339)),
		"next_execution_unix": types.Int64Value(allTimes[0].Unix()),
		"is_regular":          types.BoolValue(isRegular),
		"interval_seconds":    types.Int64Value(intervalSeconds),
		"next_executions":     nextExecutionsList,
	})
	if diags.HasError() {
		resp.Error = function.NewFuncError("internal error constructing result object")
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}

// detectAndParse auto-detects cron format by field count and returns a unified cronSchedule.
// 5 fields → Unix (robfig/cron); 6-7 fields → Quartz (reugn/go-quartz, year field dropped if present).
func detectAndParse(expression string) (cronSchedule, string, error) {
	fields := strings.Fields(expression)
	switch len(fields) {
	case 5:
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		s, err := parser.Parse(expression)
		if err != nil {
			return nil, "", err
		}
		return &robfigAdapter{s}, "unix", nil
	case 6, 7:
		expr := expression
		if len(fields) == 7 {
			expr = strings.Join(fields[:6], " ")
		}
		t, err := quartz.NewCronTrigger(expr)
		if err != nil {
			return nil, "", err
		}
		return &quartzAdapter{t}, "quartz", nil
	default:
		return nil, "", fmt.Errorf("expected 5 fields (unix) or 6-7 fields (quartz), got %d", len(fields))
	}
}

func computeNextExecutions(schedule cronSchedule, from time.Time, count int) ([]time.Time, error) {
	times := make([]time.Time, count)
	t := from
	for i := range times {
		next, err := schedule.Next(t)
		if err != nil {
			return nil, err
		}
		times[i] = next
		t = next
	}
	return times, nil
}

// checkRegularity returns whether all intervals between consecutive executions are equal.
// Returns (true, intervalSeconds) if regular, (false, 0) otherwise.
// Examples: "*/5 * * * *" → regular 300s; "0 9 1 * *" → irregular (months differ in length).
func checkRegularity(times []time.Time) (bool, int64) {
	if len(times) < 2 {
		return false, 0
	}
	first := int64(times[1].Sub(times[0]).Seconds())
	for i := 1; i < len(times)-1; i++ {
		if int64(times[i+1].Sub(times[i]).Seconds()) != first {
			return false, 0
		}
	}
	return true, first
}
