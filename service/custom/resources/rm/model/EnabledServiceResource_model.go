package model

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnabledServiceModel struct {
	ServiceParam types.String   `tfsdk:"service"`
	ProjectParam types.String   `tfsdk:"project"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}
