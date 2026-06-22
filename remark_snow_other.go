package metartafparser

import (
	"regexp"
)

// ============================================================
// SnowIncreaseCommand: SNINCR inc/total
// ============================================================

var snowIncrRe = regexp.MustCompile(`^SNINCR (\d+)\/(\d+)`)

type snowIncreaseRemarkCommand struct {
	locale Locale
}

func newSnowIncreaseRemarkCommand(locale Locale) *snowIncreaseRemarkCommand {
	return &snowIncreaseRemarkCommand{locale: locale}
}

func (c *snowIncreaseRemarkCommand) CanParse(code string) bool { return snowIncrRe.MatchString(code) }

func (c *snowIncreaseRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := snowIncrRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	inch := atoi(m[1])
	depth := atoi(m[2])
	desc := remarkDesc(c.locale, "Remark.Snow.Increasing.Rapidly", inch, depth)
	remarks = append(remarks, Remark{
		Type:           RemarkTypeSnowIncrease,
		Description:    desc,
		Raw:            m[0],
		InchesLastHour: &inch,
		TotalDepth:     &depth,
	})
	return trimAfterRegex(code, snowIncrRe), remarks, nil
}

// ============================================================
// SnowDepthCommand: 4/ddd
// ============================================================

var snowDepthRe = regexp.MustCompile(`^4\/(\d{3})`)

type snowDepthRemarkCommand struct {
	locale Locale
}

func newSnowDepthRemarkCommand(locale Locale) *snowDepthRemarkCommand {
	return &snowDepthRemarkCommand{locale: locale}
}

func (c *snowDepthRemarkCommand) CanParse(code string) bool { return snowDepthRe.MatchString(code) }

func (c *snowDepthRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := snowDepthRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	depth := atoi(m[1])
	desc := remarkDesc(c.locale, "Remark.Snow.Depth", depth)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSnowDepth,
		Description: desc,
		Raw:         m[0],
		Value:       float64Ptr(float64(depth)),
	})
	return trimAfterRegex(code, snowDepthRe), remarks, nil
}

// ============================================================
// SunshineDurationCommand: 98ddd
// ============================================================

var sunshineRe = regexp.MustCompile(`^98(\d{3})`)

type sunshineDurationRemarkCommand struct {
	locale Locale
}

func newSunshineDurationRemarkCommand(locale Locale) *sunshineDurationRemarkCommand {
	return &sunshineDurationRemarkCommand{locale: locale}
}

func (c *sunshineDurationRemarkCommand) CanParse(code string) bool {
	return sunshineRe.MatchString(code)
}

func (c *sunshineDurationRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := sunshineRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	duration := atoi(m[1])
	desc := remarkDesc(c.locale, "Remark.Sunshine.Duration", duration)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSunshineDuration,
		Description: desc,
		Raw:         m[0],
		Min:         float64Ptr(float64(duration)),
	})
	return trimAfterRegex(code, sunshineRe), remarks, nil
}

// ============================================================
// WaterEquivalentSnowCommand: 933ddd
// ============================================================

var waterEquivSnowRe = regexp.MustCompile(`^933(\d{3})`)

type waterEquivalentSnowRemarkCommand struct {
	locale Locale
}

func newWaterEquivalentSnowRemarkCommand(locale Locale) *waterEquivalentSnowRemarkCommand {
	return &waterEquivalentSnowRemarkCommand{locale: locale}
}

func (c *waterEquivalentSnowRemarkCommand) CanParse(code string) bool {
	return waterEquivSnowRe.MatchString(code)
}

func (c *waterEquivalentSnowRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := waterEquivSnowRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	amount := float64(atoi(m[1])) / 10
	desc := remarkDesc(c.locale, "Remark.Water.Equivalent.Snow.Ground", amount)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeWaterEquivalentSnow,
		Description: desc,
		Raw:         m[0],
		Amount:      &amount,
	})
	return trimAfterRegex(code, waterEquivSnowRe), remarks, nil
}

// ============================================================
// NextForecastByCommand: NXT FCST BY ddhhnnZ
// ============================================================

var nextFcstByRe = regexp.MustCompile(`^NXT FCST BY (\d{2})(\d{2})(\d{2})Z`)

type nextForecastByRemarkCommand struct {
	locale Locale
}

func newNextForecastByRemarkCommand(locale Locale) *nextForecastByRemarkCommand {
	return &nextForecastByRemarkCommand{locale: locale}
}

func (c *nextForecastByRemarkCommand) CanParse(code string) bool {
	return nextFcstByRe.MatchString(code)
}

func (c *nextForecastByRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := nextFcstByRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	day := atoi(m[1])
	hour := atoi(m[2])
	minute := atoi(m[3])
	desc := remarkDesc(c.locale, "Remark.Next.Forecast.By", day, hour, minute)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeNextForecastBy,
		Description: desc,
		Raw:         m[0],
		Day:         &day,
		Hour:        &hour,
		Minute:      &minute,
	})
	return trimAfterRegex(code, nextFcstByRe), remarks, nil
}
