package main

import (
	"fmt"
	"strconv"
)

type Symbol struct {
	name string
}

func main() {
	for _, s := range []string{"3.1415", "-10", "0", "1e-5", "1.3e+5", ".223", "+.45"} {
		f, err := strconv.ParseFloat(s, 64)
		fmt.Printf("%v => %v (%v)\n", s, f, err)
	}
}
