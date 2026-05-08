package provider

import "github.com/hashicorp/terraform-plugin-framework/attr"

func IsValueSet(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}
