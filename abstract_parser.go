package metartafparser

import (
	"regexp"
	"strings"
)

type abstractParser struct {
	Locale         Locale
	CommonSupplier *commonCommandSupplier
}

func newAbstractParser(locale Locale) *abstractParser {
	return &abstractParser{
		Locale:         locale,
		CommonSupplier: newCommonCommandSupplier(),
	}
}

const cavokStr = "CAVOK"

const remarkStr = "RMK"

var intensityRe = regexp.MustCompile(`^(-|\+|VC)`)

func (p *abstractParser) parseWeatherCondition(input string) *WeatherCondition {
	var intensity *Intensity
	if m := intensityRe.FindString(input); m != "" {
		input = input[len(m):]
		switch m {
		case "-":
			v := IntensityLight
			intensity = &v
		case "+":
			v := IntensityHeavy
			intensity = &v
		case "VC":
			v := IntensityInVicinity
			intensity = &v
		}
	} else {
		v := IntensityModerate
		intensity = &v
	}

	var descriptive *Descriptive
	for _, d := range allDescriptiveValues {
		if strings.HasPrefix(input, string(d)) {
			v := d
			descriptive = &v
			input = input[len(v):]
			break
		}
	}

	wc := WeatherCondition{
		Intensity:   intensity,
		Descriptive: descriptive,
	}

	for {
		matched := false
		for _, phen := range allPhenomenonValues {
			if descriptive != nil && string(*descriptive) == string(phen) {
				continue
			}
			phenStr := string(phen)
			if strings.HasPrefix(input, "/"+phenStr) || strings.HasPrefix(input, phenStr) {
				wc.Phenomena = append(wc.Phenomena, phen)
				if strings.HasPrefix(input, "/"+phenStr) {
					input = input[len("/"+phenStr):]
				} else {
					input = input[len(phenStr):]
				}
				matched = true
				break
			}
		}
		if !matched {
			break
		}
	}

	input = strings.ReplaceAll(input, "/", "")
	if input != "" {
		return nil
	}

	if intensity != nil && *intensity == IntensityHeavy &&
		len(wc.Phenomena) == 1 && wc.Phenomena[0] == PhenomenonFunnelCloud {
		wc.Phenomena[0] = PhenomenonTornado
		wc.Intensity = nil
	}

	return &wc
}

func (p *abstractParser) generalParse(c *Container, input string) bool {
	if input == cavokStr {
		c.Cavok = boolPtr(true)
		p6 := ValueIndicatorGreaterThan
		c.Visibility = &Visibility{
			Distance: Distance{
				Indicator: &p6,
				Value:     9999,
				Unit:      DistanceUnitMeters,
			},
		}
		return true
	}

	wc := p.parseWeatherCondition(input)
	if wc != nil && isWeatherConditionValid(*wc) {
		c.WeatherConditions = append(c.WeatherConditions, *wc)
		return true
	}

	for _, cmd := range p.CommonSupplier.commands {
		if !cmd.CanParse(input) {
			continue
		}
		parsed, _, parseErr := cmd.Parse(input)
		if parseErr != nil {
			return false
		}
		mergeContainer(c, parsed)
		return true
	}

	return false
}
