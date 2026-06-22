package metartafparser

import (
	"regexp"
)

// ============================================================
// WindPeakCommand: PK WND dddss[/hhmm]
// ============================================================

var windPeakRe = regexp.MustCompile(`^PK WND (\d{3})(\d{2,3})\/(\d{2})?(\d{2})`)

type windPeakRemarkCommand struct {
	locale Locale
}

func newWindPeakRemarkCommand(locale Locale) *windPeakRemarkCommand {
	return &windPeakRemarkCommand{locale: locale}
}

func (c *windPeakRemarkCommand) CanParse(code string) bool { return windPeakRe.MatchString(code) }

func (c *windPeakRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := windPeakRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	degrees := atoi(m[1])
	speed := atoi(m[2])
	desc := remarkDesc(c.locale, "Remark.PeakWind", degrees, speed, m[3], m[4])
	rm := Remark{
		Type:        RemarkTypeWindPeak,
		Description: desc,
		Raw:         m[0],
		Degrees:     &degrees,
		Speed:       &speed,
		StartMinute: intPtr(atoi(m[4])),
	}
	if m[3] != "" {
		rm.StartHour = intPtr(atoi(m[3]))
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, windPeakRe), remarks, nil
}

// ============================================================
// WindShiftFropaCommand: WSHFT [hh]mm FROPA
// ============================================================

var wshftFropaRe = regexp.MustCompile(`^WSHFT (\d{2})?(\d{2}) FROPA`)

type windShiftFropaRemarkCommand struct {
	locale Locale
}

func newWindShiftFropaRemarkCommand(locale Locale) *windShiftFropaRemarkCommand {
	return &windShiftFropaRemarkCommand{locale: locale}
}

func (c *windShiftFropaRemarkCommand) CanParse(code string) bool {
	return wshftFropaRe.MatchString(code)
}

func (c *windShiftFropaRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := wshftFropaRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	desc := remarkDesc(c.locale, "Remark.WindShift.FROPA", m[1], m[2])
	rm := Remark{
		Type:        RemarkTypeWindShiftFropa,
		Description: desc,
		Raw:         m[0],
		StartMinute: intPtr(atoi(m[2])),
	}
	if m[1] != "" {
		rm.StartHour = intPtr(atoi(m[1]))
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, wshftFropaRe), remarks, nil
}

// ============================================================
// WindShiftCommand: WSHFT [hh]mm
// ============================================================

var wshftRe = regexp.MustCompile(`^WSHFT (\d{2})?(\d{2})`)

type windShiftRemarkCommand struct {
	locale Locale
}

func newWindShiftRemarkCommand(locale Locale) *windShiftRemarkCommand {
	return &windShiftRemarkCommand{locale: locale}
}

func (c *windShiftRemarkCommand) CanParse(code string) bool {
	if wshftFropaRe.MatchString(code) {
		return false
	}
	return wshftRe.MatchString(code)
}

func (c *windShiftRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := wshftRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	desc := remarkDesc(c.locale, "Remark.WindShift.0", m[1], m[2])
	rm := Remark{
		Type:        RemarkTypeWindShift,
		Description: desc,
		Raw:         m[0],
		StartMinute: intPtr(atoi(m[2])),
	}
	if m[1] != "" {
		rm.StartHour = intPtr(atoi(m[1]))
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, wshftRe), remarks, nil
}

// ============================================================
// VirgaDirectionCommand: VIRGA <dir>
// ============================================================

var virgaDirRe = regexp.MustCompile(`^VIRGA ([A-Z]{2})`)

type virgaDirectionRemarkCommand struct {
	locale Locale
}

func newVirgaDirectionRemarkCommand(locale Locale) *virgaDirectionRemarkCommand {
	return &virgaDirectionRemarkCommand{locale: locale}
}

func (c *virgaDirectionRemarkCommand) CanParse(code string) bool { return virgaDirRe.MatchString(code) }

func (c *virgaDirectionRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := virgaDirRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	dirLoc := locStr(c.locale, "Converter."+m[1])
	desc := remarkDesc(c.locale, "Remark.Virga.Direction", dirLoc)
	dir := m[1]
	remarks = append(remarks, Remark{
		Type:        RemarkTypeVirgaDirection,
		Description: desc,
		Raw:         m[0],
		Direction:   &dir,
	})
	return trimAfterRegex(code, virgaDirRe), remarks, nil
}
