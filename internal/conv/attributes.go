package conv

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func GetAttributesTypes[T interface{ GetType() attr.Type }](attrs map[string]T) map[string]attr.Type {
	result := make(map[string]attr.Type, len(attrs))
	for k, v := range attrs {
		result[k] = v.GetType()
	}
	return result
}
