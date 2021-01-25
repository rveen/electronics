package electronics

import (
	"strings"
)

// Rth returns the default thermal resistance of the component (to ambient) based
// on the package.
// Source:
// - http://www.yageo.com/exep/pages/download/literatures/PYu-R_Mount_7.pdf
func Rth(pkg string) float64 {

	pkg = strings.ToLower(pkg)

	switch pkg {
	case "0201":
		return 800
	case "0603":
		return 400
	case "0805":
		return 250
	case "1206":
		return 200
	case "1210":
		return 125
	case "1218":
		return 100
	case "2010":
		return 80
	case "2512":
		return 50
	}

	if strings.HasPrefix(pkg, "sot23") {
		return 357
	}

	return 0
}
