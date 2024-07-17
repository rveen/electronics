package electronics

import (
	"strings"
)

// Rth returns the default thermal resistance of the component (to ambient) based
// on the package.
// Source:
// - https://www.vishay.com/docs/53048/pprachp.pdf
func Rth(pkg string) float64 {

	pkg = strings.ToLower(pkg)

	switch pkg {
	case "0201":
		return 1000   // To be confirmed
	case "0402":
		return 870		
	case "0603":
		return 550
	case "0612":
		return 220	// TBC
	case "0805":
		return 440
	case "1206":
		return 220
	case "1210":
		return 140
	case "1218":
		return 125    // TBC
	case "2010":
		return 110
	case "2512":
		return 50  	// TBD
	}

	if strings.HasPrefix(pkg, "sot23") {
		return 357
	}

	return 0
}
