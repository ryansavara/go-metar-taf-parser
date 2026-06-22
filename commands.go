package metartafparser

import (
	"regexp"
	"strconv"
)

// commonCommander operates on a generic Container.
type commonCommander interface {
	Identifier() string
	Parse(input string) (*Container, string, error)
	CanParse(input string) bool
}

// metarCommander operates on a Metar.
type metarCommander interface {
	CanParse(s string) bool
	Execute(m *Metar, s string) error
}

// tafCommander operates on a Container (could be TAF or TAFTrend).
type tafCommander interface {
	CanParse(s string) bool
	Execute(c *Container, s string) error
}

// --- cloudCommand ---

var cloudRegex = regexp.MustCompile(`^([A-Z]{3})(?:\/{3}|(\d{3}))?(?:\/{3}|(?:([A-Z]{2,3})(?:\/([A-Z]{2,3}))?))?$`)

type cloudCommand struct{}

func isValidCloudQuantity(s string) bool {
	switch CloudQuantity(s) {
	case CloudQuantitySKC, CloudQuantityFEW, CloudQuantitySCT, CloudQuantityBKN, CloudQuantityOVC, CloudQuantityNSC:
		return true
	}
	return false
}

func isValidPhenomenon(s string) bool {
	switch Phenomenon(s) {
	case PhenomenonRain, PhenomenonDrizzle, PhenomenonSnow, PhenomenonSnowGrains,
		PhenomenonIcePellets, PhenomenonIceCrystals, PhenomenonHail, PhenomenonSmallHail,
		PhenomenonUnknownPrecip, PhenomenonFog, PhenomenonVolcanicAsh, PhenomenonMist,
		PhenomenonHaze, PhenomenonWidespreadDust, PhenomenonSmoke, PhenomenonSand,
		PhenomenonSpray, PhenomenonSquall, PhenomenonSandWhirls, PhenomenonThunderstorm,
		PhenomenonDuststorm, PhenomenonSandstorm, PhenomenonFunnelCloud, PhenomenonNoSignificantWeather:
		return true
	}
	return false
}

func isValidCloudType(s string) bool {
	switch CloudType(s) {
	case CloudTypeCB, CloudTypeTCU, CloudTypeCI, CloudTypeCC, CloudTypeCS, CloudTypeAC, CloudTypeST, CloudTypeCU, CloudTypeAS, CloudTypeNS, CloudTypeSC:
		return true
	}
	return false
}

func (cloudCommand) Identifier() string { return "cloud" }

func (cloudCommand) Parse(input string) (*Container, string, error) {
	m := cloudRegex.FindStringSubmatch(input)
	if m == nil {
		return nil, input, errCommandNotHandled
	}
	if !isValidCloudQuantity(m[1]) {
		return nil, input, errCommandNotHandled
	}
	if m[3] != "" && !isValidCloudType(m[3]) {
		return nil, input, errCommandNotHandled
	}
	if m[4] != "" && !isValidCloudType(m[4]) {
		return nil, input, errCommandNotHandled
	}

	cl := Cloud{Quantity: CloudQuantity(m[1])}
	if m[2] != "" {
		h, _ := strconv.Atoi(m[2])
		h *= 100
		cl.Height = &h
	}
	if m[3] != "" {
		t := CloudType(m[3])
		cl.Type = &t
	}
	if m[4] != "" {
		st := CloudType(m[4])
		cl.SecondaryType = &st
	}

	return &Container{Clouds: []Cloud{cl}}, "", nil
}

func (cloudCommand) CanParse(input string) bool {
	if input == "NSW" || input == "NIL" {
		return false
	}
	m := cloudRegex.FindStringSubmatch(input)
	if m == nil {
		return false
	}
	if !isValidCloudQuantity(m[1]) {
		return false
	}
	if m[3] != "" && !isValidCloudType(m[3]) {
		return false
	}
	if m[4] != "" && !isValidCloudType(m[4]) {
		return false
	}
	return true
}

// --- mainVisibilityCommand ---

var mainVisRegex = regexp.MustCompile(`^(\d{4})(|NDV)$`)

type mainVisibilityCommand struct{}

func (mainVisibilityCommand) Identifier() string { return "mainVisibility" }

func (mainVisibilityCommand) Parse(input string) (*Container, string, error) {
	m := mainVisRegex.FindStringSubmatch(input)
	if m == nil {
		return nil, input, errCommandNotHandled
	}

	d, err := convertVisibility(m[1])
	if err != nil {
		return nil, input, err
	}
	if m[2] == "NDV" {
		d.Ndv = true
	}
	return &Container{Visibility: &Visibility{Distance: d}}, "", nil
}

func (mainVisibilityCommand) CanParse(input string) bool {
	return mainVisRegex.MatchString(input)
}

