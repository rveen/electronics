package electronics

import (
	"math"
	"strconv"
	"strings"
)

func Values(s string) []float64 {

	ss := strings.Fields(s)
	var ff []float64

	for _, s := range ss {
		ff = append(ff, Value(s))
	}
	return ff
}

func Value(s string) float64 {

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}

	s = strings.ToLower(s)

	var v1 []rune
	var v2 []rune
	var k []rune

	first := true

	for i, c := range s {
		if (c >= '0' && c <= '9') || c == '.' {
			if first {
				v1 = append(v1, c)
			} else {
				v2 = append(v2, c)
			}
		} else if c == ' ' || c == '\t' {
			continue
		} else {
			if i == 0 {
				return math.NaN()
			}
			k = append(k, c)
			first = false
		}
	}

	n1, _ := strconv.ParseFloat(string(v1), 64)

	n2 := 0.0
	if len(v2) != 0 {
		n2, _ = strconv.ParseFloat(string(v2), 64)
	}

	n1 = n1 + n2/10.0

	if len(k) == 0 {
		return n1
	}

	ks := string(k)

	if ks == "meg" {
		return n1 * 1e6
	}

	switch k[0] {
	case 'k':
		return n1 * 1e3
	case 'm':
		// TODO this should return n1*1e-3
		return n1 * 1e6
	case 'u', 'µ':
		return n1 * 1e-6
	case 'n':
		return n1 * 1e-9
	case 'p':
		return n1 * 1e-12
	case 'f':
		return n1 * 1e-15
	default:
		return n1
	}
}
