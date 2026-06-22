package metartafparser

import (
	"regexp"
)

// ============================================================
// TowerVisibilityCommand: TWR VIS <dist>
// ============================================================

var twrVisRe = regexp.MustCompile(`^TWR VIS ((\d)*( )?(\d?\/?\d))`)

type towerVisibilityRemarkCommand struct {
	locale Locale
}

func newTowerVisibilityRemarkCommand(locale Locale) *towerVisibilityRemarkCommand {
	return &towerVisibilityRemarkCommand{locale: locale}
}

func (c *towerVisibilityRemarkCommand) CanParse(code string) bool { return twrVisRe.MatchString(code) }

func (c *towerVisibilityRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := twrVisRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	dist, err := convertFractionalAmount(m[1])
	if err != nil {
		return code, remarks, err
	}
	desc := remarkDesc(c.locale, "Remark.Tower.Visibility", m[1])
	remarks = append(remarks, Remark{
		Type:        RemarkTypeTowerVisibility,
		Description: desc,
		Raw:         m[0],
		Value:       &dist,
	})
	return trimAfterRegex(code, twrVisRe), remarks, nil
}

// ============================================================
// SurfaceVisibilityCommand: SFC VIS <dist>
// ============================================================

var sfcVisRe = regexp.MustCompile(`^SFC VIS ((\d)*( )?(\d?\/?\d))`)

type surfaceVisibilityRemarkCommand struct {
	locale Locale
}

func newSurfaceVisibilityRemarkCommand(locale Locale) *surfaceVisibilityRemarkCommand {
	return &surfaceVisibilityRemarkCommand{locale: locale}
}

func (c *surfaceVisibilityRemarkCommand) CanParse(code string) bool {
	return sfcVisRe.MatchString(code)
}

func (c *surfaceVisibilityRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := sfcVisRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	dist, err := convertFractionalAmount(m[1])
	if err != nil {
		return code, remarks, err
	}
	desc := remarkDesc(c.locale, "Remark.Surface.Visibility", m[1])
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSurfaceVisibility,
		Description: desc,
		Raw:         m[0],
		Value:       &dist,
	})
	return trimAfterRegex(code, sfcVisRe), remarks, nil
}

// ============================================================
// PrevailingVisibilityCommand: VIS <min>V<max>
// ============================================================

var prevailVisRe = regexp.MustCompile(`^VIS ((\d)*( )?(\d?\/?\d))V((\d)*( )?(\d?\/?\d))`)

type prevailingVisibilityRemarkCommand struct {
	locale Locale
}

func newPrevailingVisibilityRemarkCommand(locale Locale) *prevailingVisibilityRemarkCommand {
	return &prevailingVisibilityRemarkCommand{locale: locale}
}

func (c *prevailingVisibilityRemarkCommand) CanParse(code string) bool {
	return prevailVisRe.MatchString(code)
}

func (c *prevailingVisibilityRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := prevailVisRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	minVis, err := convertFractionalAmount(m[1])
	if err != nil {
		return code, remarks, err
	}
	maxVis, err := convertFractionalAmount(m[5])
	if err != nil {
		return code, remarks, err
	}
	desc := remarkDesc(c.locale, "Remark.Variable.Prevailing.Visibility", m[1], m[5])
	remarks = append(remarks, Remark{
		Type:        RemarkTypePrevailingVisibility,
		Description: desc,
		Raw:         m[0],
		Min:         float64Ptr(minVis),
		Max:         float64Ptr(maxVis),
	})
	return trimAfterRegex(code, prevailVisRe), remarks, nil
}

// ============================================================
// SecondLocationVisibilityCommand: VIS <dist> <location>
// ============================================================

var secondLocVisRe = regexp.MustCompile(`^VIS ((\d)*( )?(\d?\/?\d)) (\w+)`)

type secondLocationVisibilityRemarkCommand struct {
	locale Locale
}

func newSecondLocationVisibilityRemarkCommand(locale Locale) *secondLocationVisibilityRemarkCommand {
	return &secondLocationVisibilityRemarkCommand{locale: locale}
}