// --- windCommand ---

var windRegex = regexp.MustCompile(`^(VRB|000|[0-3]\d{2})(\d{2})G?(\d{2,3})?(KT|MPS|KM\/H)?`)

type windCommand struct{}

func (windCommand) Identifier() string { return "wind" }

func (windCommand) Parse(input string) (*Container, string, error) {
	m := windRegex.FindStringSubmatch(input)
	if m == nil {
		return nil, input, errCommandNotHandled
	}

	unit := SpeedUnit("KT")
	if m[4] != "" {
		unit = SpeedUnit(m[4])
	}

	speed, _ := strconv.Atoi(m[2])
	w := Wind{
		Speed:     speed,
		Direction: degreesToCardinal(m[1]),
		Unit:      unit,
	}
	if m[1] != "VRB" {
		d, _ := strconv.Atoi(m[1])
		w.Degrees = &d
	}
	if m[3] != "" {
		g, _ := strconv.Atoi(m[3])
		w.Gust = &g
	}

	return &Container{Wind: &w}, "", nil
}

func (windCommand) CanParse(input string) bool {
	return windRegex.MatchString(input)
}

// --- windVariationCommand ---

var windVarRegex = regexp.MustCompile(`^(\d{3})V(\d{3})`)

type windVariationCommand struct{}

func (windVariationCommand) Identifier() string { return "windVariation" }

func (windVariationCommand) Parse(input string) (*Container, string, error) {
	m := windVarRegex.FindStringSubmatch(input)
	if m == nil {
		return nil, input, errCommandNotHandled
	}
	minV, _ := strconv.Atoi(m[1])
	maxV, _ := strconv.Atoi(m[2])
	return &Container{
		Wind: &Wind{
			MinVariation: &minV,
			MaxVariation: &maxV,
		},
	}, "", nil
}

func (windVariationCommand) CanParse(input string) bool {
	return windVarRegex.MatchString(input)
}

// --- windShearCommand ---

var windShearRegex = regexp.MustCompile(`^WS(\d{3})\/(\w{3})(\d{2})G?(\d{2,3})?(KT|MPS|KM\/H)`)

type windShearCommand struct{}

func (windShearCommand) Identifier() string { return "windShear" }

func (windShearCommand) Parse(input string) (*Container, string, error) {
	m := windShearRegex.FindStringSubmatch(input)
	if m == nil {
		return nil, input, errCommandNotHandled
	}

	unit := SpeedUnit(m[5])
	speed, _ := strconv.Atoi(m[3])
	height, _ := strconv.Atoi(m[1])
	height *= 100

	ws := WindShear{
		Speed:     speed,
		Direction: degreesToCardinal(m[2]),
		Unit:      unit,
		Height:    height,
	}
	if m[2] != "VRB" {
		d, _ := strconv.Atoi(m[2])
		ws.Degrees = &d
	}
	if m[4] != "" {
		g, _ := strconv.Atoi(m[4])
		ws.Gust = &g
	}

	return &Container{WindShear: &ws}, "", nil
}

func (windShearCommand) CanParse(input string) bool {
	return windShearRegex.MatchString(input)
}

// --- verticalVisibilityCommand ---

var vertVisRegex = regexp.MustCompile(`^VV(\d{3})$`)

type verticalVisibilityCommand struct{}

func (verticalVisibilityCommand) Identifier() string { return "verticalVisibility" }

func (verticalVisibilityCommand) Parse(input string) (*Container, string, error) {
	m := vertVisRegex.FindStringSubmatch(input)
	if m == nil {
		return nil, input, errCommandNotHandled
	}
	h, _ := strconv.Atoi(m[1])
	h *= 100
	return &Container{VerticalVisibility: &h}, "", nil
}

func (verticalVisibilityCommand) CanParse(input string) bool {
	return vertVisRegex.MatchString(input)
}

// --- minimalVisibilityCommand ---

var minVisRegex = regexp.MustCompile(`^(?i)(\d{4})(N|E|S|W|NE|NW|SE|SW|NNE|NNW|ENE|ESE|SSE|SSW|WNW|WSW)$`)

type minimalVisibilityCommand struct{}

func (minimalVisibilityCommand) Identifier() string { return "minimalVisibility" }

func (minimalVisibilityCommand) Parse(input string) (*Container, string, error) {
	m := minVisRegex.FindStringSubmatch(input)
	if m == nil {
		return nil, input, errCommandNotHandled
	}
	v, _ := strconv.Atoi(m[1])
	return &Container{
		Visibility: &Visibility{
			Min: &VisibilityMin{
				Value:     v,
				Direction: Direction(m[2]),
			},
		},
	}, "", nil
}

