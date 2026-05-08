package cmp

import "github.com/hashicorp/terraform-plugin-framework/attr"

// Or returns the first known value from the provided list of values.
func Or[T attr.Value](vals ...T) T {
	var zero T
	for _, val := range vals {
		if !val.IsNull() && !val.IsUnknown() {
			return val
		}
	}
	return zero
}
