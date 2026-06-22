package metartafparser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var errMatchNotFound = errors.New("match not found")

type remarkCommander interface {
	CanParse(code string) bool
	Execute(code string, remarks []Remark) (string, []Remark, error)
}

// --- helper for remark commands ---

func remarkDesc(locale Locale, path string, args ...any) *string {
	return formatMsg(localeGet(path, locale), args...)
}

func locStr(locale Locale, path string) string {
	s := localeGet(path, locale)
	if s == nil {
		return ""
	}
	return *s
}

func trimAfterRegex(code string, re *regexp.Regexp) string {
	loc := re.FindStringIndex(code)
	if loc == nil {
		return strings.TrimSpace(code)
	}
	return strings.TrimSpace(code[loc[1]:])
}

// ============================================================
// DefaultCommand - catch-all, always canParse = true
// ============================================================

type defaultRemarkCommand struct {
	locale Locale
}

func newDefaultRemarkCommand(locale Locale) *defaultRemarkCommand {
	return &defaultRemarkCommand{locale: locale}
}

func (c *defaultRemarkCommand) CanParse(_ string) bool { return true }

func (c *defaultRemarkCommand) Execute(code string, remarks []Remark) (string, []Remark, error) {
	parts := splitFields(code, 1)
	token := parts[0]
	var rem *string
	if len(parts) > 1 {
		rem = localeGet("Remark."+token, c.locale)
	}

	if isKnownRemarkType(token) {
		remarks = append(remarks, Remark{
			Type:        RemarkType(token),
			Description: rem,
			Raw:         token,
		})
	} else {
		if len(remarks) > 0 && remarks[len(remarks)-1].Type == RemarkTypeUnknown {
			remarks[len(remarks)-1].Raw = remarks[len(remarks)-1].Raw + " " + token
		} else {
			remarks = append(remarks, Remark{
				Type: RemarkTypeUnknown,
				Raw:  token,
			})
		}
	}

	remaining := ""
	if len(parts) > 1 {
		remaining = parts[1]
	}
	return remaining, remarks, nil
}

// ============================================================
// remarkCommandSupplier
// ============================================================

type remarkCommandSupplier struct {
	defaultCommand *defaultRemarkCommand
	commands       []remarkCommander
}

func newRemarkCommandSupplier(locale Locale) *remarkCommandSupplier {
	return &remarkCommandSupplier{
		defaultCommand: newDefaultRemarkCommand(locale),
		commands: []remarkCommander{
			newWindPeakRemarkCommand(locale),
			newWindShiftFropaRemarkCommand(locale),
			newWindShiftRemarkCommand(locale),
			newTowerVisibilityRemarkCommand(locale),
			newSurfaceVisibilityRemarkCommand(locale),
			newPrevailingVisibilityRemarkCommand(locale),
			newSecondLocationVisibilityRemarkCommand(locale),
			newSectorVisibilityRemarkCommand(locale),
			newTornadicActivityBegEndRemarkCommand(locale),
			newTornadicActivityBegRemarkCommand(locale),
			newTornadicActivityEndRemarkCommand(locale),
			newPrecipitationBegEndRemarkCommand(locale),
			newPrecipitationBegRemarkCommand(locale),
			newPrecipitationEndRemarkCommand(locale),
			newThunderStormLocationMovingRemarkCommand(locale),
			newThunderStormLocationRemarkCommand(locale),
			newSmallHailSizeRemarkCommand(locale),
			newHailSizeRemarkCommand(locale),
			newSnowPelletsRemarkCommand(locale),
			newVirgaDirectionRemarkCommand(locale),
			newCeilingHeightRemarkCommand(locale),
			newObscurationRemarkCommand(locale),
			newVariableSkyHeightRemarkCommand(locale),
			newVariableSkyRemarkCommand(locale),
			newCeilingSecondLocationRemarkCommand(locale),
			newSeaLevelPressureRemarkCommand(locale),
			newSnowIncreaseRemarkCommand(locale),
			newHourlyMaximumMinimumTemperatureRemarkCommand(locale),
			newHourlyMaximumTemperatureRemarkCommand(locale),
			newHourlyMinimumTemperatureRemarkCommand(locale),
			newHourlyPrecipitationAmountRemarkCommand(locale),
			newHourlyTemperatureDewPointRemarkCommand(locale),
			newHourlyPressureRemarkCommand(locale),
			newIceAccretionRemarkCommand(locale),
			newPrecipitationAmount36HourRemarkCommand(locale),
			newPrecipitationAmount24HourRemarkCommand(locale),
			newSnowDepthRemarkCommand(locale),
			newSunshineDurationRemarkCommand(locale),
			newWaterEquivalentSnowRemarkCommand(locale),
			newNextForecastByRemarkCommand(locale),
		},
	}
}

func (s *remarkCommandSupplier) Get(code string) (remarkCommander, error) {
	for _, cmd := range s.commands {
		if cmd.CanParse(code) {
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("no command found for: %s", code) //nolint:err113
}