func (minimalVisibilityCommand) CanParse(input string) bool {
	return minVisRegex.MatchString(input)
}

// --- mainVisibilityNauticalMilesCommand ---

var visNMRegex = regexp.MustCompile(`^(P|M)?(\d*)(\s)?((\d\/\d)?SM)$`)

type mainVisibilityNauticalMilesCommand struct{}

func (mainVisibilityNauticalMilesCommand) Identifier() string { return "mainVisibilityNM" }

func (mainVisibilityNauticalMilesCommand) Parse(input string) (*Container, string, error) {
	if !visNMRegex.MatchString(input) {
		return nil, input, errCommandNotHandled
	}
	d, err := convertNauticalMilesVisibility(input)
	if err != nil {
		return nil, input, err
	}
	return &Container{Visibility: &Visibility{Distance: d}}, "", nil
}

func (mainVisibilityNauticalMilesCommand) CanParse(input string) bool {
	return visNMRegex.MatchString(input)
}

// --- commonCommandSupplier ---

type commonCommandSupplier struct {
	commands []commonCommander
}

func newCommonCommandSupplier() *commonCommandSupplier {
	return &commonCommandSupplier{
		commands: []commonCommander{
			windShearCommand{},
			windCommand{},
			windVariationCommand{},
			mainVisibilityCommand{},
			mainVisibilityNauticalMilesCommand{},
			minimalVisibilityCommand{},
			verticalVisibilityCommand{},
			cloudCommand{},
		},
	}
}

func (s *commonCommandSupplier) Get(input string) commonCommander {
	for _, cmd := range s.commands {
		if cmd.CanParse(input) {
			return cmd
		}
	}
	return nil
}

// --- Metar Commands ---

// altimeterCommand: Q1023 -> 1023 hPa.
var altRegex = regexp.MustCompile(`^Q(\d{4})$`)

type altimeterCommand struct{}

func (altimeterCommand) CanParse(s string) bool { return altRegex.MatchString(s) }
func (altimeterCommand) Execute(m *Metar, s string) error {
	parts := altRegex.FindStringSubmatch(s)
	if parts == nil {
		return NewUnexpectedParseError("Match not found")
	}
	v := float64(atoi(parts[1]))
	m.Altimeter = &Altimeter{Value: v, Unit: AltimeterUnitHPa}
	return nil
}

// altimeterMercuryCommand: A2992 -> 29.92 inHg.
var altMercuryRegex = regexp.MustCompile(`^A(\d{4})$`)

type altimeterMercuryCommand struct{}

func (altimeterMercuryCommand) CanParse(s string) bool { return altMercuryRegex.MatchString(s) }
func (altimeterMercuryCommand) Execute(m *Metar, s string) error {
	parts := altMercuryRegex.FindStringSubmatch(s)
	if parts == nil {
		return NewUnexpectedParseError("Match not found")
	}
	v := float64(atoi(parts[1])) / 100
	m.Altimeter = &Altimeter{Value: v, Unit: AltimeterUnitInHg}
	return nil
}

// temperatureCommand: 10/08 -> temp=10, dew=8.
var tempRegex = regexp.MustCompile(`^(M?\d{2})\/(M?\d{2})$`)

type temperatureCommand struct{}

func (temperatureCommand) CanParse(s string) bool { return tempRegex.MatchString(s) }
func (temperatureCommand) Execute(m *Metar, s string) error {
	parts := tempRegex.FindStringSubmatch(s)
	if parts == nil {
		return NewUnexpectedParseError("Match not found")
	}
	t, err := convertTemperature(parts[1])
	if err != nil {
		return err
	}
	d, err := convertTemperature(parts[2])
	if err != nil {
		return err
	}
	m.Temperature = &t
	m.DewPoint = &d
	return nil
}

// runwayCommand: R12/1000, R12/1000V1200, R12/1000U, R12//////.
var runwayGenericRegex = regexp.MustCompile(`^R(\d{2}\w?\/)`)
var runwayDepositRegex = regexp.MustCompile(`^R(\d{2}\w?)\/([/\d])([/\d])(\/\/|\d{2})(\/\/|\d{2})$`)
var runwayRegex = regexp.MustCompile(`^R(\d{2}\w?)\/([MP])?(\d{4})(?:([UDN])|(FT)(?:\/([UDN]))?)$`)
var runwayMaxRangeRegex = regexp.MustCompile(`^R(\d{2}\w?)\/(\d{4})V([MP])?(\d{3,4})(?:([UDN])|(FT)(?:\/([UDN]))?)$`)

type runwayCommand struct{}

func (runwayCommand) CanParse(s string) bool {
	return runwayGenericRegex.MatchString(s)
}

