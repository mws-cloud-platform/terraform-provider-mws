package acctest

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// NewEmptyProviderConfig creates a new HCL file with the empty provider
// configuration.
func NewEmptyProviderConfig(providerName string) *hclwrite.File {
	f := hclwrite.NewEmptyFile()
	f.Body().AppendNewBlock("provider", []string{providerName})
	return f
}
