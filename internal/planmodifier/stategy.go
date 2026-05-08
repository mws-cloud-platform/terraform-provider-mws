package planmodifier

import (
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const ErrStrategyNotFound = consterr.Error("strategy not found")

func GetStandardPlanModifierNameByType(p StandardType) (string, error) {
	switch p {
	case StandardRequiresReplaceIfConfigured:
		return "RequiresReplaceIfConfigured()", nil
	case StandardUseStateForUnknown:
		return "UseStateForUnknown()", nil

	default:
		return "", ErrStrategyNotFound
	}
}

func GetPlanModifierNameByType(p CustomType) (string, error) {
	switch p {
	case CustomRequiresReplaceIfRemoved:
		return "RequiresReplaceIfRemoved()", nil
	default:
		return "", ErrStrategyNotFound
	}
}
