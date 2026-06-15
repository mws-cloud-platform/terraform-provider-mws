package acctest

import "go.mws.cloud/util-toolset/pkg/utils/consterr"

const (
	// ErrTerraformInvalidConfig is returned when the terraform configuration is invalid.
	ErrTerraformInvalidConfig = consterr.Error("invalid terraform configuration")
)
