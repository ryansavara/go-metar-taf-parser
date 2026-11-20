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

	minMaxLines := make([][]string, 0, len(lines))
	minMaxLines = append(minMaxLines, firstLine[idx+1:])
	minMaxLines = append(minMaxLines, lines[1:]...)
	p.parseMaxMinTemperatures(taf, minMaxLines)

	for li := 1; li < len(lines); li++ {
		if err := p.parseTAFLine(taf, lines[li]); err != nil {
			return nil, err
		}
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
		t, err := parseDeliveryTime(firstLine[idx])
		if err != nil {
			return "", flagState{}, nil, nil, nil, Validity{}, 0, err
		}
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
	validity, err := parseValidity(validityRaw)
	if err != nil {
		return "", flagState{}, nil, nil, nil, Validity{}, 0, err
	}

	return station, flags, day, hour, minute, validity, idx, nil
}

func (p *tafParser) parseTAFFirstLine(taf *TAF, firstLine []string, startIdx int) {
	for i := startIdx + 1; i < len(firstLine); i++ {
		tok := firstLine[i]
		tafCmd := p.TafSupplier.Get(tok)
		switch {
		case tok == remarkStr:
			parseRemarkIntoTAF(taf, firstLine, i, p.Locale)
			return
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

type tempGroup struct {
	temp       Temperature
	isMax      bool
	isMin      bool
	isMilitary bool
}

func (p *tafParser) parseMaxMinTemperatures(taf *TAF, lines [][]string) {
	var groups []tempGroup

	for _, line := range lines {
		for _, tok := range line {
			if tok == remarkStr {
				break
			}
			switch {
			case strings.HasPrefix(tok, "TX"):
				if t, err := parseTemperature(tok); err == nil {
					groups = append(groups, tempGroup{temp: t, isMax: true})
				}
			case strings.HasPrefix(tok, "TN"):
				if t, err := parseTemperature(tok); err == nil {
					groups = append(groups, tempGroup{temp: t, isMin: true})
				}
			case isMilitaryTempGroup(tok):
				if t, err := parseMilitaryTemperature(tok); err == nil {
					groups = append(groups, tempGroup{temp: t, isMilitary: true})
				}
			}
		}
	}

	assignMaxMinTemperatures(taf, groups)
}

func isMilitaryTempGroup(tok string) bool {
	if len(tok) < 5 || tok[0] != 'T' {
		return false
	}
	if strings.HasPrefix(tok, "TX") || strings.HasPrefix(tok, "TN") || strings.HasPrefix(tok, "TEMPO") || strings.HasPrefix(tok, "TAF") {
		return false
	}
	parts := splitN(tok, "/")
	if len(parts) != 2 {
		return false
	}
	if len(parts[1]) < 4 || parts[1][len(parts[1])-1] != 'Z' {
		return false
	}
	tempPart := parts[0][1:]
	if len(tempPart) == 0 {
		return false
	}
	if tempPart[0] == 'M' {
		tempPart = tempPart[1:]
	}
	for _, c := range tempPart {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func parseMilitaryTemperature(input string) (Temperature, error) {
	parts := splitN(input, "/")
	if len(parts) < 2 {
		return Temperature{}, NewParseError("Invalid military temperature format: " + input)
	}
	if len(parts[1]) < 4 || parts[1][len(parts[1])-1] != 'Z' {
		return Temperature{}, NewParseError("Invalid day/hour in military temperature: " + input)
	}
	t, err := convertTemperature(parts[0][1:])
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

func assignMaxMinTemperatures(taf *TAF, groups []tempGroup) {
	var maxTemp *Temperature
	var minTemp *Temperature

	for _, g := range groups {
		if g.isMax {
			maxTemp = &g.temp
		}
		if g.isMin {
			minTemp = &g.temp
		}
	}

	militaryGroups := make([]Temperature, 0)
	for _, g := range groups {
		if g.isMilitary {
			militaryGroups = append(militaryGroups, g.temp)
		}
	}

	switch len(militaryGroups) {
	case 1:
		if maxTemp == nil {
			maxTemp = &militaryGroups[0]
		} else if minTemp == nil {
			if militaryGroups[0].Temperature < maxTemp.Temperature {
				minTemp = &militaryGroups[0]
			} else {
				minTemp = &militaryGroups[0]
			}
		}
	case 2:
		var higher, lower Temperature
		if militaryGroups[0].Temperature > militaryGroups[1].Temperature {
			higher = militaryGroups[0]
			lower = militaryGroups[1]
		} else {
			higher = militaryGroups[1]
			lower = militaryGroups[0]
		}
		if maxTemp == nil {
			maxTemp = &higher
		}
		if minTemp == nil {
			minTemp = &lower
		}
	}

	if maxTemp != nil {
		taf.MaxTemperature = maxTemp
	}
	if minTemp != nil {
		taf.MinTemperature = minTemp
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

func parseValidity(input string) (Validity, error) {
	parts := splitN(input, "/")
	startDay := atoi(parts[0][0:2])
	startHour := atoi(parts[0][2:])
	endDay := atoi(parts[1][0:2])
	endHour := atoi(parts[1][2:])
	if startDay < 1 || startDay > 31 {
		return Validity{}, NewTimestampOutOfBoundsError("start day out of range: " + input)
	}
	if startHour < 0 || startHour > 23 {
		return Validity{}, NewTimestampOutOfBoundsError("start hour out of range: " + input)
	}
	if endDay < 1 || endDay > 31 {
		return Validity{}, NewTimestampOutOfBoundsError("end day out of range: " + input)
	}
	if endHour < 0 || endHour > 24 {
		return Validity{}, NewTimestampOutOfBoundsError("end hour out of range: " + input)
	}
	return Validity{
		StartDay:  startDay,
		StartHour: startHour,
		EndDay:    endDay,
		EndHour:   endHour,
	}, nil
}

func parseFromValidity(input string) (FMValidity, error) {
	startDay := atoi(input[2:4])
	startHour := atoi(input[4:6])
	startMinutes := atoi(input[6:8])
	if startDay < 1 || startDay > 31 {
		return FMValidity{}, NewTimestampOutOfBoundsError("start day out of range: " + input)
	}
	if startHour < 0 || startHour > 23 {
		return FMValidity{}, NewTimestampOutOfBoundsError("start hour out of range: " + input)
	}
	if startMinutes < 0 || startMinutes > 59 {
		return FMValidity{}, NewTimestampOutOfBoundsError("start minutes out of range: " + input)
	}
	return FMValidity{
		StartDay:     startDay,
		StartHour:    startHour,
		StartMinutes: startMinutes,
	}, nil
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
	inRemark := false
	for i := range len(s) {
		if s[i] != ' ' {
			continue
		}
		rest := s[i+1:]
		if !inRemark && isRemarkTokenStart(rest) {
			inRemark = true
			continue
		}
		if inRemark {
			continue
		}
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

func isRemarkTokenStart(s string) bool {
	return s == remarkStr || strings.HasPrefix(s, remarkStr+" ")
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

func findRemarkIndex(tokens []string) int {
	for i, tok := range tokens {
		if tok == remarkStr {
			return i
		}
	}
	return -1
}

func (p *tafParser) parseTAFLine(taf *TAF, tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}

	rmkIdx := findRemarkIndex(tokens)
	var rawTokens []string
	if rmkIdx >= 0 {
		rawTokens = tokens[:rmkIdx]
	} else {
		rawTokens = tokens
	}

	idx := 0
	var trend TAFTrend

	switch {
	case fmPattern.MatchString(tokens[0]):
		fmVal, err := parseFromValidity(tokens[0])
		if err != nil {
			return err
		}
		trend = TAFTrend{
			BaseTAFTrend: BaseTAFTrend{
				Type:     WeatherChangeTypeFM,
				Raw:      strings.Join(rawTokens, " "),
				Validity: fmVal,
			},
		}
	case strings.HasPrefix(tokens[0], "PROB"):
		if len(tokens) < 2 {
			return nil
		}
		prob := atoi(tokens[0][4:])

		validityIdx := idx + 1
		var foundValidity *Validity
		for i := validityIdx; i < len(tokens); i++ {
			if validityPattern.MatchString(tokens[i]) {
				v, err := parseValidity(tokens[i])
				if err != nil {
					return err
				}
				foundValidity = &v
				break
			}
		}
		if foundValidity == nil {
			return nil
		}

		trend = TAFTrend{
			BaseTAFTrend: BaseTAFTrend{
				Type:        WeatherChangeTypePROB,
				Probability: &prob,
				Raw:         strings.Join(rawTokens, " "),
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
					Raw:         strings.Join(rawTokens, " "),
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
				v, err := parseValidity(tokens[i])
				if err != nil {
					return err
				}
				foundValidity = &v
				break
			}
		}
		if foundValidity == nil {
			return nil
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
			return nil
		}

		trend = TAFTrend{
			BaseTAFTrend: BaseTAFTrend{
				Type:     wt,
				Raw:      strings.Join(rawTokens, " "),
				Validity: *foundValidity,
			},
		}
	}

	p.parseTAFTrend(idx, tokens, &trend, taf)
	taf.Trends = append(taf.Trends, trend)
	return nil
}

func (p *tafParser) parseTAFTrend(index int, tokens []string, trend *TAFTrend, taf *TAF) {
	for i := index; i < len(tokens); i++ {
		tok := tokens[i]
		tafCmd := p.TafSupplier.Get(tok)
		switch {
		case tok == remarkStr:
			parseRemarkIntoTAFTrend(trend, tokens, i, p.Locale)
			if taf.Remark == nil {
				parseRemarkIntoTAF(taf, tokens, i, p.Locale)
			}
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
