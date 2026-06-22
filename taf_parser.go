package metartafparser

import (
	"regexp"
	"strings"
)

const tafStr = "TAF"

type tafParser struct {
	*abstractParser

	TafSupplier *tafCommandSupplier
}

func newTAFParser(locale Locale) *tafParser {
	return &tafParser{
		abstractParser: newAbstractParser(locale),
		TafSupplier:    newTAFCommandSupplier(),
	}
}

var validityPattern = regexp.MustCompile(`^\d{4}\/\d{4}$`)
var partialPattern = regexp.MustCompile(`^PART (\d) OF (\d) `)
var fmPattern = regexp.MustCompile(`^FM(\d{2})(\d{2})(\d{2})$`)

func (p *tafParser) Parse(input string) (*TAF, error) {
	err := p.checkPartial(input)
	if err != nil {
		return nil, err
	}

	lines := p.extractLinesTokens(input)
	if len(lines) == 0 || len(lines[0]) == 0 {
		return nil, NewParseError("Empty TAF input")
	}

	firstLine := lines[0]

	station, flags, day, hour, minute, validity, idx, err := p.parseTAFHeader(firstLine)
	if err != nil {
		return nil, err
	}

	taf := &TAF{
		Station:    station,
		Day:        day,
		Hour:       hour,
		Minute:     minute,
		Validity:   validity,
		Message:    input,
		InitialRaw: strings.Join(firstLine, " "),
	}

	applyFlagsToTAF(taf, flags)

	p.parseTAFFirstLine(taf, firstLine, idx)

	minMaxLines := [][]string{firstLine[idx+1:]}
	if len(lines) > 1 {
		minMaxLines = append(minMaxLines, lines[len(lines)-1])
	}
	p.parseMaxMinTemperatures(taf, minMaxLines)

	for li := 1; li < len(lines); li++ {
		p.parseTAFLine(taf, lines[li])
	}

	return taf, nil
}

func (p *tafParser) checkPartial(input string) error {
	m := partialPattern.FindStringSubmatch(input)
	if m != nil {
		part := atoi(m[1])
		total := atoi(m[2])
		return NewPartialWeatherStatementError(m[0], part, total)
	}
	return nil
}

func (p *tafParser) parseTAFHeader(firstLine []string) (string, flagState, *int, *int, *int, Validity, int, error) {
	idx := 0
	var flags flagState

	if idx < len(firstLine) && firstLine[idx] == tafStr {
		idx++
	}
	if idx+1 < len(firstLine) && firstLine[idx] != tafStr && firstLine[idx+1] == tafStr {
		idx += 2
	} else if idx < len(firstLine) && firstLine[idx] == tafStr {
		idx++
	}
	if idx < len(firstLine) {
		if f := findFlags(firstLine[idx]); f != nil {
			flags = *f
			idx++
		}
	}
	if idx < len(firstLine) && firstLine[idx] == tafStr {
		idx++
	}
	if idx < len(firstLine) {
		if f := findFlags(firstLine[idx]); f != nil {
			flags = mergeFlags(flags, *f)
			idx++
		}
	}

	if idx >= len(firstLine) {
		return "", flagState{}, nil, nil, nil, Validity{}, 0, NewParseError("No station found")
	}

	station := firstLine[idx]
	idx++

	var day, hour, minute *int
	if idx < len(firstLine) {
		t := parseDeliveryTime(firstLine[idx])
		if t != nil {
			day = t.Day
			hour = t.Hour
			minute = t.Minute
			idx++
		}
	}

	if idx >= len(firstLine) {
		return "", flagState{}, nil, nil, nil, Validity{}, 0, NewParseError("No validity found")
	}

	validityRaw := firstLine[idx]
	if !validityPattern.MatchString(validityRaw) {
		return "", flagState{}, nil, nil, nil, Validity{}, 0, NewParseError("Invalid validity format: " + validityRaw)
	}
	validity := parseValidity(validityRaw)

	return station, flags, day, hour, minute, validity, idx, nil
}

func (p *tafParser) parseTAFFirstLine(taf *TAF, firstLine []string, startIdx int) {
	for i := startIdx + 1; i < len(firstLine); i++ {
		tok := firstLine[i]
		tafCmd := p.TafSupplier.Get(tok)
		switch {
		case tok == remarkStr:
			parseRemarkIntoTAF(taf, firstLine, i, p.Locale)
		case tafCmd != nil:
			err := tafCmd.Execute(&taf.Container, tok)
			if err != nil {
				p.generalParse(&taf.Container, tok)
			}
		default:
			p.generalParse(&taf.Container, tok)
			parseFlagsIntoTAF(taf, tok)
		}
	}
}

