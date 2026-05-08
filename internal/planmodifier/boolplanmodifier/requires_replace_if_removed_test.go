package boolplanmodifier

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRequiresReplaceIfRemovedModifierPlanModifyBool(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.BoolAttribute{},
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

	testPlan := func(value types.Bool) tfsdk.Plan {
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
	testConfig := func(value types.Bool) tfsdk.Config {
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

	testState := func(value types.Bool) tfsdk.State {
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
		request  planmodifier.BoolRequest
		ifFunc   func() planmodifier.Bool
		expected *planmodifier.BoolResponse
	}{
		"state-null": {
			// resource creation
			request: planmodifier.BoolRequest{
				Plan:       testPlan(types.BoolUnknown()),
				PlanValue:  types.BoolUnknown(),
				State:      nullState,
				StateValue: types.BoolNull(),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.BoolResponse{
				PlanValue:       types.BoolUnknown(),
				RequiresReplace: false,
			},
		},
		"plan-null": {
			// resource destroy
			request: planmodifier.BoolRequest{
				Plan:       nullPlan,
				PlanValue:  types.BoolNull(),
				State:      testState(types.BoolValue(true)),
				StateValue: types.BoolValue(true),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.BoolResponse{
				PlanValue:       types.BoolNull(),
				RequiresReplace: false,
			},
		},
		"value-removed": {
			request: planmodifier.BoolRequest{
				ConfigValue: types.BoolNull(),
				Config:      testConfig(types.BoolNull()),
				Plan:        testPlan(types.BoolNull()),
				PlanValue:   types.BoolNull(),
				State:       testState(types.BoolValue(true)),
				StateValue:  types.BoolValue(true),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.BoolResponse{
				PlanValue:       types.BoolNull(),
				RequiresReplace: true,
			},
		},

		"config-null-plan-notnull": {
			request: planmodifier.BoolRequest{
				ConfigValue: types.BoolNull(),
				Config:      testConfig(types.BoolNull()),
				Plan:        testPlan(types.BoolValue(true)),
				PlanValue:   types.BoolValue(true),
				State:       testState(types.BoolValue(true)),
				StateValue:  types.BoolValue(true),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.BoolResponse{
				PlanValue:       types.BoolValue(true),
				RequiresReplace: false,
			},
		},

		"config-notnull-plan-null": {
			request: planmodifier.BoolRequest{
				Config:      testConfig(types.BoolValue(true)),
				ConfigValue: types.BoolValue(true),
				PlanValue:   types.BoolNull(),
				Plan:        testPlan(types.BoolNull()),
				State:       testState(types.BoolValue(true)),
				StateValue:  types.BoolValue(true),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.BoolResponse{
				PlanValue:       types.BoolNull(),
				RequiresReplace: false,
			},
		},

		"planvalue-statevalue-different-if-false": {
			request: planmodifier.BoolRequest{
				Plan:       testPlan(types.BoolValue(false)),
				PlanValue:  types.BoolValue(false),
				State:      testState(types.BoolValue(true)),
				StateValue: types.BoolValue(true),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.BoolResponse{
				PlanValue:       types.BoolValue(false),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-equal": {
			request: planmodifier.BoolRequest{
				Plan:       testPlan(types.BoolValue(false)),
				PlanValue:  types.BoolValue(false),
				State:      testState(types.BoolValue(false)),
				StateValue: types.BoolValue(false),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.BoolResponse{
				PlanValue:       types.BoolValue(false),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.BoolResponse{
				PlanValue: testCase.request.PlanValue,
			}
			RequiresReplaceIfRemoved().PlanModifyBool(t.Context(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
