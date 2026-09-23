// Package clearance implements clearance dimensioning per IEC 60664-4.
package isolation

import (
	"errors"
	// "fmt"
	"math"
)

// table1 is IEC 60664-4:2005 Table 1: minimum clearances in air at
// atmospheric pressure for inhomogeneous field conditions.
// u: peak voltage in kV. d: clearance in mm.
var table1 = []struct{ u, d float64 }{
	{0.6, 0.065},
	{0.8, 0.18},
	{1.0, 0.5},
	{1.2, 1.4},
	{1.4, 2.35},
	{1.6, 4.0},
	{1.8, 6.7},
	{2.0, 11.0},
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
