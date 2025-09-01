package isolation

/* Based on EVS-EN IEC 60664-1:2020+A1:2025 */

import (
	"golib/mathx"
	"log"
	"math"
)

var riv []float64 = []float64{330, 500, 800, 1500, 2500, 4000, 6000, 8000, 12000, 15000}

func (f Isolation) RatedImpulseVoltage(v float64, ovc int) float64 {
	return RatedImpulseVoltage(v, ovc)
}

func (f Isolation) Clearance(typ string, volt, alt float64, pollution int, reinforced bool) float64 {
	return math.Round(Clearance(typ, volt, alt, pollution, reinforced)*100) / 100
}

func (f Isolation) Creepage(v float64, pollution int, material string, reinforced bool) float64 {
	return Creepage(v, pollution, material, reinforced)
}

func (f Isolation) CreepagePcb(v float64, pollution int, reinforced bool) float64 {
	return CreepagePcb(v, pollution, reinforced)
}

// Table F.1
// v = line to neutral voltage
func RatedImpulseVoltage(v float64, ovc int) float64 {

	if ovc < 1 || ovc > 4 {
		return -1
	}

	if v > 1500 {
		return -1
	}
	if v > 1250 {
		if ovc == 3 {
			return 10000
		}
		return riv[ovc+5]
	}
	if v >= 1000 {
		return riv[ovc+4]
	}
	if v >= 600 {
		return riv[ovc+3]
	}
	if v >= 300 {
		return riv[ovc+2]
	}
	if v >= 150 {
		return riv[ovc+1]
	}
	if v >= 100 {
		return riv[ovc]
	}

	return riv[ovc-1]
}

// Tables F.2, F.8, F.9 (heterogeneous)

var (
	f2v []float64 = []float64{330, 400, 500, 600, 800, 1000, 1200, 1500, 2000, 2500, 3000, 4000, 5000, 6000, 8000, 10000, 12000, 15000, 20000, 25000, 30000, 40000, 50000}
	f2c []float64 = []float64{0.01, 0.02, 0.04, 0.06, 0.10, 0.15, 0.25, 0.5, 1, 1.5, 2, 3, 4, 5.5, 8, 11, 14, 18, 25, 33, 40, 60, 75}
	f8v []float64 = []float64{40, 60, 100, 120, 150, 200, 250, 330, 400, 500, 600, 800, 1000, 1200, 1500, 2000, 2500, 3000, 4000, 5000, 6000, 8000, 10000, 12000, 15000, 20000, 25000, 30000, 40000, 50000}
	f8c []float64 = []float64{0.001, 0.002, 0.003, 0.004, 0.005, 0.006, 0.008, 0.01, 0.02, 0.04, 0.06, 0.13, 0.26, 0.42, 0.76, 1.27, 1.8, 2.4, 3.8, 5.7, 7.9, 11, 15, 2, 19, 25, 34, 44, 55, 77, 100}
	f9c []float64 = []float64{0.001, 0.002, 0.003, 0.004, 0.005, 0.006, 0.008, 0.01, 0.02, 0.04, 0.06, 0.13, 0.26, 0.42, 0.76, 1.27, 2, 3.2, 11, 24, 64, 184, 290, 320, 0, 0, 0, 0, 0, 0}
)

// typ:
// p (pulse, 1.25/50 us) = ambient peaks, h (long, hi-pot) = peaks are part of the signal, d (long, avoid partial discharge)
func Clearance(typ string, volt, altitude float64, pollution int, reinforced bool) float64 {

	log.Printf("typ %s volt %f alt %f pol %d reinf %T\n", typ, volt, altitude, pollution, reinforced)

	if reinforced {
		volt *= 1.6
	}

	if typ == "h" {
		return mathx.Interpolate(volt, f8v, f8c) * AltitudeFactor(altitude)
	}

	if typ == "d" {
		return mathx.Interpolate(volt, f8v, f9c) * AltitudeFactor(altitude)
	}

	if typ != "p" {
		return -1
	}

	// pulse (impulse)
	c := mathx.Interpolate(volt, f2v, f2c)
	if math.IsNaN(c) {
		return c
	}

	switch pollution {
	case 4:
		if c < 1.6 {
			c = 1.6
		}
	case 3:
		if c < 0.8 {
			c = 0.8
		}
	case 2:
		if c < 0.2 {
			c = 0.2
		}
	}

	return c * AltitudeFactor(altitude)
}

// Altitude factor. Tables A2 and F.10
var (
	alt  []float64 = []float64{0, 200, 500, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 15000, 20000}
	kalt []float64 = []float64{0.784, 0.803, 0.833, 0.884, 1, 1.14, 1.29, 1.48, 1.7, 1.95, 2.25, 2.62, 3.02, 6.67, 14.5}
)

func AltitudeFactor(altitude float64) float64 {
	return mathx.Interpolate(altitude, alt, kalt)
}

// IEC 60664-1:2020, 6.2.2.1.4
func TestVoltage(volt, clearance, altitude float64) float64 {

	m := 0.9163 // clearance < 0.01

	switch {
	case clearance > 10:
		m = 0.9243
	case clearance > 1:
		m = 0.8539
	case clearance > 0.0625:
		m = 0.6361
	case clearance > 0.01:
		m = 0.3305
	}

	return volt / math.Pow(1.0/AltitudeFactor(altitude), m)

}

