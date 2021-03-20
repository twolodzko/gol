package main

import (
	"fmt"
	"reflect"
)

type Symbol struct {
	name string
}

func main() {
	x := Symbol{"foo"}
	fmt.Printf("%v\n", reflect.TypeOf(x))
}
