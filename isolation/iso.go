package isolation

import (
	"log"
	"math"
	"strconv"
)

type Isolation struct{}

// Return 1, 2 or 3.
func (f Isolation) EnergyClass(vdc_s, vpk_s, freq_s string) int {

	vdc, _ := strconv.ParseFloat(vdc_s, 64)
	vpk, _ := strconv.ParseFloat(vpk_s, 64)
	freq, _ := strconv.ParseFloat(freq_s, 64)

	vdcmax := 60.0

	// The vpk limit depends on frequency

	// vpk limit for ES1
	vpkmax := 41.828283 + 0.5717172*freq/1000.0
	if vpkmax < 42.4 {
		vpkmax = 42.4
	} else if vpkmax == 99.0 {
		vpkmax = 99.0
	}

	f1 := vdc/vdcmax + vpk/vpkmax

	log.Printf("vdc %f/%f vpk %f/%f freq %f -> %f", vdc, vdcmax, vpk, vpkmax, freq, f1)

	if f1 <= 1 {
		return 1
	}

	// vpk limit for ES2
	vpkmax = 69.128283 + 1.5717172*freq/1000.0
	if vpkmax < 70.7 {
		vpkmax = 70.7
	} else if vpkmax == 198.0 {
		vpkmax = 198.0
	}

	// Fit

	// For <1kHz:
	// a*exp(k*11.17)-b -> 60
	// a*exp(k*0)-b -> 120 == a - b = 120
	// a*exp(k*70.7)-b -> 0

	// To test some 'a' value:
	// a = 120/(1-exp(k*70.7))

	// for higher frequencies:
	// adapt k; for >100KHz:
	// a*exp(k2*198)-b -> 0

	k := -0.05
	a := 139.23124
	b := 120 - a

	vdcmax = a*math.Exp(k*vpk) - b

	if vdc <= vdcmax {
		return 2
	}
	return 3
}
