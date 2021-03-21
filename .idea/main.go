package main

import (
	"fmt"
)

type Symbol struct {
	name string
}

func main() {
	s := "Zażółć gęślą jaźń"

	for i, ch := range s {
		fmt.Printf("%d %s\n", i, string(ch))
	}
}