func (c *secondLocationVisibilityRemarkCommand) CanParse(code string) bool {
	if prevailVisRe.MatchString(code) {
		return false
	}
	return secondLocVisRe.MatchString(code)
}

func (c *secondLocationVisibilityRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := secondLocVisRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	dist, err := convertFractionalAmount(m[1])
	if err != nil {
		return code, remarks, err
	}
	desc := remarkDesc(c.locale, "Remark.Second.Location.Visibility", m[1], m[5])
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSecondLocationVisibility,
		Description: desc,
		Raw:         m[0],
		Value:       &dist,
	})
	return trimAfterRegex(code, secondLocVisRe), remarks, nil
}

// ============================================================
// SectorVisibilityCommand: VIS <dir> <dist>
// ============================================================

var sectorVisRe = regexp.MustCompile(`^VIS ([A-Z]{1,2}) ((\d)*( )?(\d?\/?\d))`)

type sectorVisibilityRemarkCommand struct {
	locale Locale
}

func newSectorVisibilityRemarkCommand(locale Locale) *sectorVisibilityRemarkCommand {
	return &sectorVisibilityRemarkCommand{locale: locale}
}

func (c *sectorVisibilityRemarkCommand) CanParse(code string) bool {
	if prevailVisRe.MatchString(code) || secondLocVisRe.MatchString(code) {
		return false
	}
	return sectorVisRe.MatchString(code)
}

func (c *sectorVisibilityRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := sectorVisRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	dist, err := convertFractionalAmount(m[2])
	if err != nil {
		return code, remarks, err
	}
	dirLoc := locStr(c.locale, "Converter."+m[1])
	desc := remarkDesc(c.locale, "Remark.Sector.Visibility", dirLoc, m[2])
	remarks = append(remarks, Remark{
		Type:        RemarkTypeSectorVisibility,
		Description: desc,
		Raw:         m[0],
		Value:       &dist,
	})
	if m[1] != "" {
		d := m[1]
		remarks[len(remarks)-1].Direction = &d
	}
	return trimAfterRegex(code, sectorVisRe), remarks, nil
}

// ============================================================
// CeilingHeightCommand: CIG hhhVhhh
// ============================================================

var ceilingHeightRe = regexp.MustCompile(`^CIG (\d{3})V(\d{3})`)

type ceilingHeightRemarkCommand struct {
	locale Locale
}

func newCeilingHeightRemarkCommand(locale Locale) *ceilingHeightRemarkCommand {
	return &ceilingHeightRemarkCommand{locale: locale}
}

func (c *ceilingHeightRemarkCommand) CanParse(code string) bool {
	return ceilingHeightRe.MatchString(code)
}

func (c *ceilingHeightRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := ceilingHeightRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	minVal := atoi(m[1]) * 100
	maxVal := atoi(m[2]) * 100
	desc := remarkDesc(c.locale, "Remark.Ceiling.Height", minVal, maxVal)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeCeilingHeight,
		Description: desc,
		Raw:         m[0],
		Min:         float64Ptr(float64(minVal)),
		Max:         float64Ptr(float64(maxVal)),
	})
	return trimAfterRegex(code, ceilingHeightRe), remarks, nil
}

// ============================================================
// ObscurationCommand: <phen> <qty>hgt
// ============================================================

var obscurationRe = regexp.MustCompile(`^([A-Z]{2}) ([A-Z]{3})(\d{3})`)

type obscurationRemarkCommand struct {
	locale Locale
}

func newObscurationRemarkCommand(locale Locale) *obscurationRemarkCommand {
	return &obscurationRemarkCommand{locale: locale}
}

func (c *obscurationRemarkCommand) CanParse(code string) bool { return obscurationRe.MatchString(code) }

