package metartafparser

import (
	"math"
	"strconv"
	"strings"
)

func degreesToCardinal(input string) Direction {
	d, err := strconv.ParseFloat(input, 64)
	if err != nil || math.IsInf(d, 0) || math.IsNaN(d) {
		return DirectionVRB
	}
	dirs := []Direction{
		DirectionN, DirectionNNE, DirectionNE, DirectionENE,
		DirectionE, DirectionESE, DirectionSE, DirectionSSE,
		DirectionS, DirectionSSW, DirectionSW, DirectionWSW,
		DirectionW, DirectionWNW, DirectionNW, DirectionNNW,
	}
	ix := int(math.Floor((d + 11.25) / 22.5))
	return dirs[ix%16]
}

func convertVisibility(input string) (Distance, error) {
	v, err := strconv.Atoi(input)
	if err != nil {
		return Distance{}, err
	}
	if input == "9999" {
		return Distance{
			Indicator: valueIndicatorPtr(ValueIndicatorGreaterThan),
			Value:     float64(v),
			Unit:      DistanceUnitMeters,
		}, nil
	}
	return Distance{
		Value: float64(v),
		Unit:  DistanceUnitMeters,
	}, nil
}

func convertNauticalMilesVisibility(input string) (Distance, error) {
	var indicator *ValueIndicator
	idx := 0
	if strings.HasPrefix(input, "P") {
		indicator = valueIndicatorPtr(ValueIndicatorGreaterThan)
		idx = 1
	} else if strings.HasPrefix(input, "M") {
		indicator = valueIndicatorPtr(ValueIndicatorLessThan)
		idx = 1
	}
	valueStr := input[idx : len(input)-2]
	val, err := convertFractionalAmount(valueStr)
	if err != nil {
		return Distance{}, err
	}
	return Distance{
		Indicator: indicator,
		Value:     val,
		Unit:      DistanceUnitStatuteMiles,
	}, nil
}

func convertFractionalAmount(input string) (float64, error) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) < 2 {
		return parseFraction(parts[0])
	}
	whole, err := strconv.ParseFloat(parts[0], 64)
	if err != nil || math.IsInf(whole, 0) || math.IsNaN(whole) {
		return 0, strconv.ErrRange
	}
	frac, err := parseFraction(parts[1])
	if err != nil {
		return 0, err
	}
	return whole + frac, nil
}

func parseFraction(input string) (float64, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) < 2 {
		v, err := strconv.ParseFloat(parts[0], 64)
		if err != nil || math.IsInf(v, 0) || math.IsNaN(v) {
			return 0, strconv.ErrRange
		}
		return v, nil
	}
	top, err := strconv.ParseFloat(parts[0], 64)
	if err != nil || math.IsInf(top, 0) || math.IsNaN(top) {
		return 0, strconv.ErrRange
	}
	bottom, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || math.IsInf(bottom, 0) || math.IsNaN(bottom) {
		return 0, strconv.ErrRange
	}
	if bottom == 0 {
		return 0, nil
	}
	return math.Round(top/bottom*100) / 100, nil
}

func convertTemperature(input string) (float64, error) {
	if strings.HasPrefix(input, "M") {
		parts := splitN(input, "M")
		v, err := strconv.ParseFloat(parts[1], 64)
		if err != nil || math.IsInf(v, 0) || math.IsNaN(v) {
			return 0, strconv.ErrRange
		}
		return -v, nil
	}
	v, err := strconv.ParseFloat(input, 64)
	if err != nil || math.IsInf(v, 0) || math.IsNaN(v) {
		return 0, strconv.ErrRange
	}
	return v, nil
}

func convertTemperatureRemarks(sign, temperature string) float64 {
	temp, _ := strconv.ParseFloat(temperature, 64)
	if math.IsInf(temp, 0) || math.IsNaN(temp) {
		return 0
	}
	temp /= 10
	if sign == "0" {
		return temp
	}
	return -temp
}

func convertPrecipitationAmount(amount string) float64 {
	v, _ := strconv.ParseFloat(amount, 64)
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return 0
	}
	return v / 100
}
