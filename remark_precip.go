package metartafparser

import (
	"regexp"
	"strings"
)

// ============================================================
// PrecipitationBegEndCommand: [desc]phenB[hh]mmE[hh]mm
// ============================================================

var precipBegEndRe = regexp.MustCompile(`^(([A-Z]{2})?([A-Z]{2})B(\d{2})?(\d{2})E(\d{2})?(\d{2}))`)

type precipitationBegEndRemarkCommand struct {
	locale Locale
}

func newPrecipitationBegEndRemarkCommand(locale Locale) *precipitationBegEndRemarkCommand {
	return &precipitationBegEndRemarkCommand{locale: locale}
}

func (c *precipitationBegEndRemarkCommand) CanParse(code string) bool {
	m := precipBegEndRe.FindStringSubmatch(code)
	if m == nil {
		return false
	}
	return !strings.HasPrefix(m[0], "Q")
}

func (c *precipitationBegEndRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := precipBegEndRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	descLoc := ""
	if m[2] != "" {
		descLoc = locStr(c.locale, "Descriptive."+m[2])
	}
	phenLoc := locStr(c.locale, "Phenomenon."+m[3])
	desc := remarkDesc(c.locale, "Remark.Precipitation.Beg.End", descLoc, phenLoc, m[4], m[5], m[6], m[7])
	rm := Remark{
		Type:        RemarkTypePrecipitationBegEnd,
		Description: desc,
		Raw:         m[0],
		StartMinute: intPtr(atoi(m[5])),
		EndMinute:   intPtr(atoi(m[7])),
	}
	if m[2] != "" {
		d := Descriptive(m[2])
		rm.Descriptive = &d
	}
	if m[3] != "" {
		p := Phenomenon(m[3])
		rm.Phenomenon = &p
	}
	if m[4] != "" {
		rm.StartHour = intPtr(atoi(m[4]))
	}
	if m[6] != "" {
		rm.EndHour = intPtr(atoi(m[6]))
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, precipBegEndRe), remarks, nil
}

// ============================================================
// PrecipitationBegCommand: [desc]phenB[hh]mm
// ============================================================

var precipBegRe = regexp.MustCompile(`^(([A-Z]{2})?([A-Z]{2})B(\d{2})?(\d{2}))`)

type precipitationBegRemarkCommand struct {
	locale Locale
}

func newPrecipitationBegRemarkCommand(locale Locale) *precipitationBegRemarkCommand {
	return &precipitationBegRemarkCommand{locale: locale}
}

func (c *precipitationBegRemarkCommand) CanParse(code string) bool {
	if precipBegEndRe.MatchString(code) {
		return false
	}
	m := precipBegRe.FindStringSubmatch(code)
	if m == nil {
		return false
	}
	return !strings.HasPrefix(m[0], "Q")
}

func (c *precipitationBegRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := precipBegRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	descLoc := ""
	if m[2] != "" {
		descLoc = locStr(c.locale, "Descriptive."+m[2])
	}
	phenLoc := locStr(c.locale, "Phenomenon."+m[3])
	desc := remarkDesc(c.locale, "Remark.Precipitation.Beg.0", descLoc, phenLoc, m[4], m[5])
	rm := Remark{
		Type:        RemarkTypePrecipitationBeg,
		Description: desc,
		Raw:         m[0],
		StartMinute: intPtr(atoi(m[5])),
	}
	if m[2] != "" {
		d := Descriptive(m[2])
		rm.Descriptive = &d
	}
	if m[3] != "" {
		p := Phenomenon(m[3])
		rm.Phenomenon = &p
	}
	if m[4] != "" {
		rm.StartHour = intPtr(atoi(m[4]))
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, precipBegRe), remarks, nil
}

// ============================================================
// PrecipitationEndCommand: [desc]phenE[hh]mm
// ============================================================

var precipEndRe = regexp.MustCompile(`^(([A-Z]{2})?([A-Z]{2})E(\d{2})?(\d{2}))`)

type precipitationEndRemarkCommand struct {
	locale Locale
}

func newPrecipitationEndRemarkCommand(locale Locale) *precipitationEndRemarkCommand {
	return &precipitationEndRemarkCommand{locale: locale}
}

func (c *precipitationEndRemarkCommand) CanParse(code string) bool {
	if precipBegEndRe.MatchString(code) {
		return false
	}
	m := precipEndRe.FindStringSubmatch(code)
	if m == nil {
		return false
	}
	return !strings.HasPrefix(m[0], "Q")
}

