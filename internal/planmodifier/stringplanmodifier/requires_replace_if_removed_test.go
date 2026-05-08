package stringplanmodifier

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRequiresReplaceIfRemovedModifierPlanModifyString(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.StringAttribute{},
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

	testPlan := func(value types.String) tfsdk.Plan {
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
	testConfig := func(value types.String) tfsdk.Config {
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

	testState := func(value types.String) tfsdk.State {
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
		request  planmodifier.StringRequest
		ifFunc   func() planmodifier.String
		expected *planmodifier.StringResponse
	}{
		"state-null": {
			// resource creation
			request: planmodifier.StringRequest{
				Plan:       testPlan(types.StringUnknown()),
				PlanValue:  types.StringUnknown(),
				State:      nullState,
				StateValue: types.StringNull(),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.StringResponse{
				PlanValue:       types.StringUnknown(),
				RequiresReplace: false,
			},
		},
		"plan-null": {
			// resource destroy
			request: planmodifier.StringRequest{
				Plan:       nullPlan,
				PlanValue:  types.StringNull(),
				State:      testState(types.StringValue("test")),
				StateValue: types.StringValue("test"),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.StringResponse{
				PlanValue:       types.StringNull(),
				RequiresReplace: false,
			},
		},
		"value-removed": {
			request: planmodifier.StringRequest{
				ConfigValue: types.StringNull(),
				Config:      testConfig(types.StringNull()),
				Plan:        testPlan(types.StringNull()),
				PlanValue:   types.StringNull(),
				State:       testState(types.StringValue("test")),
				StateValue:  types.StringValue("test"),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.StringResponse{
				PlanValue:       types.StringNull(),
				RequiresReplace: true,
			},
		},

		"config-null-plan-notnull": {
			request: planmodifier.StringRequest{
				ConfigValue: types.StringNull(),
				Config:      testConfig(types.StringNull()),
				Plan:        testPlan(types.StringValue("test")),
				PlanValue:   types.StringValue("test"),
				State:       testState(types.StringValue("test")),
				StateValue:  types.StringValue("test"),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.StringResponse{
				PlanValue:       types.StringValue("test"),
				RequiresReplace: false,
			},
		},

		"config-notnull-plan-null": {
			request: planmodifier.StringRequest{
				Config:      testConfig(types.StringValue("test")),
				ConfigValue: types.StringValue("test"),
				PlanValue:   types.StringNull(),
				Plan:        testPlan(types.StringNull()),
				State:       testState(types.StringValue("test")),
				StateValue:  types.StringValue("test"),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.StringResponse{
				PlanValue:       types.StringNull(),
				RequiresReplace: false,
			},
		},

		"planvalue-statevalue-different-if-false": {
			request: planmodifier.StringRequest{
				Plan:       testPlan(types.StringValue("other")),
				PlanValue:  types.StringValue("other"),
				State:      testState(types.StringValue("test")),
				StateValue: types.StringValue("test"),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.StringResponse{
				PlanValue:       types.StringValue("other"),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-equal": {
			request: planmodifier.StringRequest{
				Plan:       testPlan(types.StringValue("other")),
				PlanValue:  types.StringValue("other"),
				State:      testState(types.StringValue("other")),
				StateValue: types.StringValue("other"),
			},
			ifFunc: RequiresReplaceIfRemoved,

			expected: &planmodifier.StringResponse{
				PlanValue:       types.StringValue("other"),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.StringResponse{
				PlanValue: testCase.request.PlanValue,
			}

			RequiresReplaceIfRemoved().PlanModifyString(t.Context(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
