package metartafparser

import (
	"fmt"
	"regexp"
	"strconv"
)

// ============================================================
// SeaLevelPressureCommand: SLPppp
// ============================================================

var slpRe = regexp.MustCompile(`^SLP(\d{2})(\d)`)

type seaLevelPressureRemarkCommand struct {
	locale Locale
}

func newSeaLevelPressureRemarkCommand(locale Locale) *seaLevelPressureRemarkCommand {
	return &seaLevelPressureRemarkCommand{locale: locale}
}

func (c *seaLevelPressureRemarkCommand) CanParse(code string) bool { return slpRe.MatchString(code) }

func (c *seaLevelPressureRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := slpRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	ppp, _ := strconv.Atoi(m[1] + m[2])
	var press float64
	if ppp >= 500 {
		press = float64(9000+ppp) / 10
	} else {
		press = float64(10000+ppp) / 10
	}
	desc := remarkDesc(c.locale, "Remark.Sea.Level.Pressure", press)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSeaLevelPressure,
		Description: desc,
		Raw:         m[0],
		Value:       &press,
	})
	return trimAfterRegex(code, slpRe), remarks, nil
}

// ============================================================
// HourlyMaximumMinimumTemperatureCommand: 4sTTTsTTT
// ============================================================

var hourlyMaxMinTempRe = regexp.MustCompile(`^4([01])(\d{3})([01])(\d{3})`)

type hourlyMaximumMinimumTemperatureRemarkCommand struct {
	locale Locale
}

func newHourlyMaximumMinimumTemperatureRemarkCommand(locale Locale) *hourlyMaximumMinimumTemperatureRemarkCommand {
	return &hourlyMaximumMinimumTemperatureRemarkCommand{locale: locale}
}

func (c *hourlyMaximumMinimumTemperatureRemarkCommand) CanParse(code string) bool {
	return hourlyMaxMinTempRe.MatchString(code)
}

func (c *hourlyMaximumMinimumTemperatureRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := hourlyMaxMinTempRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	maxTemp := convertTemperatureRemarks(m[1], m[2])
	minTemp := convertTemperatureRemarks(m[3], m[4])
	maxStr := fmt.Sprintf("%.1f", maxTemp)
	minStr := fmt.Sprintf("%.1f", minTemp)
	desc := remarkDesc(c.locale, "Remark.Hourly.Maximum.Minimum.Temperature", maxStr, minStr)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeHourlyMaximumMinimumTemperature,
		Description: desc,
		Raw:         m[0],
		Max:         float64Ptr(maxTemp),
		Min:         float64Ptr(minTemp),
	})
	return trimAfterRegex(code, hourlyMaxMinTempRe), remarks, nil
}

// ============================================================
// HourlyMaximumTemperatureCommand: 1sTTT
// ============================================================

var hourlyMaxTempRe = regexp.MustCompile(`^1([01])(\d{3})`)

type hourlyMaximumTemperatureRemarkCommand struct {
	locale Locale
}

func newHourlyMaximumTemperatureRemarkCommand(locale Locale) *hourlyMaximumTemperatureRemarkCommand {
	return &hourlyMaximumTemperatureRemarkCommand{locale: locale}
}

func (c *hourlyMaximumTemperatureRemarkCommand) CanParse(code string) bool {
	return hourlyMaxTempRe.MatchString(code)
}

func (c *hourlyMaximumTemperatureRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := hourlyMaxTempRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	maxTemp := convertTemperatureRemarks(m[1], m[2])
	maxStr := fmt.Sprintf("%.1f", maxTemp)
	desc := remarkDesc(c.locale, "Remark.Hourly.Maximum.Temperature", maxStr)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeHourlyMaximumTemperature,
		Description: desc,
		Raw:         m[0],
		Max:         float64Ptr(maxTemp),
	})
	return trimAfterRegex(code, hourlyMaxTempRe), remarks, nil
}

// ============================================================
// HourlyMinimumTemperatureCommand: 2sTTT
// ============================================================

var hourlyMinTempRe = regexp.MustCompile(`^2([01])(\d{3})`)

