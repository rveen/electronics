// Package clearance implements clearance dimensioning per IEC 60664-4.
package isolation

import (
	"errors"
	"log"
	"math"
)

// table1 is IEC 60664-4:2005 Table 1: minimum clearances in air at
// atmospheric pressure for inhomogeneous field conditions.
// u: peak voltage in kV. d: clearance in mm.
var table1 = []struct{ u, d float64 }{
	{600, 0.065},
	{800, 0.18},
	{1000, 0.5},
	{1200, 1.4},
	{1400, 2.35},
	{1600, 4.0},
	{1800, 6.7},
	{2000, 11.0},
}

// ErrAboveTable is returned for peak voltages above the last row of Table 1.
var ErrAboveTable = errors.New("iec60664-4: U_peak above 2.0 kV is not covered by Table 1")

// InhomogeneousClearance returns the minimum clearance in mm for the peak
// voltage uPeak in kV, per IEC 60664-4 Table 1 (clause 4.4.3).
//
// Up to 0.6 kV the table value 0.065 mm applies (footnote b: no data below
// 0.6 kV). Between rows the result is interpolated linearly (footnote a).
func ClearanceInhomogenuousHF(uPeak float64) float64 {

	if math.IsNaN(uPeak) || uPeak < 0 {
		return math.NaN()
		// return 0, fmt.Errorf("iec60664-4: invalid U_peak %v kV", uPeak)
	}
	if uPeak <= table1[0].u {
		return table1[0].d
	}
	for i := 1; i < len(table1); i++ {
		p, q := table1[i-1], table1[i]
		if uPeak <= q.u {
			return p.d + (uPeak-p.u)*(q.d-p.d)/(q.u-p.u)
		}
	}
	return math.NaN()
	// return 0, fmt.Errorf("%w (got %.3f kV)", ErrAboveTable, uPeak)
}

/* --------------------------------------------------------------------------- */

var nan = math.NaN()

// t2Freq holds the upper frequency limit of each column of IEC 60664-4:2005
// Table 2, in Hz. The first column covers 30 kHz < f <= 100 kHz.
var t2Freq = []float64{0.1e6, 0.2e6, 0.4e6, 0.7e6, 1e6, 2e6, 3e6}

// t2Volt holds the row voltages of Table 2, U_peak in V.
var t2Volt = []float64{100, 200, 300, 400, 500, 600, 700, 800, 900,
	1000, 1100, 1200, 1300, 1400, 1500, 1600, 1700, 1800}

// t2 holds the creepage distances in mm for pollution degree 1.
// NaN marks an empty cell.
var t2 = [][]float64{
	{0.0167, nan, nan, nan, nan, nan, 0.3},
	{0.042, nan, nan, nan, nan, 0.15, 2.8},
	{0.083, 0.09, 0.09, 0.09, 0.09, 0.8, 20},
	{0.125, 0.13, 0.15, 0.19, 0.35, 4.5, nan},
	{0.183, 0.19, 0.25, 0.4, 1.5, 20, nan},
	{0.267, 0.27, 0.4, 0.85, 5, nan, nan},
	{0.358, 0.38, 0.68, 1.9, 20, nan, nan},
	{0.45, 0.55, 1.1, 3.8, nan, nan, nan},
	{0.525, 0.82, 1.9, 8.7, nan, nan, nan},
	{0.6, 1.15, 3, 18, nan, nan, nan},
	{0.683, 1.7, 5, nan, nan, nan, nan},
	{0.85, 2.4, 8.2, nan, nan, nan, nan},
	{1.2, 3.5, nan, nan, nan, nan, nan},
	{1.65, 5, nan, nan, nan, nan, nan},
	{2.3, 7.3, nan, nan, nan, nan, nan},
	{3.15, nan, nan, nan, nan, nan, nan},
	{4.4, nan, nan, nan, nan, nan, nan},
	{6.1, nan, nan, nan, nan, nan, nan},
}

// CreepageHF returns the minimum creepage distance in mm per
// IEC 60664-4:2005 Table 2.
//
//	volt:            U_peak in V (up to 1800)
//	freq:            frequency in Hz (above 30 kHz, up to 3 MHz)
//	pollutionDegree: "1", "2" or "3" (a "PD" prefix is accepted)
//
// A voltage below 100 V uses the 100 V row. A voltage between rows is
// rounded up to the next row; the table does not
// permit interpolation between rows. Values between columns are
// interpolated linearly (footnote b). The result is NaN for inputs outside the table, for cells beyond the
// last filled column of a row, and for unknown pollution degrees.
func CreepageHF(volt, freq float64, pollutionDegree int, reinforced bool) float64 {

	k := pollutionFactor(pollutionDegree)
	if math.IsNaN(k) || math.IsNaN(volt) || math.IsNaN(freq) {
		return nan
	}
	if volt < 0 || volt > t2Volt[len(t2Volt)-1] {
		return nan
	}
	if freq <= 30e3 || freq > t2Freq[len(t2Freq)-1] {
		return nan
	}

	// Round up to the first row with t2Volt[i] >= volt.
	i := 0
	for t2Volt[i] < volt {
		i++
	}

	k2 := 1.0
	if reinforced {
		k2 = 2.0
	}

	return k * k2 * rowAt(i, freq)
}

// pollutionFactor returns the multiplication factor of Table 2, footnote a.
func pollutionFactor(pd int) float64 {
	switch pd {
	case 1:
		return 1.0
	case 2:
		return 1.2
	case 3:
		return 1.4
	}
	return 0
}

// rowAt returns the creepage distance of row i at frequency f (Hz),
// interpolating between the nearest filled cells of that row.
// It returns NaN if no filled cell exists at or above f.
func rowAt(i int, f float64) float64 {
	row := t2[i]
	if f <= t2Freq[0] {
		return row[0]
	}
	// j: first column with t2Freq[j] >= f.
	j := 1
	for t2Freq[j] < f {
		j++
	}
	right := j
	for right < len(row) && math.IsNaN(row[right]) {
		right++
	}
	if right == len(row) {
		return nan
	}
	left := j - 1
	for math.IsNaN(row[left]) {
		left--
	}
	return lerp(f, t2Freq[left], t2Freq[right], row[left], row[right])
}

// lerp interpolates linearly; NaN in either endpoint yields NaN.
func lerp(x, x0, x1, y0, y1 float64) float64 {
	if x == x1 {
		return y1
	}
	return y0 + (x-x0)*(y1-y0)/(x1-x0)
}
