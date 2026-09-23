package electronics

import (
	"math"
	"testing"
)

func TestSpiceValue(t *testing.T) {

	for _, tc := range []struct {
		in   string
		want float64
	}{
		{"0", 0},
		{"10", 10},
		{"-5", -5},
		{"+2.5", 2.5},
		{".5", 0.5},
		{"1e-3", 1e-3},
		{"1E3", 1e3},
		{"1.5e2k", 150e3},
		{"1m", 1e-3},
		{"1M", 1e-3},
		{"1meg", 1e6},
		{"1MEG", 1e6},
		{"1Meg", 1e6},
		{"1mil", 25.4e-6},
		{"10k", 10e3},
		{"10K", 10e3},
		{"1kOhm", 1e3},
		{"10uF", 10e-6},
		{"100n", 100e-9},
		{"22p", 22e-12},
		{"3f", 3e-15},
		{"1a", 1e-18},
		{"2g", 2e9},
		{"1t", 1e12},
		{"4.7µ", 4.7e-6},
		{"4.7μ", 4.7e-6},
		{"5V", 5},
		{"1eV", 1},
		{"10Ω", 10},
		{"4k7", 4700},
		{"4k75", 4750},
		{"2u2", 2.2e-6},
		{"4r7", 4.7},
		{"4R7", 4.7},
		{"10r", 10},
		{"-4k7", -4700},
		{" 1k ", 1e3},
	} {
		got, err := SpiceValue(tc.in)
		if err != nil {
			t.Errorf("SpiceValue(%q): unexpected error %v", tc.in, err)
			continue
		}
		if !close(got, tc.want) {
			t.Errorf("SpiceValue(%q) = %g, want %g", tc.in, got, tc.want)
		}
	}
}

func TestSpiceValueErrors(t *testing.T) {

	for _, in := range []string{"", "k", "abc", "-", ".", "{r1}", "10 k", "1k}", "1,5", "1.5.3"} {
		if v, err := SpiceValue(in); err == nil || !math.IsNaN(v) {
			t.Errorf("SpiceValue(%q) = %g, %v; want NaN and an error", in, v, err)
		}
	}
}

func close(a, b float64) bool {
	if a == b {
		return true
	}
	return math.Abs(a-b) <= 1e-12*math.Max(math.Abs(a), math.Abs(b))
}
