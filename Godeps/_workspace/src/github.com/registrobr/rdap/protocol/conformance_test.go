package protocol

import (
	"reflect"
	"testing"
)

func TestConformanceSetConformance(t *testing.T) {
	expected := []string{"rdap_level0", "nicbr_level0"}

	var c Conformance
	c.SetConformance(expected)

	if !reflect.DeepEqual(c.Levels, expected) {
		t.Errorf("Unexpected conformance levels. Expected “%#v” and got “%#v”", expected, c.Levels)
	}
}
