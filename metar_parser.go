package metartafparser

import (
	"regexp"
	"strings"
)

type metarParser struct {
	*abstractParser

	MetarSupplier *metarCommandSupplier
}

func newMetarParser(locale Locale) *metarParser {
	return &metarParser{
		abstractParser: newAbstractParser(locale),
		MetarSupplier:  newMetarCommandSupplier(),
	}
}

var visMNPrefixPattern = regexp.MustCompile(`^(P|M)\d*$`)

func (p *metarParser) Parse(input string) (*Metar, error) {
	tokens := tokenize(input)
	tokens = mergePMVisibilityPrefixes(tokens)
	if len(tokens) == 0 {
		return nil, NewParseError("Empty input")
	}

	idx := 0

	var metarType *MetarType
	if t := parseMetarType(tokens[idx]); t != nil {
		metarType = t
		idx++
	}

	var flags struct {
		amendment, auto, canceled, corrected, nilVal bool
	}
	if idx+1 < len(tokens) && isStation(tokens[idx+1]) {
		if f := findFlags(tokens[idx]); f != nil {
			flags = *f
			idx++
		}
	}

	if idx >= len(tokens) {
		return nil, NewParseError("Unexpected end of input after type/flags")
	}

	station := tokens[idx]
	idx++

	if idx >= len(tokens) {
		return nil, NewParseError("Unexpected end of input after station")
	}

	timeVals := parseDeliveryTime(tokens[idx])
	idx++

	metar := &Metar{
		Type:    metarType,
		Station: station,
		Message: input,
	}
	if timeVals != nil {
		metar.Day = timeVals.Day
		metar.Hour = timeVals.Hour
		metar.Minute = timeVals.Minute
	}
	applyFlagsToMetar(metar, flags)

	for idx < len(tokens) {
		tok := tokens[idx]

		if !p.generalParse(&metar.Container, tok) && !parseFlagsIntoMetar(metar, tok) {
			switch tok {
			case "NOSIG":
				metar.Nosig = boolPtr(true)
			case string(WeatherChangeTypeTEMPO), string(WeatherChangeTypeINTER), string(WeatherChangeTypeBECMG):
				startIdx := idx
				var trendType WeatherChangeType
				switch tok {
				case string(WeatherChangeTypeTEMPO):
					trendType = WeatherChangeTypeTEMPO
				case string(WeatherChangeTypeINTER):
					trendType = WeatherChangeTypeINTER
				case string(WeatherChangeTypeBECMG):
					trendType = WeatherChangeTypeBECMG
				}
				trend := MetarTrend{
					Type: trendType,
				}
				idx = p.parseMetarTrend(idx, &trend, tokens)
				trend.Raw = strings.Join(tokens[startIdx:idx+1], " ")
				metar.Trends = append(metar.Trends, trend)
			case remarkStr:
				parseRemarkIntoMetar(metar, tokens, idx, p.Locale)
				return metar, nil
			default:
				cmd := p.MetarSupplier.Get(tok)
				if cmd != nil {
					err := cmd.Execute(metar, tok)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		idx++
	}

	return metar, nil
}

func mergePMVisibilityPrefixes(tokens []string) []string {
	//nolint:intrange // range would evaluate len(tokens) once; slice mutates inside loop
	for i := 0; i < len(tokens)-1; i++ {
		if visMNPrefixPattern.MatchString(tokens[i]) && visNMRegex.MatchString(tokens[i+1]) {
			if len(tokens[i]) > 1 {
				tokens[i] = tokens[i] + " " + tokens[i+1]
			} else {
				tokens[i] += tokens[i+1]
			}
			tokens = append(tokens[:i+1], tokens[i+2:]...)
		}
	}
	return tokens
}

func (p *metarParser) parseMetarTrend(index int, trend *MetarTrend, tokens []string) int {
	i := index + 1
	for i < len(tokens) && tokens[i] != string(WeatherChangeTypeTEMPO) && tokens[i] != string(WeatherChangeTypeINTER) && tokens[i] != string(WeatherChangeTypeBECMG) && tokens[i] != remarkStr {
		tok := tokens[i]
		if strings.HasPrefix(tok, "FM") || strings.HasPrefix(tok, "TL") || strings.HasPrefix(tok, "AT") {
			tt := MetarTrendTime{}
			prefix := tok[:2]
			switch prefix {
			case "FM":
				tt.Type = TimeIndicatorFM
			case "TL":
				tt.Type = TimeIndicatorTL
			case "AT":
				tt.Type = TimeIndicatorAT
			}
			if len(tok) >= 4 {
				h := atoi(tok[2:4])
				tt.Hour = &h
			}
			if len(tok) >= 6 {
				m := atoi(tok[4:6])
				tt.Minute = &m
			}
			trend.Times = append(trend.Times, tt)
		} else {
			p.generalParse(&trend.Container, tok)
		}
		i++
	}
	return i - 1
}

type flagState struct {
	amendment, auto, canceled, corrected, nilVal bool
}

func findFlags(flag string) *flagState {
	switch flag {
	case "AMD":
		return &flagState{amendment: true}
	case "AUTO":
		return &flagState{auto: true}
	case "CNL":
		return &flagState{canceled: true}
	case "COR":
		return &flagState{corrected: true}
	case "NIL":
		return &flagState{nilVal: true}
	}
	return nil
}

func applyFlagsToMetar(m *Metar, flags flagState) {
	if flags.amendment {
		m.Amendment = boolPtr(true)
	}
	if flags.auto {
		m.Auto = boolPtr(true)
	}
	if flags.canceled {
		m.Canceled = boolPtr(true)
	}
	if flags.corrected {
		m.Corrected = boolPtr(true)
	}
	if flags.nilVal {
		m.Nil = boolPtr(true)
	}
}

func parseFlagsIntoMetar(m *Metar, flag string) bool {
	f := findFlags(flag)
	if f == nil {
		return false
	}
	applyFlagsToMetar(m, *f)
	return true
}

// DeliveryTime stores the parsed day, hour, and minute from a report timestamp.
type DeliveryTime struct {
	// Day of month (1-31).
	Day *int
	// Hour in UTC (0-23).
	Hour *int
	// Minute (0-59).
	Minute *int
}

var deliveryTimeRe = regexp.MustCompile(`^\d{6}Z?$`)

func parseDeliveryTime(s string) *DeliveryTime {
	if !deliveryTimeRe.MatchString(s) {
		return nil
	}
	day := atoi(s[0:2])
	hour := atoi(s[2:4])
	minute := atoi(s[4:6])
	return &DeliveryTime{
		Day:    &day,
		Hour:   &hour,
		Minute: &minute,
	}
}

func parseMetarType(s string) *MetarType {
	switch s {
	case "METAR":
		t := MetarTypeMETAR
		return &t
	case "SPECI":
		t := MetarTypeSPECI
		return &t
	}
	return nil
}

func parseRemarkTokens(tokens []string, index int, locale Locale) ([]Remark, string) {
	remarkCode := strings.Join(tokens[index+1:], " ")
	parser := newRemarkParser(locale)
	remarks := parser.Parse(remarkCode)
	descParts := make([]string, len(remarks))
	for i, r := range remarks {
		if r.Description != nil {
			descParts[i] = *r.Description
		} else {
			descParts[i] = r.Raw
		}
	}
	return remarks, strings.Join(descParts, " ")
}

func parseRemarkIntoMetar(m *Metar, tokens []string, index int, locale Locale) {
	remarks, combined := parseRemarkTokens(tokens, index, locale)
	m.Remarks = remarks
	m.Remark = &combined
}
