package acctest

import (
	"testing"

	"go.mws.cloud/util-toolset/pkg/testing/golden"
)

func TestNewEmptyProviderConfig(t *testing.T) {
	dir := golden.NewDir(t,
		golden.WithPath("testdata/"+t.Name()),
		golden.WithRecreateOnUpdate(),
	)
	config := NewEmptyProviderConfig("mws")
	dir.Bytes(t, "provider.tf", config.Bytes())
}
