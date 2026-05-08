package objectplanmodifier

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

func TestRequiresReplaceIfModifierPlanModifyObject(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{"testattr": types.StringType},
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

	testPlan := func(value types.Object) tfsdk.Plan {
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
	testConfig := func(value types.Object) tfsdk.Config {
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
	testState := func(value types.Object) tfsdk.State {
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
		request  planmodifier.ObjectRequest
		ifFunc   func() planmodifier.Object
		expected *planmodifier.ObjectResponse
	}{
		"state-null": {
			// resource creation
			request: planmodifier.ObjectRequest{
				Plan:       testPlan(types.ObjectUnknown(map[string]attr.Type{"testattr": types.StringType})),
				PlanValue:  types.ObjectUnknown(map[string]attr.Type{"testattr": types.StringType}),
				State:      nullState,
				StateValue: types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectUnknown(map[string]attr.Type{"testattr": types.StringType}),
				RequiresReplace: false,
			},
		},
		"plan-null": {
			// resource destroy
			request: planmodifier.ObjectRequest{
				Plan:       nullPlan,
				PlanValue:  types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				State:      testState(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")})),
				StateValue: types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				RequiresReplace: false,
			},
		},
		"value-removed": {
			request: planmodifier.ObjectRequest{
				ConfigValue: types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				Config:      testConfig(types.ObjectNull(map[string]attr.Type{"testattr": types.StringType})),
				Plan:        testPlan(types.ObjectNull(map[string]attr.Type{"testattr": types.StringType})),
				PlanValue:   types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				State:       testState(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")})),
				StateValue:  types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				RequiresReplace: true,
			},
		},

		"config-null-plan-notnull": {
			request: planmodifier.ObjectRequest{
				ConfigValue: types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				Config:      testConfig(types.ObjectNull(map[string]attr.Type{"testattr": types.StringType})),
				Plan:        testPlan(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")})),
				PlanValue:   types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
				State:       testState(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")})),
				StateValue:  types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
				RequiresReplace: false,
			},
		},

		"config-notnull-plan-null": {
			request: planmodifier.ObjectRequest{
				Config:      testConfig(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")})),
				ConfigValue: types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
				PlanValue:   types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				Plan:        testPlan(types.ObjectNull(map[string]attr.Type{"testattr": types.StringType})),
				State:       testState(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")})),
				StateValue:  types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-different-if-false": {
			request: planmodifier.ObjectRequest{
				Plan:       testPlan(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")})),
				PlanValue:  types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
				State:      testState(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")})),
				StateValue: types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("other")}),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-equal": {
			request: planmodifier.ObjectRequest{
				Plan:       testPlan(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")})),
				PlanValue:  types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")}),
				State:      testState(types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")})),
				StateValue: types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectValueMust(map[string]attr.Type{"testattr": types.StringType}, map[string]attr.Value{"testattr": types.StringValue("test")}),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.ObjectResponse{
				PlanValue: testCase.request.PlanValue,
			}

			RequiresReplaceIfRemoved().PlanModifyObject(t.Context(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
