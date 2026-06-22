package metartafparser

import (
	"regexp"
	"strings"
)

// ============================================================
// TornadicActivityBegEndCommand: TORNADO/FUNNEL CLOUD/WATERSPOUT B[hh]mmE[hh]mm [dist] [dir]
// ============================================================

var tornBegEndRe = regexp.MustCompile(`^(TORNADO|FUNNEL CLOUD|WATERSPOUT) (B(\d{2})?(\d{2}))(E(\d{2})?(\d{2}))( (\d+)? ([A-Z]{1,2})?)?`)

type tornadicActivityBegEndRemarkCommand struct {
	locale Locale
}

func newTornadicActivityBegEndRemarkCommand(locale Locale) *tornadicActivityBegEndRemarkCommand {
	return &tornadicActivityBegEndRemarkCommand{locale: locale}
}

func (c *tornadicActivityBegEndRemarkCommand) CanParse(code string) bool {
	return tornBegEndRe.MatchString(code)
}

func (c *tornadicActivityBegEndRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := tornBegEndRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	tornType := strings.ReplaceAll(m[1], " ", "")
	tornLoc := locStr(c.locale, "Remark."+tornType)
	dirLoc := locStr(c.locale, "Converter."+m[8])
	desc := remarkDesc(c.locale, "Remark.Tornadic.Activity.BegEnd", tornLoc, m[3], m[4], m[6], m[7], m[9], dirLoc)
	rm := Remark{
		Type:        RemarkTypeTornadicActivityBegEnd,
		Description: desc,
		Raw:         m[0],
		StartMinute: intPtr(atoi(m[4])),
		EndMinute:   intPtr(atoi(m[7])),
	}
	if m[3] != "" {
		rm.StartHour = intPtr(atoi(m[3]))
	}
	if m[6] != "" {
		rm.EndHour = intPtr(atoi(m[6]))
	}
	if m[9] != "" {
		rm.Value = float64Ptr(float64(atoi(m[9])))
	}
	if m[10] != "" {
		dir := m[10]
		rm.Direction = &dir
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, tornBegEndRe), remarks, nil
}

// ============================================================
// TornadicActivityBegCommand: TORNADO/FUNNEL CLOUD/WATERSPOUT B[hh]mm [dist] [dir]
// ============================================================

var tornBegRe = regexp.MustCompile(`^(TORNADO|FUNNEL CLOUD|WATERSPOUT) (B(\d{2})?(\d{2}))( (\d+)? ([A-Z]{1,2})?)?`)

type tornadicActivityBegRemarkCommand struct {
	locale Locale
}

func newTornadicActivityBegRemarkCommand(locale Locale) *tornadicActivityBegRemarkCommand {
	return &tornadicActivityBegRemarkCommand{locale: locale}
}

func (c *tornadicActivityBegRemarkCommand) CanParse(code string) bool {
	if tornBegEndRe.MatchString(code) {
		return false
	}
	return tornBegRe.MatchString(code)
}

func (c *tornadicActivityBegRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := tornBegRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	tornType := strings.ReplaceAll(m[1], " ", "")
	tornLoc := locStr(c.locale, "Remark."+tornType)
	dirLoc := locStr(c.locale, "Converter."+m[7])
	desc := remarkDesc(c.locale, "Remark.Tornadic.Activity.Beginning", tornLoc, m[3], m[4], m[6], dirLoc)
	rm := Remark{
		Type:        RemarkTypeTornadicActivityBeg,
		Description: desc,
		Raw:         m[0],
		StartMinute: intPtr(atoi(m[4])),
	}
	if m[3] != "" {
		rm.StartHour = intPtr(atoi(m[3]))
	}
	if m[6] != "" {
		rm.Value = float64Ptr(float64(atoi(m[6])))
	}
	if m[7] != "" {
		dir := m[7]
		rm.Direction = &dir
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, tornBegRe), remarks, nil
}

// ============================================================
// TornadicActivityEndCommand: TORNADO/FUNNEL CLOUD/WATERSPOUT E[hh]mm [dist] [dir]
// ============================================================

var tornEndRe = regexp.MustCompile(`^(TORNADO|FUNNEL CLOUD|WATERSPOUT) (E(\d{2})?(\d{2}))( (\d+)? ([A-Z]{1,2})?)?`)

type tornadicActivityEndRemarkCommand struct {
	locale Locale
}

func newTornadicActivityEndRemarkCommand(locale Locale) *tornadicActivityEndRemarkCommand {
	return &tornadicActivityEndRemarkCommand{locale: locale}
}

func (c *tornadicActivityEndRemarkCommand) CanParse(code string) bool {
	if tornBegEndRe.MatchString(code) {
		return false
	}
	return tornEndRe.MatchString(code)
}

func (c *tornadicActivityEndRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := tornEndRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	tornType := strings.ReplaceAll(m[1], " ", "")
	tornLoc := locStr(c.locale, "Remark."+tornType)
	dirLoc := locStr(c.locale, "Converter."+m[7])
	desc := remarkDesc(c.locale, "Remark.Tornadic.Activity.Ending", tornLoc, m[3], m[4], m[6], dirLoc)
	rm := Remark{
		Type:        RemarkTypeTornadicActivityEnd,
		Description: desc,
		Raw:         m[0],
		EndMinute:   intPtr(atoi(m[4])),
	}
	if m[3] != "" {
		rm.EndHour = intPtr(atoi(m[3]))
	}
	if m[6] != "" {
		rm.Value = float64Ptr(float64(atoi(m[6])))
	}
	if m[7] != "" {
		dir := m[7]
		rm.Direction = &dir
	}
	remarks = append(remarks, rm)
	return trimAfterRegex(code, tornEndRe), remarks, nil
}