func (c *obscurationRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := obscurationRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	if !isValidCloudQuantity(m[2]) {
		return code, remarks, &commandExecutionError{}
	}
	if !isValidPhenomenon(m[1]) {
		return code, remarks, &commandExecutionError{}
	}
	phenLoc := locStr(c.locale, "Phenomenon."+m[1])
	qtyLoc := locStr(c.locale, "CloudQuantity."+m[2])
	height := atoi(m[3]) * 100
	desc := remarkDesc(c.locale, "Remark.Obscuration", qtyLoc, height, phenLoc)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeObscuration,
		Description: desc,
		Raw:         m[0],
		Min:         float64Ptr(float64(height)),
	})
	return trimAfterRegex(code, obscurationRe), remarks, nil
}

// ============================================================
// VariableSkyHeightCommand: <qty>hgt V <qty>
// ============================================================

var varSkyHgtRe = regexp.MustCompile(`^([A-Z]{3})(\d{3}) V ([A-Z]{3})`)

type variableSkyHeightRemarkCommand struct {
	locale Locale
}

func newVariableSkyHeightRemarkCommand(locale Locale) *variableSkyHeightRemarkCommand {
	return &variableSkyHeightRemarkCommand{locale: locale}
}

func (c *variableSkyHeightRemarkCommand) CanParse(code string) bool {
	return varSkyHgtRe.MatchString(code)
}

func (c *variableSkyHeightRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := varSkyHgtRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	height := atoi(m[2]) * 100
	qty1Loc := locStr(c.locale, "CloudQuantity."+m[1])
	qty2Loc := locStr(c.locale, "CloudQuantity."+m[3])
	desc := remarkDesc(c.locale, "Remark.Variable.Sky.Condition.Height", height, qty1Loc, qty2Loc)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeVariableSkyHeight,
		Description: desc,
		Raw:         m[0],
		Min:         float64Ptr(float64(height)),
	})
	return trimAfterRegex(code, varSkyHgtRe), remarks, nil
}

// ============================================================
// VariableSkyCommand: <qty> V <qty>
// ============================================================

var varSkyRe = regexp.MustCompile(`^([A-Z]{3}) V ([A-Z]{3})`)

type variableSkyRemarkCommand struct {
	locale Locale
}

func newVariableSkyRemarkCommand(locale Locale) *variableSkyRemarkCommand {
	return &variableSkyRemarkCommand{locale: locale}
}

func (c *variableSkyRemarkCommand) CanParse(code string) bool {
	if varSkyHgtRe.MatchString(code) {
		return false
	}
	return varSkyRe.MatchString(code)
}

func (c *variableSkyRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := varSkyRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	qty1Loc := locStr(c.locale, "CloudQuantity."+m[1])
	qty2Loc := locStr(c.locale, "CloudQuantity."+m[2])
	desc := remarkDesc(c.locale, "Remark.Variable.Sky.Condition.0", qty1Loc, qty2Loc)
	remarks = append(remarks, Remark{
		Type:        RemarkTypeVariableSky,
		Description: desc,
		Raw:         m[0],
	})
	return trimAfterRegex(code, varSkyRe), remarks, nil
}

// ============================================================
// CeilingSecondLocationCommand: CIG hhh <loc>
// ============================================================

var ceilingSecondLocRe = regexp.MustCompile(`^CIG (\d{3}) (\w+)`)

type ceilingSecondLocationRemarkCommand struct {
	locale Locale
}

func newCeilingSecondLocationRemarkCommand(locale Locale) *ceilingSecondLocationRemarkCommand {
	return &ceilingSecondLocationRemarkCommand{locale: locale}
}

func (c *ceilingSecondLocationRemarkCommand) CanParse(code string) bool {
	if ceilingHeightRe.MatchString(code) {
		return false
	}
	return ceilingSecondLocRe.MatchString(code)
}

func (c *ceilingSecondLocationRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	m := ceilingSecondLocRe.FindStringSubmatch(code)
	if m == nil {
		return code, remarks, errMatchNotFound
	}
	height := atoi(m[1]) * 100
	desc := remarkDesc(c.locale, "Remark.Ceiling.Second.Location", height, m[2])
	remarks = append(remarks, Remark{
		Type:        RemarkTypeCeilingSecondLocation,
		Description: desc,
		Raw:         m[0],
		Min:         float64Ptr(float64(height)),
	})
	return trimAfterRegex(code, ceilingSecondLocRe), remarks, nil
}
