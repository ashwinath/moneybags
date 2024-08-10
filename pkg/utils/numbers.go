package utils

import "math"

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}
