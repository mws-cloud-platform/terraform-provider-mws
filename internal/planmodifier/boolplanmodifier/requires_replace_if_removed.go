package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIfRemoved пересоздает ресурс если:
//
//   - В состоянии ресурса (`StateValue`) уже есть значение.
//   - В плане (`PlanValue`) и конфигурации (`ConfigValue`) значение отсутствует (оба равны `null`).
//
// Это позволяет корректно обрабатывать случаи, когда поле удаляется из
// конфигурации: вместо простого удаления значения Terraform пересоздаёт ресурс,
// чтобы обеспечить согласованность состояния.
func RequiresReplaceIfRemoved() planmodifier.Bool {
	return boolplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.BoolRequest, resp *boolplanmodifier.RequiresReplaceIfFuncResponse) {
			removed := req.PlanValue.IsNull() && req.ConfigValue.IsNull()
			if !req.StateValue.IsNull() && removed {
				resp.RequiresReplace = true
			}
		},
		"If the value of this attribute is removed, Terraform will destroy and recreate the resource.",
		"If the value of this attribute is removed, Terraform will destroy and recreate the resource.",
	)
}
