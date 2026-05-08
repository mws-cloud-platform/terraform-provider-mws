package maplanmodifier

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

func TestRequiresReplaceIfModifierPlanModifyMap(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.MapAttribute{
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

	testPlan := func(value types.Map) tfsdk.Plan {
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
	testConfig := func(value types.Map) tfsdk.Config {
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
	testState := func(value types.Map) tfsdk.State {
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
		request  planmodifier.MapRequest
		ifFunc   func() planmodifier.Map
		expected *planmodifier.MapResponse
	}{
		"state-null": {
			// resource creation
			request: planmodifier.MapRequest{
				Plan:       testPlan(types.MapUnknown(types.StringType)),
				PlanValue:  types.MapUnknown(types.StringType),
				State:      nullState,
				StateValue: types.MapNull(types.StringType),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.MapResponse{
				PlanValue:       types.MapUnknown(types.StringType),
				RequiresReplace: false,
			},
		},

		"plan-null": {
			// resource destroy
			request: planmodifier.MapRequest{
				Plan:       nullPlan,
				PlanValue:  types.MapNull(types.StringType),
				State:      testState(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")})),
				StateValue: types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.MapResponse{
				PlanValue:       types.MapNull(types.StringType),
				RequiresReplace: false,
			},
		},
		"value-removed": {
			request: planmodifier.MapRequest{
				ConfigValue: types.MapNull(types.StringType),
				Config:      testConfig(types.MapNull(types.StringType)),
				Plan:        testPlan(types.MapNull(types.StringType)),
				PlanValue:   types.MapNull(types.StringType),
				State:       testState(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")})),
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.MapResponse{
				PlanValue:       types.MapNull(types.StringType),
				RequiresReplace: true,
			},
		},

		"config-null-plan-notnull": {
			request: planmodifier.MapRequest{
				ConfigValue: types.MapNull(types.StringType),
				Config:      testConfig(types.MapNull(types.StringType)),
				Plan:        testPlan(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")})),
				PlanValue:   types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")}),
				State:       testState(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")})),
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.MapResponse{
				PlanValue:       types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")}),
				RequiresReplace: false,
			},
		},

		"config-notnull-plan-null": {
			request: planmodifier.MapRequest{
				Config:      testConfig(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")})),
				ConfigValue: types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")}),
				PlanValue:   types.MapNull(types.StringType),
				Plan:        testPlan(types.MapNull(types.StringType)),
				State:       testState(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")})),
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.MapResponse{
				PlanValue:       types.MapNull(types.StringType),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-different-if-false": {
			request: planmodifier.MapRequest{
				Plan:       testPlan(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")})),
				PlanValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")}),
				State:      testState(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")})),
				StateValue: types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.MapResponse{
				PlanValue:       types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("other")}),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-equal": {
			request: planmodifier.MapRequest{
				Plan:       testPlan(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")})),
				PlanValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")}),
				State:      testState(types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")})),
				StateValue: types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")}),
			},
			ifFunc: RequiresReplaceIfRemoved,
			expected: &planmodifier.MapResponse{
				PlanValue:       types.MapValueMust(types.StringType, map[string]attr.Value{"testkey": types.StringValue("test")}),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.MapResponse{
				PlanValue: testCase.request.PlanValue,
			}

			RequiresReplaceIfRemoved().PlanModifyMap(t.Context(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
