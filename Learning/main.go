package main

import (
	"fmt"
)

func main() {
	var x float64
	x = 50
	fmt.Println(x)
	y := calculator(x, x)
	fmt.Println(y)
}