func (c *precipitationEndRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := precipEndRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	descLoc := ""
	if m[2] != "" {
		descLoc = locStr(c.locale, "Descriptive."+m[2])
	}
	phenLoc := locStr(c.locale, "Phenomenon."+m[3])
	desc := remarkDesc(c.locale, "Remark.Precipitation.End", descLoc, phenLoc, m[4], m[5])
	rm := Remark{
		Type:        RemarkTypePrecipitationEnd,
		Description: desc,
		Raw:         m[0],
		EndMinute:   intPtr(atoi(m[5])),
	}
	if m[2] != "" {
		d := Descriptive(m[2])
		rm.Descriptive = &d
	}
	if m[3] != "" {
		p := Phenomenon(m[3])
		rm.Phenomenon = &p
	}
	if m[4] != "" {
		rm.EndHour = intPtr(atoi(m[4]))
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, precipEndRe), remarks, nil
}

// ============================================================
// HourlyPrecipitationAmountCommand: Ppppp
// ============================================================

var hourlyPrecipRe = regexp.MustCompile(`^P(\d{4})`)

type hourlyPrecipitationAmountRemarkCommand struct {
	locale Locale
}

func newHourlyPrecipitationAmountRemarkCommand(locale Locale) *hourlyPrecipitationAmountRemarkCommand {
	return &hourlyPrecipitationAmountRemarkCommand{locale: locale}
}

func (c *hourlyPrecipitationAmountRemarkCommand) CanParse(code string) bool {
	return hourlyPrecipRe.MatchString(code)
}

func (c *hourlyPrecipitationAmountRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := hourlyPrecipRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	amount := float64(atoi(m[1]))
	amountVal := amount / 100
	desc := remarkDesc(c.locale, "Remark.Precipitation.Amount.Hourly", amountVal)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeHourlyPrecipitationAmount,
		Description: desc,
		Raw:         m[0],
		Amount:      &amountVal,
	})
	return trimAfterRegex(code, hourlyPrecipRe), remarks, nil
}

// ============================================================
// PrecipitationAmount36HourCommand: [36]pppp
// ============================================================

var precipAmount36Re = regexp.MustCompile(`^([36])(\d{4})`)

type precipitationAmount36HourRemarkCommand struct {
	locale Locale
}

func newPrecipitationAmount36HourRemarkCommand(locale Locale) *precipitationAmount36HourRemarkCommand {
	return &precipitationAmount36HourRemarkCommand{locale: locale}
}

func (c *precipitationAmount36HourRemarkCommand) CanParse(code string) bool {
	return precipAmount36Re.MatchString(code)
}

func (c *precipitationAmount36HourRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := precipAmount36Re.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	period := atoi(m[1])
	amount := convertPrecipitationAmount(m[2])
	desc := remarkDesc(c.locale, "Remark.Precipitation.Amount.3.6", period, amount)
	remarks = append(remarks, Remark{
		Type:        RemarkTypePrecipitationAmount36Hour,
		Description: desc,
		Raw:         m[0],
		Amount:      &amount,
	})
	return trimAfterRegex(code, precipAmount36Re), remarks, nil
}

// ============================================================
// PrecipitationAmount24HourCommand: 7pppp
// ============================================================

var precipAmount24Re = regexp.MustCompile(`^7(\d{4})`)

type precipitationAmount24HourRemarkCommand struct {
	locale Locale
}

func newPrecipitationAmount24HourRemarkCommand(locale Locale) *precipitationAmount24HourRemarkCommand {
	return &precipitationAmount24HourRemarkCommand{locale: locale}
}

func (c *precipitationAmount24HourRemarkCommand) CanParse(code string) bool {
	return precipAmount24Re.MatchString(code)
}

func (c *precipitationAmount24HourRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := precipAmount24Re.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	amount := convertPrecipitationAmount(m[1])
	desc := remarkDesc(c.locale, "Remark.Precipitation.Amount.24", amount)
	remarks = append(remarks, Remark{
		Type:        RemarkTypePrecipitationAmount24Hour,
		Description: desc,
		Raw:         m[0],
		Amount:      &amount,
	})
	return trimAfterRegex(code, precipAmount24Re), remarks, nil
}

// ============================================================
// IceAccretionCommand: lhammm
// ============================================================

var iceAccretionRe = regexp.MustCompile(`^l(\d)(\d{3})`)

type iceAccretionRemarkCommand struct {
	locale Locale
}

func newIceAccretionRemarkCommand(locale Locale) *iceAccretionRemarkCommand {
	return &iceAccretionRemarkCommand{locale: locale}
}

func (c *iceAccretionRemarkCommand) CanParse(code string) bool {
	return iceAccretionRe.MatchString(code)
}

func (c *iceAccretionRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := iceAccretionRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	period := atoi(m[1])
	amountRaw := atoi(m[2])
	amount := float64(amountRaw) / 100
	desc := remarkDesc(c.locale, "Remark.Ice.Accretion.Amount", amountRaw, period)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeIceAccretion,
		Description: desc,
		Raw:         m[0],
		Amount:      &amount,
	})
	return trimAfterRegex(code, iceAccretionRe), remarks, nil
}