func (p *tafParser) parseMaxMinTemperatures(taf *TAF, lines [][]string) {
	for _, line := range lines {
		for _, tok := range line {
			if tok == remarkStr {
				break
			}
			if strings.HasPrefix(tok, "TX") {
				t, err := parseTemperature(tok)
				if err == nil {
					taf.MaxTemperature = &t
				}
			} else if strings.HasPrefix(tok, "TN") {
				t, err := parseTemperature(tok)
				if err == nil {
					taf.MinTemperature = &t
				}
			}
		}
	}
}

func parseTemperature(input string) (Temperature, error) {
	parts := splitN(input, "/")
	if len(parts) < 2 {
		return Temperature{}, NewParseError("Invalid temperature format: " + input)
	}
	if len(parts[1]) < 4 {
		return Temperature{}, NewParseError("Invalid day/hour in temperature: " + input)
	}
	t, err := convertTemperature(parts[0][2:])
	if err != nil {
		return Temperature{}, err
	}
	day := atoi(parts[1][0:2])
	hour := atoi(parts[1][2:4])
	return Temperature{
		Temperature: t,
		Day:         day,
		Hour:        hour,
	}, nil
}

func parseValidity(input string) Validity {
	parts := splitN(input, "/")
	return Validity{
		StartDay:  atoi(parts[0][0:2]),
		StartHour: atoi(parts[0][2:]),
		EndDay:    atoi(parts[1][0:2]),
		EndHour:   atoi(parts[1][2:]),
	}
}

func parseFromValidity(input string) FMValidity {
	return FMValidity{
		StartDay:     atoi(input[2:4]),
		StartHour:    atoi(input[4:6]),
		StartMinutes: atoi(input[6:8]),
	}
}

func (p *tafParser) extractLinesTokens(tafCode string) [][]string {
	singleLine := strings.ReplaceAll(tafCode, "\n", " ")
	cleanLine := regexp.MustCompile(`\s{2,}`).ReplaceAllString(singleLine, " ")

	parts := splitTAFTrendLines(cleanLine)

	result := make([][]string, 0, len(parts))
	for _, line := range parts {
		tokens := tokenize(line)
		tokens = mergePMVisibilityPrefixes(tokens)
		result = append(result, tokens)
	}
	return result
}

func splitTAFTrendLines(s string) []string {
	var parts []string
	start := 0
	for i := range len(s) {
		if s[i] != ' ' {
			continue
		}
		rest := s[i+1:]
		if isTAFTrendStart(rest) {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])

	var joined []string
	for i := 0; i < len(parts); i++ {
		if probPattern.MatchString(parts[i]) && i+1 < len(parts) {
			next := parts[i+1]
			if strings.HasPrefix(next, "TEMPO") || strings.HasPrefix(next, "INTER") {
				joined = append(joined, parts[i]+" "+next)
				i++
				continue
			}
		}
		joined = append(joined, parts[i])
	}
	return joined
}

func isTAFTrendStart(s string) bool {
	if s == string(WeatherChangeTypeTEMPO) || s == string(WeatherChangeTypeINTER) || s == string(WeatherChangeTypeBECMG) ||
		strings.HasPrefix(s, string(WeatherChangeTypeTEMPO)+" ") || strings.HasPrefix(s, string(WeatherChangeTypeINTER)+" ") || strings.HasPrefix(s, string(WeatherChangeTypeBECMG)+" ") {
		return true
	}
	if strings.HasPrefix(s, "FM") && len(s) >= 3 {
		if len(s) >= 4 && s[2] >= 'A' && s[2] <= 'Z' && s[3] >= 'A' && s[3] <= 'Z' && (len(s) == 4 || s[4] == ' ') {
			return false
		}
		return true
	}
	if strings.HasPrefix(s, "PROB") {
		if len(s) >= 7 && isDigitsOnly(s[4:6]) && (len(s) == 6 || s[6] == ' ') {
			return true
		}
	}
	return false
}

