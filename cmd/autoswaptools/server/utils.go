package server

import "math"

func toDcrnCoin(a int64) float64 {
	return float64(a) / math.Pow10(8)
}