/*
var polyAlt []float64 = []float64{7.7602e-01, 9.5192e-05, 8.1820e-09, -9.5830e-14, 5.8214e-17}

func AltituteFactorPoly(altitude float64) float64 {
	return mathx.PolyVal(altitude, polyAlt)
}*/

var (
	x_vcr  []float64 = []float64{10, 12.5, 16, 20, 25, 32, 40, 50, 63, 80, 100, 125, 160, 200, 250, 320, 400, 500, 630, 800, 1000, 1250, 1600, 2000, 2500, 3200, 4000, 5000, 6000, 8000, 10000}
	y_pcb1 []float64 = []float64{0.025, 0.025, 0.025, 0.025, 0.025, 0.025, 0.025, 0.025, 0.04, 0.063, 0.1, 0.16, 0.25, 0.4, 0.56, 0.75, 1, 1.3, 1.8, 2.4, 3.2}
	y_pcb2 []float64 = []float64{0.04, 0.04, 0.04, 0.04, 0.04, 0.04, 0.04, 0.04, 0.063, 0.1, 0.16, 0.25, 0.4, 0.63, 1, 1.6, 2, 2.5, 3.2, 4, 5}
	y_p1   []float64 = []float64{0.08, 0.09, 0.1, 0.11, 0.125, 0.14, 0.16, 0.18, 0.2, 0.22, 0.25, 0.28, 0.32, 0.42, 0.56, 0.75, 1, 1.3, 1.8, 2.4, 3.2, 4.2, 5.6, 7.5, 10, 12.5, 16, 20, 32, 40}

	y_p2m1 []float64 = []float64{0.4, 0.42, 0.45, 0.48, 0.5, 0.53, 0.56, 0.6, 0.63, 0.67, 0.71, 0.75, 0.8, 1, 1.25, 1.6, 2, 2.5, 3.2, 4, 5, 6.3, 8, 10, 12.5, 16, 20, 25, 32, 40, 50}
	y_p2m2 []float64 = []float64{0.4, 0.42, 0.45, 0.48, 0.5, 0.53, 0.8, 0.85, 0.9, 1, 1.05, 1.1, 1.4, 1.8, 2.2, 2.8, 3.6, 4.5, 5.6, 7.1, 9, 11, 14, 18, 22, 28, 36, 45, 56, 71}
	y_p2m3 []float64 = []float64{0.4, 0.42, 0.45, 0.48, 0.5, 0.53, 1.1, 1.2, 1.25, 1.3, 1.4, 1.5, 1.6, 2, 2.5, 3.2, 4, 5, 6.3, 8, 10, 12.5, 16, 20, 25, 32, 40, 50, 63, 80, 100}

	y_p3m1 []float64 = []float64{1, 1.05, 1.1, 1.2, 1.25, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 2, 2.5, 3.2, 4, 5, 6.3, 8, 10, 12.5, 16, 20, 25, 32, 40, 50, 63, 80, 100, 125}
	y_p3m2 []float64 = []float64{1, 1.05, 1.1, 1.2, 1.25, 1.3, 1.6, 1.7, 1.8, 1.9, 2, 2.1, 2.2, 2.8, 3.6, 4.5, 5.6, 7.1, 9, 11, 14, 18, 22, 28, 36, 45, 56, 71, 90, 110, 140}
	y_p3m3 []float64 = []float64{1, 1.05, 1.1, 1.2, 1.25, 1.3, 1.8, 1.9, 2, 2.1, 2.2, 2.4, 2.5, 3.2, 4, 5, 6.3, 8, 10, 12.5, 16, 20, 25, 32, 40, 50, 63, 80, 100, 125, 160}
)

func CreepagePcb(v float64, pollution int, reinforced bool) float64 {

	k := 1.0
	if reinforced {
		k = 2.0
	}

	switch pollution {
	case 1:
		return mathx.Interpolate(v, x_vcr, y_pcb1) * k
	case 2:
		return mathx.Interpolate(v, x_vcr, y_pcb2) * k
	}

	return math.NaN()
}

func Creepage(v float64, pollution int, material string, reinforced bool) float64 {

	k := 1.0
	if reinforced {
		k = 2.0
	}

	switch pollution {
	case 1:
		return mathx.Interpolate(v, x_vcr, y_p1) * k
	case 2:
		switch material {
		case "1":
			return mathx.Interpolate(v, x_vcr, y_p2m1) * k
		case "2":
			return mathx.Interpolate(v, x_vcr, y_p2m2) * k
		case "3a", "3b", "p2", "p3":
			return mathx.Interpolate(v, x_vcr, y_p2m3) * k
		}
	case 3:
		switch material {
		case "1":
			return mathx.Interpolate(v, x_vcr, y_p3m1) * k
		case "2":
			return mathx.Interpolate(v, x_vcr, y_p3m2) * k
		case "3a", "3b", "p2", "p3":
			return mathx.Interpolate(v, x_vcr, y_p3m3) * k
		}
	}

	return math.NaN()
}
