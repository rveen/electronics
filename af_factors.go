package electronics

import (
	"math"
)

// Acceleration factors

const (
	InvBoltzmann float64 = 11604.522
	// Boltzmann constant in eV/K
	BoltzmannEv float64 = 8.617333262e-5
	k                   = BoltzmannEv
	t0                  = 273.15
)

// Arrhenius model
func AF_Arrhenius(ea, ttest, tfield float64) float64 {
	return math.Exp(InvBoltzmann * ea * (1/(ttest+t0) - 1/(tfield+t0)))
}

// Lawson model
//
// typical ea = 0.4
func AF_Lawson(ea, ttest, rhtest, tfield, rhfield float64) float64 {

	const b = 5.57e-4

	r1 := b * (math.Pow(rhtest, 2) - math.Pow(rhfield, 2))
	r2 := -ea / k * (1.0/(ttest+t0) - 1.0/(tfield+t0))

	return math.Exp(r2 + r1)
}

// Coffin-Manson acceleration factor
func AF_CoffinManson(ttest, tfield float64) float64 {
	return math.Pow(ttest/tfield, 2.5)
}

// Modified Norris-Landzberg model (Lead free, SAC305)
// Source ZVEI Robust validation of EE modules
// _t : test, _f : field
// Constants taken from: https://www.researchgate.net/publication/275951660_Norris-Landzberg_Acceleration_Factors_and_Goldmann_Constants_for_SAC305_Lead-Free_Electronics
// Typical a = 2.65, b = 0.136, ea = 0.4
//
// Check https://watermark.silverchair.com/jom_v33_1_35.pdf
func AF_NorrisLandzberg(dt_f, dt_t, f_f, f_t, tmax_f, tmax_t float64, b, y, ea float64) float64 {
	return math.Pow(dt_t/dt_f, b) * math.Pow(f_f/f_t, y) * math.Exp(InvBoltzmann*ea*(1/(tmax_f+t0)-1/(tmax_t+t0)))
}

// Basquin's model
func AF_Basquin(grms float64) float64 {
	return math.Pow(grms*2, 1.5)
}

// Peck’s model
func AF_Peck(ea, rh, ttest, tfield float64) float64 {
	return math.Pow(rh/70, 4.4) * AF_Arrhenius(ea, ttest, tfield)
}
