package cmp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"go.mws.cloud/terraform-provider-mws/internal/cmp"
)

func TestOr(t *testing.T) {
	null := types.StringNull()
	foo := types.StringValue("foo")
	bar := types.StringValue("bar")

	for _, tc := range []struct {
		vals     []types.String
		expected types.String
	}{
		{[]types.String{}, null},
		{[]types.String{null}, null},
		{[]types.String{foo}, foo},
		{[]types.String{null, foo}, foo},
		{[]types.String{foo, null}, foo},
		{[]types.String{foo, bar}, foo},
		{[]types.String{null, foo, bar}, foo},
	} {
		actual := cmp.Or(tc.vals...)
		require.Equal(t, tc.expected, actual)
	}
}
