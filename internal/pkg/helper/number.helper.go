package helper

import (
	"math"
	"strconv"
)

func Abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func RoundWithPrecision(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

func AbsAndRoundWithPrecision(value float64, precision int) float64 {
	return RoundWithPrecision(Abs(value), precision)
}

func FmtFloatWithPrecision(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}

func FmtRound(value float64, precision int) string {
	return FmtFloatWithPrecision(RoundWithPrecision(value, precision), precision)
}

func FmtRoundAbs(value float64, precision int) string {
	return FmtFloatWithPrecision(RoundWithPrecision(Abs(value), precision), precision)
}