func (runwayCommand) Execute(metar *Metar, s string) error {
	if matches := runwayDepositRegex.FindStringSubmatch(s); matches != nil {
		dt := DepositType(matches[2])
		dc := DepositCoverage(matches[3])
		metar.RunwaysInfo = append(metar.RunwaysInfo, RunwayInfo{
			Deposit: &RunwayInfoDeposit{
				Name:            matches[1],
				DepositType:     &dt,
				Coverage:        &dc,
				Thickness:       matches[4],
				BrakingCapacity: matches[5],
			},
		})
		return nil
	}
	if matches := runwayRegex.FindStringSubmatch(s); matches != nil {
		ri := RunwayInfoRange{
			Name:     matches[1],
			MinRange: atoi(matches[3]),
			Unit:     RunwayInfoUnitMeters,
		}
		if matches[2] != "" {
			vi := ValueIndicator(matches[2])
			ri.Indicator = &vi
		}
		if matches[5] != "" {
			ut := RunwayInfoUnit(matches[5])
			ri.Unit = ut
		}
		t := trendFromMatch(matches[4], matches[6])
		ri.Trend = t
		metar.RunwaysInfo = append(metar.RunwaysInfo, RunwayInfo{Range: &ri})
		return nil
	}
	if matches := runwayMaxRangeRegex.FindStringSubmatch(s); matches != nil {
		ri := RunwayInfoRange{
			Name:     matches[1],
			MinRange: atoi(matches[2]),
			MaxRange: intPtr(atoi(matches[4])),
			Unit:     RunwayInfoUnitMeters,
		}
		if matches[3] != "" {
			vi := ValueIndicator(matches[3])
			ri.Indicator = &vi
		}
		if matches[6] != "" {
			ut := RunwayInfoUnit(matches[6])
			ri.Unit = ut
		}
		t := trendFromMatch(matches[5], matches[7])
		ri.Trend = t
		metar.RunwaysInfo = append(metar.RunwaysInfo, RunwayInfo{Range: &ri})
	}
	return nil
}

func trendFromMatch(trendStr, trendAfterFT string) *RunwayInfoTrend {
	if trendStr != "" {
		t := RunwayInfoTrend(trendStr)
		return &t
	}
	if trendAfterFT != "" {
		t := RunwayInfoTrend(trendAfterFT)
		return &t
	}
	return nil
}

// --- metarCommandSupplier ---

type metarCommandSupplier struct {
	commands []metarCommander
}

func newMetarCommandSupplier() *metarCommandSupplier {
	return &metarCommandSupplier{
		commands: []metarCommander{
			runwayCommand{},
			temperatureCommand{},
			altimeterCommand{},
			altimeterMercuryCommand{},
		},
	}
}

func (s *metarCommandSupplier) Get(input string) metarCommander {
	for _, cmd := range s.commands {
		if cmd.CanParse(input) {
			return cmd
		}
	}
	return nil
}

// --- TAF Commands ---

var icingRegex = regexp.MustCompile(`^6(\d)(\d{3})(\d)$`)
var turbRegex = regexp.MustCompile(`^5(\d|X)(\d{3})(\d)$`)

type icingCommand struct{}
type turbulenceCommand struct{}

func (icingCommand) CanParse(s string) bool      { return icingRegex.MatchString(s) }
func (turbulenceCommand) CanParse(s string) bool { return turbRegex.MatchString(s) }

func (icingCommand) Execute(c *Container, s string) error {
	m := icingRegex.FindStringSubmatch(s)
	if m == nil {
		return NewUnexpectedParseError("Match not found")
	}
	c.Icing = append(c.Icing, Icing{
		Intensity:  IcingIntensity(m[1]),
		BaseHeight: atoi(m[2]) * 100,
		Depth:      atoi(m[3]) * 1000,
	})
	return nil
}

func (turbulenceCommand) Execute(c *Container, s string) error {
	m := turbRegex.FindStringSubmatch(s)
	if m == nil {
		return NewUnexpectedParseError("Match not found")
	}
	c.Turbulence = append(c.Turbulence, Turbulence{
		Intensity:  TurbulenceIntensity(m[1]),
		BaseHeight: atoi(m[2]) * 100,
		Depth:      atoi(m[3]) * 1000,
	})
	return nil
}

// --- tafCommandSupplier ---

type tafCommandSupplier struct {
	commands []tafCommander
}

func newTAFCommandSupplier() *tafCommandSupplier {
	return &tafCommandSupplier{
		commands: []tafCommander{
			turbulenceCommand{},
			icingCommand{},
		},
	}
}

func (s *tafCommandSupplier) Get(input string) tafCommander {
	for _, cmd := range s.commands {
		if cmd.CanParse(input) {
			return cmd
		}
	}
	return nil
}

// --- helpers for commands ---

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
