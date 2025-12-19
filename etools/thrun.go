package etools

import (
	"fmt"
	"math"

	"github.com/rveen/electronics"
	m "github.com/rveen/margaid"
)

type Etools struct{}

func (E Etools) ThRunaway(i0, vr, t0 float64) string {

	fmt.Printf("ThRunaway %f %f %f\n", i0, vr, t0)

	// Current leakage at 25ºC
	// i0 := 50e-6

	// Reverse voltage in actual circuit
	// vr := 35.0

	// Temperature dependency of leakage current
	// We assume 2^(T/10) = e^(T/k)
	k := 10 / math.Log(2.05)

	// Ambient temperature (sink or temperature ground)
	// t0 := 75.0
	tmax := 150.0

	// Thermal resistance junction to ambient
	rth := 300.0 // ºC/W
	ok := false

	for n := 0; n < 100; n++ {
		rth -= 5
		ok = !Runaway(vr, i0, t0, rth, tmax, k)
		if ok {
			break
		}
	}

	if !ok || rth <= 0 {
		return "No solution"
	}
	return fmt.Sprintf("%.0f", rth)
}

func DeltaTj(p, rth float64) float64 {
	return p * rth
}

func Ileak(i0, t, k float64) float64 {
	return i0 * math.Exp((t-25)/k)
}

func Ileak2(i0, t float64) float64 {
	return i0 * electronics.AF_Arrhenius(0.74, 25, t)
}

func Runaway(vr, i0, t0, rth, tmax, k float64) bool {

	tj := t0

	tjp := tj
	pp := i0 * vr

	for tj < tmax {

		i := Ileak(i0, tj, k)
		p := i * vr
		tj = DeltaTj(p, rth) + t0

		dtdp := (tj - tjp) / (p - pp)
		// fmt.Printf("dT/dP = %f\n", dtdp)
		if dtdp < rth/2 {
			break
		}

		if math.IsNaN(dtdp) {
			return true
		}

		tjp = tj
		pp = p
	}
	if tj < tmax {
		return false
	}
	return true
}

func (E Etools) LineChart(i0, k float64) string {

	series := m.NewSeries()

	for t := 25.0; t < 150; t++ {
		y := Ileak(i0, t, k) * 1000
		// y := Ileak2(i0, t) * 1000
		series.Add(m.MakeValue(t, y))
	}

	diagram := m.New(600, 400,
		m.WithAutorange(m.XAxis, series),
		m.WithAutorange(m.YAxis, series),
		// m.WithRange(m.YAxis, 0, 1e-3),
		m.WithProjection(m.YAxis, m.Log),
		// m.WithInset(70),
		m.WithPadding(2),
		// m.WithColorScheme(90),
	)

	diagram.Smooth(series, m.UsingAxes(m.XAxis, m.YAxis), m.UsingStrokeWidth(3.14))
	diagram.Axis(series, m.XAxis, diagram.ValueTicker('f', 1, 10), true, "Temp [ºC]")
	diagram.Axis(series, m.YAxis, diagram.ValueTicker('f', 2, 10), true, "Current [mA]")

	diagram.Frame()
	diagram.Title("Reverse current")

	return diagram.String()
}
