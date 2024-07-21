package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func LogAttributeChangeModifier(f func(context.Context, planmodifier.StringRequest) diag.Diagnostics, description, markdownDescription string) planmodifier.String {
	return logAttributeChangeModifier{
		ifFunc:              f,
		description:         description,
		markdownDescription: markdownDescription,
	}
}

// logAttributeChangeModifier is an plan modifier that logs the attribute change.
type logAttributeChangeModifier struct {
	ifFunc              func(context.Context, planmodifier.StringRequest) diag.Diagnostics
	description         string
	markdownDescription string
}

// Description returns a human-readable description of the plan modifier.
func (m logAttributeChangeModifier) Description(_ context.Context) string {
	return m.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m logAttributeChangeModifier) MarkdownDescription(_ context.Context) string {
	return m.markdownDescription
}

// PlanModifyString implements the plan modification logic.
func (m logAttributeChangeModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do not log on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// Do not log on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Do not log if the plan and state values are equal.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.Append(m.ifFunc(ctx, req)...)
}
