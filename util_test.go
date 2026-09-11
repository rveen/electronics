package electronics

import (
	"math"
	"testing"
)

// Value follows the BOM convention: M (and m) is mega. See SpiceValue for
// SPICE semantics.
func TestValue(t *testing.T) {

	for _, tc := range []struct {
		in   string
		want float64
	}{
		{"10", 10},
		{"1e-3", 1e-3},
		{"10k", 10e3},
		{"1M", 1e6},
		{"1meg", 1e6},
		{"100u", 100e-6},
		{"10n", 10e-9},
		{"22p", 22e-12},
		{"5%", 5},
		{"4k7", 4700},
		{"4k75", 4750},
		{"2u2", 2.2e-6},
		{"10 k", 10e3},
	} {
		if got := Value(tc.in); !close(got, tc.want) {
			t.Errorf("Value(%q) = %g, want %g", tc.in, got, tc.want)
		}
	}

	if v := Value("abc"); !math.IsNaN(v) {
		t.Errorf("Value(%q) = %g, want NaN", "abc", v)
	}
}
