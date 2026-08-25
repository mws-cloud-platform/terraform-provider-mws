package conv

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// MustKnownValue возвращает заполненное значение для заданного типа
// устанавливая zero value. Паникует при некорректном типе атрибута.
func MustKnownValue(t attr.Type) attr.Value {
	switch t := t.(type) {
	case basetypes.StringType:
		return types.StringValue("")
	case basetypes.BoolType:
		return types.BoolValue(false)
	case basetypes.Int64Type:
		return types.Int64Value(0)
	case basetypes.Float64Type:
		return types.Float64Value(0)
	case basetypes.ListType:
		return types.ListValueMust(t.ElemType, []attr.Value{})
	case basetypes.SetType:
		return types.SetValueMust(t.ElemType, []attr.Value{})
	case basetypes.MapType:
		return types.MapValueMust(t.ElemType, map[string]attr.Value{})
	case basetypes.ObjectType:
		return MustKnownObjectValue(t.AttrTypes)
	default:
		panic(fmt.Sprintf("unsupported attr type %T", t))
	}
}

// MustKnownObjectValue возвращает заполненное значение объекта для схемы tf-модели.
// Паникует при некорректном типе атрибута.
func MustKnownObjectValue(attrTypes map[string]attr.Type) basetypes.ObjectValue {
	return types.ObjectValueMust(attrTypes, knownAttrs(attrTypes))
}

func knownAttrs(attrTypes map[string]attr.Type) map[string]attr.Value {
	attrs := make(map[string]attr.Value, len(attrTypes))
	for name, attrType := range attrTypes {
		attrs[name] = MustKnownValue(attrType)
	}
	return attrs
}
