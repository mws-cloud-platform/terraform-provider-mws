package listplanmodifier

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRequiresReplaceIfModifierPlanModifyList(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.ListAttribute{
				ElementType: types.StringType,
			},
		},
	}

	nullPlan := tfsdk.Plan{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(t.Context()),
			nil,
		),
	}

	nullState := tfsdk.State{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(t.Context()),
			nil,
		),
	}

	testPlan := func(value types.List) tfsdk.Plan {
		tfValue, err := value.ToTerraformValue(t.Context())

		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Plan{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(t.Context()),
				map[string]tftypes.Value{
					"testattr": tfValue,
				},
			),
		}
	}
	testConfig := func(value types.List) tfsdk.Config {
		tfValue, err := value.ToTerraformValue(t.Context())

		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Config{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(t.Context()),
				map[string]tftypes.Value{
					"testattr": tfValue,
				},
			),
		}
	}
	testState := func(value types.List) tfsdk.State {
		tfValue, err := value.ToTerraformValue(t.Context())

		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.State{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(t.Context()),
				map[string]tftypes.Value{
					"testattr": tfValue,
				},
			),
		}
	}

	testCases := map[string]struct {
		request  planmodifier.ListRequest
		ifFunc   func() planmodifier.List
		expected *planmodifier.ListResponse
	}{
		"state-null": {
			// resource creation
			request: planmodifier.ListRequest{
				Plan:       testPlan(types.ListUnknown(types.StringType)),
				PlanValue:  types.ListUnknown(types.StringType),
				State:      nullState,
				StateValue: types.ListNull(types.StringType),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ListResponse{
				PlanValue:       types.ListUnknown(types.StringType),
				RequiresReplace: false,
			},
		},
		"plan-null": {
			// resource destroy
			request: planmodifier.ListRequest{
				Plan:       nullPlan,
				PlanValue:  types.ListNull(types.StringType),
				State:      testState(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")})),
				StateValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ListResponse{
				PlanValue:       types.ListNull(types.StringType),
				RequiresReplace: false,
			},
		},
		"value-removed": {
			request: planmodifier.ListRequest{
				ConfigValue: types.ListNull(types.StringType),
				Config:      testConfig(types.ListNull(types.StringType)),
				Plan:        testPlan(types.ListNull(types.StringType)),
				PlanValue:   types.ListNull(types.StringType),
				State:       testState(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")})),
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ListResponse{
				PlanValue:       types.ListNull(types.StringType),
				RequiresReplace: true,
			},
		},

		"config-null-plan-notnull": {
			request: planmodifier.ListRequest{
				ConfigValue: types.ListNull(types.StringType),
				Config:      testConfig(types.ListNull(types.StringType)),
				Plan:        testPlan(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")})),
				PlanValue:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")}),
				State:       testState(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")})),
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ListResponse{
				PlanValue:       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")}),
				RequiresReplace: false,
			},
		},

		"config-notnull-plan-null": {
			request: planmodifier.ListRequest{
				Config:      testConfig(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")})),
				ConfigValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")}),
				PlanValue:   types.ListNull(types.StringType),
				Plan:        testPlan(types.ListNull(types.StringType)),
				State:       testState(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")})),
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ListResponse{
				PlanValue:       types.ListNull(types.StringType),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-different-if-false": {
			request: planmodifier.ListRequest{
				Plan:       testPlan(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")})),
				PlanValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")}),
				State:      testState(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")})),
				StateValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ListResponse{
				PlanValue:       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("other")}),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-equal": {
			request: planmodifier.ListRequest{
				Plan:       testPlan(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")})),
				PlanValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
				State:      testState(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")})),
				StateValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ListResponse{
				PlanValue:       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.ListResponse{
				PlanValue: testCase.request.PlanValue,
			}

			RequiresReplaceIfRemoved().PlanModifyList(t.Context(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
