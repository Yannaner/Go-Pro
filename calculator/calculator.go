package calculator

import (
	"errors"
	"fmt"
	"math"
)

// Add returns the sum of two numbers
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference of two numbers
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of two numbers
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns the quotient of two numbers
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// Power returns a raised to the power of b
func Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// SquareRoot returns the square root of a number
func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("cannot calculate square root of negative number")
	}
	return math.Sqrt(a), nil
}

// PrintResult formats and prints calculation results
func PrintResult(operation string, result float64) {
	fmt.Printf("%s: %.2f\n", operation, result)
}
