package planmodifier

type CustomType int8
type StandardType int8

const (
	StandardRequiresReplaceIfConfigured StandardType = iota
	StandardUseStateForUnknown
)

const (
	CustomRequiresReplaceIfRemoved CustomType = iota
)