func isDigitsOnly(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

var probPattern = regexp.MustCompile(`^PROB\d{2}$`)

func (p *tafParser) parseTAFLine(taf *TAF, tokens []string) {
	if len(tokens) == 0 {
		return
	}

	idx := 0
	var trend TAFTrend

	switch {
	case fmPattern.MatchString(tokens[0]):
		fmVal := parseFromValidity(tokens[0])
		trend = TAFTrend{
			BaseTAFTrend: BaseTAFTrend{
				Type:     WeatherChangeTypeFM,
				Raw:      strings.Join(tokens, " "),
				Validity: fmVal,
			},
		}
	case strings.HasPrefix(tokens[0], "PROB"):
		if len(tokens) < 2 {
			return
		}
		prob := atoi(tokens[0][4:])

		validityIdx := idx + 1
		var foundValidity *Validity
		for i := validityIdx; i < len(tokens); i++ {
			if validityPattern.MatchString(tokens[i]) {
				v := parseValidity(tokens[i])
				foundValidity = &v
				break
			}
		}
		if foundValidity == nil {
			return
		}

		trend = TAFTrend{
			BaseTAFTrend: BaseTAFTrend{
				Type:        WeatherChangeTypePROB,
				Probability: &prob,
				Raw:         strings.Join(tokens, " "),
				Validity:    *foundValidity,
			},
		}

		if len(tokens) > 1 && (tokens[1] == "TEMPO" || tokens[1] == "INTER") {
			var wt WeatherChangeType
			if tokens[1] == "TEMPO" {
				wt = WeatherChangeTypeTEMPO
			} else {
				wt = WeatherChangeTypeINTER
			}
			trend = TAFTrend{
				BaseTAFTrend: BaseTAFTrend{
					Type:        wt,
					Probability: &prob,
					Raw:         strings.Join(tokens, " "),
					Validity:    *foundValidity,
				},
			}
			idx = 2
		} else {
			idx = 1
		}
	default:
		validityIdx := idx + 1
		var foundValidity *Validity
		for i := validityIdx; i < len(tokens); i++ {
			if validityPattern.MatchString(tokens[i]) {
				v := parseValidity(tokens[i])
				foundValidity = &v
				break
			}
		}
		if foundValidity == nil {
			return
		}

		var wt WeatherChangeType
		switch tokens[0] {
		case "TEMPO":
			wt = WeatherChangeTypeTEMPO
		case "INTER":
			wt = WeatherChangeTypeINTER
		case "BECMG":
			wt = WeatherChangeTypeBECMG
		default:
			return
		}

		trend = TAFTrend{
			BaseTAFTrend: BaseTAFTrend{
				Type:     wt,
				Raw:      strings.Join(tokens, " "),
				Validity: *foundValidity,
			},
		}
	}

	p.parseTAFTrend(idx, tokens, &trend)
	taf.Trends = append(taf.Trends, trend)
}

func (p *tafParser) parseTAFTrend(index int, tokens []string, trend *TAFTrend) {
	for i := index; i < len(tokens); i++ {
		tok := tokens[i]
		tafCmd := p.TafSupplier.Get(tok)
		switch {
		case tok == remarkStr:
			parseRemarkIntoTAFTrend(trend, tokens, i, p.Locale)
			return
		case validityPattern.MatchString(tok):
			continue
		case tafCmd != nil:
			err := tafCmd.Execute(&trend.Container, tok)
			if err != nil {
				p.generalParse(&trend.Container, tok)
			}
		default:
			p.generalParse(&trend.Container, tok)
		}
	}
}

func mergeFlags(a, b flagState) flagState {
	return flagState{
		amendment: a.amendment || b.amendment,
		auto:      a.auto || b.auto,
		canceled:  a.canceled || b.canceled,
		corrected: a.corrected || b.corrected,
		nilVal:    a.nilVal || b.nilVal,
	}
}

func applyFlagsToTAF(t *TAF, flags flagState) {
	if flags.amendment {
		t.Amendment = boolPtr(true)
	}
	if flags.auto {
		t.Auto = boolPtr(true)
	}
	if flags.canceled {
		t.Canceled = boolPtr(true)
	}
	if flags.corrected {
		t.Corrected = boolPtr(true)
	}
	if flags.nilVal {
		t.Nil = boolPtr(true)
	}
}

func parseFlagsIntoTAF(t *TAF, flag string) bool {
	f := findFlags(flag)
	if f == nil {
		return false
	}
	applyFlagsToTAF(t, *f)
	return true
}

func parseRemarkIntoTAF(t *TAF, tokens []string, index int, locale Locale) {
	remarks, combined := parseRemarkTokens(tokens, index, locale)
	t.Remarks = remarks
	t.Remark = &combined
}

func parseRemarkIntoTAFTrend(trend *TAFTrend, tokens []string, index int, locale Locale) {
	remarks, combined := parseRemarkTokens(tokens, index, locale)
	trend.Remarks = remarks
	trend.Remark = &combined
}
