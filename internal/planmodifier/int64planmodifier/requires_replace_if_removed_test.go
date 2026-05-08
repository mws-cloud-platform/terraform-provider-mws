package int64planmodifier

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRequiresReplaceIfRemovedModifierPlanModifyInt64(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.Int64Attribute{},
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

	testPlan := func(value types.Int64) tfsdk.Plan {
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
	testConfig := func(value types.Int64) tfsdk.Config {
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

	testState := func(value types.Int64) tfsdk.State {
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
		request  planmodifier.Int64Request
		ifFunc   func() planmodifier.Int64
		expected *planmodifier.Int64Response
	}{
		"state-null": {
			// resource creation
			request: planmodifier.Int64Request{
				Plan:       testPlan(types.Int64Unknown()),
				PlanValue:  types.Int64Unknown(),
				State:      nullState,
				StateValue: types.Int64Null(),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.Int64Response{
				PlanValue:       types.Int64Unknown(),
				RequiresReplace: false,
			},
		},
		"plan-null": {
			// resource destroy
			request: planmodifier.Int64Request{
				Plan:       nullPlan,
				PlanValue:  types.Int64Null(),
				State:      testState(types.Int64Value(1)),
				StateValue: types.Int64Value(1),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.Int64Response{
				PlanValue:       types.Int64Null(),
				RequiresReplace: false,
			},
		},
		"value-removed": {
			request: planmodifier.Int64Request{
				ConfigValue: types.Int64Null(),
				Config:      testConfig(types.Int64Null()),
				Plan:        testPlan(types.Int64Null()),
				PlanValue:   types.Int64Null(),
				State:       testState(types.Int64Value(1)),
				StateValue:  types.Int64Value(1),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.Int64Response{
				PlanValue:       types.Int64Null(),
				RequiresReplace: true,
			},
		},

		"config-null-plan-notnull": {
			request: planmodifier.Int64Request{
				ConfigValue: types.Int64Null(),
				Config:      testConfig(types.Int64Null()),
				Plan:        testPlan(types.Int64Value(1)),
				PlanValue:   types.Int64Value(1),
				State:       testState(types.Int64Value(1)),
				StateValue:  types.Int64Value(1),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.Int64Response{
				PlanValue:       types.Int64Value(1),
				RequiresReplace: false,
			},
		},

		"config-notnull-plan-null": {
			request: planmodifier.Int64Request{
				Config:      testConfig(types.Int64Value(1)),
				ConfigValue: types.Int64Value(1),
				PlanValue:   types.Int64Null(),
				Plan:        testPlan(types.Int64Null()),
				State:       testState(types.Int64Value(1)),
				StateValue:  types.Int64Value(1),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.Int64Response{
				PlanValue:       types.Int64Null(),
				RequiresReplace: false,
			},
		},

		"planvalue-statevalue-different-if-false": {
			request: planmodifier.Int64Request{
				Plan:       testPlan(types.Int64Value(0)),
				PlanValue:  types.Int64Value(0),
				State:      testState(types.Int64Value(1)),
				StateValue: types.Int64Value(1),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.Int64Response{
				PlanValue:       types.Int64Value(0),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-equal": {
			request: planmodifier.Int64Request{
				Plan:       testPlan(types.Int64Value(0)),
				PlanValue:  types.Int64Value(0),
				State:      testState(types.Int64Value(0)),
				StateValue: types.Int64Value(0),
			},
			ifFunc: RequiresReplaceIfRemoved,

			expected: &planmodifier.Int64Response{
				PlanValue:       types.Int64Value(0),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.Int64Response{
				PlanValue: testCase.request.PlanValue,
			}

			RequiresReplaceIfRemoved().PlanModifyInt64(t.Context(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
