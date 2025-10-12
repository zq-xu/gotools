package utils

import (
	"math"
)

// RoundFloat64 float64 Keep the decimal places
// value float64 the incoming parameter
// prec int Keep the number of digits after the decimal point
func RoundFloat64(number float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(number*shift) / shift
}

// TruncateFloat64 float64 Keep the decimal places and drop the extra digits
// value float64 the incoming parameter
// prec int Keep the number of digits after the decimal point
func TruncateFloat64(number float64, prec int) float64 {
	scale := 1.0
	for i := 0; i < prec; i++ {
		scale *= 10.0
	}
	return float64(int(number*scale)) / scale
}
