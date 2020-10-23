// Some utility functions for slices
package sliceutil

import "strconv"

// Create a slice of length "size" of float64 filled with "value"
func FillF64(value float64, size int) []float64 {
	f64 := make([]float64, size)
	for i := range f64 {
		f64[i] = value
	}
	return f64
}

// I64ToF64 convert a slice of string into a slice of float64
func StrToF64(data []string) ([]float64, error) {
	f64 := make([]float64, len(data))
	var s string
	var i int
	for i, s = range data {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		f64[i] = f
	}
	return f64, nil
}

// I64ToF64 convert a slice of int64 into a slice of float64
func I64ToF64(data []int64) []float64 {
	f64 := make([]float64, len(data))
	var f int64
	var i int
	for i, f = range data {
		f64[i] = float64(f)
	}
	return f64
}

// MapF64 transform a slice of float64 through a function of float64 -> float64
func MapF64(data []float64, f func(float64) float64) []float64 {
	f64 := make([]float64, len(data))
	var v float64
	var i int
	for i, v = range data {
		f64[i] = f(v)
	}
	return f64
}

// FilterF64 filter a slice of float64 through a function of float64 -> bool
func FilterF64(data []float64, f func(float64) bool) []float64 {
	f64 := make([]float64, 0)
	var v float64
	var i int = 0
	for i, v = range data {
		if f(v) {
			f64[i] = v
			i++
		}
	}
	return f64
}

// MinMax compute the min and max values of data
func MinMax(data []float64) (float64, float64) {
	min := data[0]
	max := min
	var d float64
	for _, d = range data {
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}
	return min, max
}
