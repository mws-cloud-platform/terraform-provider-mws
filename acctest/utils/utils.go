package utils

import (
	"fmt"
	"strings"

	helper "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// RandResourceName returns a random resource name with the given prefix.
func RandResourceName(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, helper.RandString(8))
}

// SliceStringJoin joins strings with quotes around each element.
func SliceStringJoin(in []string) string {
	out := make([]string, len(in))
	for i, z := range in {
		out[i] = "\"" + z + "\""
	}
	return "[" + strings.Join(out, ",") + "]"
}
