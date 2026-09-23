package electronics

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// SpiceValue parses a number written in SPICE notation and returns it in SI
// units.
//
// Scale factors are case-insensitive, as in SPICE: t (1e12), g (1e9),
// meg (1e6), k (1e3), m (1e-3), mil (25.4e-6), u or µ (1e-6), n (1e-9),
// p (1e-12), f (1e-15) and a (1e-18). Both m and M mean milli; mega is meg.
// Letters following the scale factor are units and are ignored ("10uF",
// "1kOhm"). The 4k7 notation, where the scale factor (or r) takes the place
// of the decimal point, is accepted too: 4k7 = 4700, 2u2 = 2.2e-6, 4r7 = 4.7.
//
// Unlike Value, which follows the BOM convention where M is mega, SpiceValue
// returns an error for anything that is not a SPICE number.
func SpiceValue(s string) (float64, error) {

	in := s
	s = strings.ToLower(strings.TrimSpace(s))

	// Leading number: sign, digits, optional fraction, optional exponent
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	digits := 0
	for i < len(s) && isDigit(s[i]) {
		i++
		digits++
	}
	integer := true
	if i < len(s) && s[i] == '.' {
		integer = false
		i++
		for i < len(s) && isDigit(s[i]) {
			i++
			digits++
		}
	}
	if digits == 0 {
		return math.NaN(), fmt.Errorf("not a SPICE number: %q", in)
	}
	if i < len(s) && s[i] == 'e' {
		j := i + 1
		if j < len(s) && (s[j] == '+' || s[j] == '-') {
			j++
		}
		k := j
		for k < len(s) && isDigit(s[k]) {
			k++
		}
		// An 'e' without exponent digits is not an exponent ("1eV")
		if k > j {
			integer = false
			i = k
		}
	}

	f, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return math.NaN(), fmt.Errorf("not a SPICE number: %q", in)
	}
	rest := s[i:]

	scale, n := spiceScale(rest)
	if n == 0 && len(rest) > 1 && rest[0] == 'r' && isDigit(rest[1]) {
		scale, n = 1, 1
	}
	rest = rest[n:]

	// 4k7 notation: digits right after the scale factor are the fraction
	if n > 0 && integer {
		j := 0
		for j < len(rest) && isDigit(rest[j]) {
			j++
		}
		if j > 0 {
			frac, _ := strconv.ParseFloat("0."+rest[:j], 64)
			if strings.HasPrefix(s, "-") {
				f -= frac
			} else {
				f += frac
			}
			rest = rest[j:]
		}
	}

	// Whatever is left must be a unit
	for _, c := range rest {
		if !(c >= 'a' && c <= 'z' || c == 'Ω' || c == 'ω' || c == '°') {
			return math.NaN(), fmt.Errorf("not a SPICE number: %q", in)
		}
	}

	return f * scale, nil
}

// spiceScale returns the scale factor at the start of s and its length in
// bytes. A length of 0 means no scale factor.
func spiceScale(s string) (float64, int) {

	for _, p := range []struct {
		prefix string
		scale  float64
	}{
		{"meg", 1e6}, {"mil", 25.4e-6},
		{"t", 1e12}, {"g", 1e9}, {"k", 1e3}, {"m", 1e-3},
		{"u", 1e-6}, {"µ", 1e-6}, {"μ", 1e-6},
		{"n", 1e-9}, {"p", 1e-12}, {"f", 1e-15}, {"a", 1e-18},
	} {
		if strings.HasPrefix(s, p.prefix) {
			return p.scale, len(p.prefix)
		}
	}
	return 1, 0
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