// ============================================================
// ThunderStormLocationMovingCommand: TS <loc> MOV <dir>
// ============================================================

var tsMovRe = regexp.MustCompile(`^TS ([A-Z]{2}) MOV ([A-Z]{2})`)

type thunderStormLocationMovingRemarkCommand struct {
	locale Locale
}

func newThunderStormLocationMovingRemarkCommand(locale Locale) *thunderStormLocationMovingRemarkCommand {
	return &thunderStormLocationMovingRemarkCommand{locale: locale}
}

func (c *thunderStormLocationMovingRemarkCommand) CanParse(code string) bool {
	return tsMovRe.MatchString(code)
}

func (c *thunderStormLocationMovingRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := tsMovRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	locLoc := locStr(c.locale, "Converter."+m[1])
	movLoc := locStr(c.locale, "Converter."+m[2])
	desc := remarkDesc(c.locale, "Remark.Thunderstorm.Location.Moving", locLoc, movLoc)
	loc := m[1]
	mov := m[2]
	remarks = append(remarks, Remark{
		Type:        RemarkTypeThunderStormLocationMoving,
		Description: desc,
		Raw:         m[0],
		Location:    &loc,
		Moving:      &mov,
	})
	return trimAfterRegex(code, tsMovRe), remarks, nil
}

// ============================================================
// ThunderStormLocationCommand: TS <loc>
// ============================================================

var tsLocRe = regexp.MustCompile(`^TS ([A-Z]{2})`)

type thunderStormLocationRemarkCommand struct {
	locale Locale
}

func newThunderStormLocationRemarkCommand(locale Locale) *thunderStormLocationRemarkCommand {
	return &thunderStormLocationRemarkCommand{locale: locale}
}

func (c *thunderStormLocationRemarkCommand) CanParse(code string) bool {
	if tsMovRe.MatchString(code) {
		return false
	}
	return tsLocRe.MatchString(code)
}

func (c *thunderStormLocationRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := tsLocRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	locLoc := locStr(c.locale, "Converter."+m[1])
	desc := remarkDesc(c.locale, "Remark.Thunderstorm.Location.0", locLoc)
	loc := m[1]
	remarks = append(remarks, Remark{
		Type:        RemarkTypeThunderStormLocation,
		Description: desc,
		Raw:         m[0],
		Location:    &loc,
	})
	return trimAfterRegex(code, tsLocRe), remarks, nil
}

// ============================================================
// SmallHailSizeCommand: GR LESS THAN <size>
// ============================================================

var smallHailRe = regexp.MustCompile(`^GR LESS THAN ((\d )?(\d\/\d)?)`)

type smallHailSizeRemarkCommand struct {
	locale Locale
}

func newSmallHailSizeRemarkCommand(locale Locale) *smallHailSizeRemarkCommand {
	return &smallHailSizeRemarkCommand{locale: locale}
}

func (c *smallHailSizeRemarkCommand) CanParse(code string) bool { return smallHailRe.MatchString(code) }

func (c *smallHailSizeRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := smallHailRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	size, err := convertFractionalAmount(m[1])
	if err != nil {
		return code, remarks, err
	}
	desc := remarkDesc(c.locale, "Remark.Hail.LesserThan", m[1])
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSmallHailSize,
		Description: desc,
		Raw:         m[0],
		Value:       &size,
	})
	return trimAfterRegex(code, smallHailRe), remarks, nil
}

// ============================================================
// HailSizeCommand: GR <size>
// ============================================================

var hailSizeRe = regexp.MustCompile(`^GR ((\d\/\d)|((\d) ?(\d\/\d)?))`)

type hailSizeRemarkCommand struct {
	locale Locale
}

func newHailSizeRemarkCommand(locale Locale) *hailSizeRemarkCommand {
	return &hailSizeRemarkCommand{locale: locale}
}

func (c *hailSizeRemarkCommand) CanParse(code string) bool {
	if smallHailRe.MatchString(code) {
		return false
	}
	return hailSizeRe.MatchString(code)
}

func (c *hailSizeRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := hailSizeRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	size, err := convertFractionalAmount(m[1])
	if err != nil {
		return code, remarks, err
	}
	desc := remarkDesc(c.locale, "Remark.Hail.0", m[1])
	remarks = append(remarks, Remark{
		Type:        RemarkTypeHailSize,
		Description: desc,
		Raw:         m[0],
		Value:       &size,
	})
	return trimAfterRegex(code, hailSizeRe), remarks, nil
}

// ============================================================
// SnowPelletsCommand: GS <LGT|MOD|HVY>
// ============================================================

var snowPelletsRe = regexp.MustCompile(`^GS (LGT|MOD|HVY)`)

type snowPelletsRemarkCommand struct {
	locale Locale
}

func newSnowPelletsRemarkCommand(locale Locale) *snowPelletsRemarkCommand {
	return &snowPelletsRemarkCommand{locale: locale}
}

func (c *snowPelletsRemarkCommand) CanParse(code string) bool { return snowPelletsRe.MatchString(code) }

func (c *snowPelletsRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := snowPelletsRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	amountLoc := locStr(c.locale, "Remark."+m[1])
	desc := remarkDesc(c.locale, "Remark.Snow.Pellets", amountLoc)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSnowPellets,
		Description: desc,
		Raw:         m[0],
	})
	return trimAfterRegex(code, snowPelletsRe), remarks, nil
}
