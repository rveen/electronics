package main

import (
	"fmt"

	"github.com/rveen/electronics"
)

func main() {

	s := "10 k}"

	f := electronics.Value(s)

	fmt.Printf("%g\n", f)

	s = "10k}"

	f = electronics.Value(s)

	fmt.Printf("%g\n", f)

}