type hourlyMinimumTemperatureRemarkCommand struct {
	locale Locale
}

func newHourlyMinimumTemperatureRemarkCommand(locale Locale) *hourlyMinimumTemperatureRemarkCommand {
	return &hourlyMinimumTemperatureRemarkCommand{locale: locale}
}

func (c *hourlyMinimumTemperatureRemarkCommand) CanParse(code string) bool {
	return hourlyMinTempRe.MatchString(code)
}

func (c *hourlyMinimumTemperatureRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := hourlyMinTempRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	minTemp := convertTemperatureRemarks(m[1], m[2])
	minStr := fmt.Sprintf("%.1f", minTemp)
	desc := remarkDesc(c.locale, "Remark.Hourly.Minimum.Temperature", minStr)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeHourlyMinimumTemperature,
		Description: desc,
		Raw:         m[0],
		Min:         float64Ptr(minTemp),
	})
	return trimAfterRegex(code, hourlyMinTempRe), remarks, nil
}

// ============================================================
// HourlyTemperatureDewPointCommand: T[s]TTT[s]TTT
// ============================================================

var hourlyTempDewRe = regexp.MustCompile(`^T([01])(\d{3})(([01])(\d{3}))?`)

type hourlyTemperatureDewPointRemarkCommand struct {
	locale Locale
}

func newHourlyTemperatureDewPointRemarkCommand(locale Locale) *hourlyTemperatureDewPointRemarkCommand {
	return &hourlyTemperatureDewPointRemarkCommand{locale: locale}
}

func (c *hourlyTemperatureDewPointRemarkCommand) CanParse(code string) bool {
	return hourlyTempDewRe.MatchString(code)
}

func (c *hourlyTemperatureDewPointRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := hourlyTempDewRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	temp := convertTemperatureRemarks(m[1], m[2])
	var desc *string
	if m[3] != "" {
		dew := convertTemperatureRemarks(m[4], m[5])
		tempStr := fmt.Sprintf("%.1f", temp)
		dewStr := fmt.Sprintf("%.1f", dew)
		desc = remarkDesc(c.locale, "Remark.Hourly.Temperature.Dew.Point", tempStr, dewStr)
	} else {
		tempStr := fmt.Sprintf("%.1f", temp)
		desc = remarkDesc(c.locale, "Remark.Hourly.Temperature.0", tempStr)
	}
	rm := Remark{
		Type:        RemarkTypeHourlyTemperatureDewPoint,
		Description: desc,
		Raw:         m[0],
		Temperature: float64Ptr(temp),
	}
	if m[3] != "" {
		dew := convertTemperatureRemarks(m[4], m[5])
		rm.DewPoint = float64Ptr(dew)
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, hourlyTempDewRe), remarks, nil
}

// ============================================================
// HourlyPressureCommand: 5cttt
// ============================================================

var hourlyPressureRe = regexp.MustCompile(`^5(\d)(\d{3})`)

type hourlyPressureRemarkCommand struct {
	locale Locale
}

func newHourlyPressureRemarkCommand(locale Locale) *hourlyPressureRemarkCommand {
	return &hourlyPressureRemarkCommand{locale: locale}
}

func (c *hourlyPressureRemarkCommand) CanParse(code string) bool {
	return hourlyPressureRe.MatchString(code)
}

func (c *hourlyPressureRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := hourlyPressureRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	codeVal := atoi(m[1])
	change := float64(atoi(m[2])) / 10
	barometerKey := fmt.Sprintf("Remark.Barometer.%d", codeVal)
	barometerDesc := locStr(c.locale, barometerKey)
	tendencyDesc := locStr(c.locale, "Remark.Pressure.Tendency")
	var fullDesc string
	if barometerDesc != "" && tendencyDesc != "" {
		fullDesc = barometerDesc + " " + tendencyDesc
	}
	desc := &fullDesc
	remarks = append(remarks, Remark{
		Type:        RemarkTypeHourlyPressure,
		Description: desc,
		Raw:         m[0],
		Value:       &change,
	})
	return trimAfterRegex(code, hourlyPressureRe), remarks, nil
}
