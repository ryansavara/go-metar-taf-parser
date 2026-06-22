package metartafparser

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// --- Helper functions ---

func assertEqual[T comparable](t *testing.T, got, want T, msg string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", msg, got, want)
	}
}
func assertNil[T any](t *testing.T, ptr *T, msg string) {
	t.Helper()
	if ptr != nil {
		t.Errorf("%s: expected nil, got %v", msg, *ptr)
	}
}
func assertNotNil[T any](t *testing.T, ptr *T, msg string) {
	t.Helper()
	if ptr == nil {
		t.Errorf("%s: expected non-nil", msg)
	}
}
func assertIntPtr(t *testing.T, got *int, want int, msg string) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %d", msg, want)
	} else if *got != want {
		t.Errorf("%s: got %d, want %d", msg, *got, want)
	}
}
func assertFloatPtr(t *testing.T, got *float64, want float64, msg string) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %f", msg, want)
	} else if *got != want {
		t.Errorf("%s: got %f, want %f", msg, *got, want)
	}
}
func assertBoolPtr(t *testing.T, got *bool, want bool, msg string) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %v", msg, want)
	} else if *got != want {
		t.Errorf("%s: got %v, want %v", msg, *got, want)
	}
}

// ============================================================
// parseWeatherCondition tests
// ============================================================

func TestParseWeatherConditionLightDrizzle(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("-DZ")
	if wc == nil {
		t.Fatal("expected non-nil weather condition")
	}
	assertNotNil(t, wc.Intensity, "intensity")
	assertEqual(t, *wc.Intensity, IntensityLight, "intensity")
	assertEqual(t, len(wc.Phenomena), 1, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonDrizzle, "phenomenon[0]")
}

func TestParseWeatherConditionLightSnowRainOrder(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("-SNRA")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertEqual(t, len(wc.Phenomena), 2, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonSnow, "phenomenon[0]")
	assertEqual(t, wc.Phenomena[1], PhenomenonRain, "phenomenon[1]")
}

func TestParseWeatherConditionShowersRainHail(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("SHRAGR")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertEqual(t, *wc.Intensity, IntensityModerate, "intensity")
	assertNotNil(t, wc.Descriptive, "descriptive")
	assertEqual(t, *wc.Descriptive, DescriptiveShowers, "descriptive")
	assertEqual(t, len(wc.Phenomena), 2, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonRain, "phenomenon[0]")
	assertEqual(t, wc.Phenomena[1], PhenomenonHail, "phenomenon[1]")
}

func TestParseWeatherConditionDrizzleMist(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("-DZ/BR")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertNotNil(t, wc.Intensity, "intensity")
	assertEqual(t, *wc.Intensity, IntensityLight, "intensity")
	assertNil(t, wc.Descriptive, "descriptive")
	assertEqual(t, len(wc.Phenomena), 2, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonDrizzle, "phenomenon[0]")
	assertEqual(t, wc.Phenomena[1], PhenomenonMist, "phenomenon[1]")
}

func TestParseWeatherConditionMistSlash(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("BR/")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertEqual(t, len(wc.Phenomena), 1, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonMist, "phenomenon[0]")
}

func TestParseWeatherConditionInvalid(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("BRd")
	if wc != nil {
		t.Fatal("expected nil for invalid input")
	}
}

func TestParseWeatherConditionHeavyFunnelCloud(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("+FC")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertNil(t, wc.Intensity, "intensity")
	assertEqual(t, len(wc.Phenomena), 1, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonTornado, "phenomenon[0]")
}

func TestParseWeatherConditionFunnelCloud(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("FC")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertEqual(t, *wc.Intensity, IntensityModerate, "intensity")
	assertEqual(t, len(wc.Phenomena), 1, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonFunnelCloud, "phenomenon[0]")
}

func TestParseWeatherConditionVicinityBlowingSnow(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("VCBLSN")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertNotNil(t, wc.Intensity, "intensity")
	assertEqual(t, *wc.Intensity, IntensityInVicinity, "intensity")
	assertNotNil(t, wc.Descriptive, "descriptive")
	assertEqual(t, *wc.Descriptive, DescriptiveBlowing, "descriptive")
	assertEqual(t, len(wc.Phenomena), 1, "phenomenons length")
	assertEqual(t, wc.Phenomena[0], PhenomenonSnow, "phenomenon[0]")
}

// ============================================================
// Tokenize tests
// ============================================================

func TestTokenize(t *testing.T) {
	input := "METAR KTTN 051853Z 04011KT 1 1/2SM VCTS SN FZFG BKN003 OVC010 M02/M02 A3006 RMK AO2 TSB40 SLP176 P0002 T10171017="
	tokens := tokenize(input)
	want := []string{"METAR", "KTTN", "051853Z", "04011KT", "1 1/2SM", "VCTS", "SN", "FZFG", "BKN003", "OVC010", "M02/M02", "A3006", "RMK", "AO2", "TSB40", "SLP176", "P0002", "T10171017"}
	if len(tokens) != len(want) {
		t.Fatalf("token count: got %d, want %d\n  got:  %q\n  want: %q", len(tokens), len(want), tokens, want)
	}
	for i := range want {
		if tokens[i] != want[i] {
			t.Errorf("token[%d]: got %q, want %q", i, tokens[i], want[i])
		}
	}
}

// ============================================================
// generalParse tests (data-driven)
// ============================================================

func testGeneralParse(t *testing.T, input string, expected bool) {
	t.Helper()
	p := newAbstractParser(DefaultLocale())
	// Pre-initialize Wind and Visibility to satisfy dependency checks
	c := &Container{
		Wind:       &Wind{},
		Visibility: &Visibility{},
	}
	got := p.generalParse(c, input)
	if got != expected {
		t.Errorf("generalParse(%q): got %v, want %v", input, got, expected)
	}
}

func TestGeneralParseWind(t *testing.T)           { testGeneralParse(t, "05009KT", true) }
func TestGeneralParseWindVar(t *testing.T)        { testGeneralParse(t, "030V113", true) }
func TestGeneralParseVisibility9999(t *testing.T) { testGeneralParse(t, "9999", true) }
func TestGeneralParseFractionalVis(t *testing.T)  { testGeneralParse(t, "6 1/2SM", true) }
func TestGeneralParseMinVis(t *testing.T)         { testGeneralParse(t, "1100w", true) }
func TestGeneralParseVerticalVis(t *testing.T)    { testGeneralParse(t, "VV002", true) }
func TestGeneralParseCavok(t *testing.T)          { testGeneralParse(t, "CAVOK", true) }
func TestGeneralParseCloudCB(t *testing.T)        { testGeneralParse(t, "SCT026CB", true) }
func TestGeneralParseInvalidCloud(t *testing.T)   { testGeneralParse(t, "ZZZ026CV", false) }
func TestGeneralParseWeatherShowers(t *testing.T) { testGeneralParse(t, "+SHGSRA", true) }
func TestGeneralParseWeatherHeavy(t *testing.T) {
	p := newAbstractParser(DefaultLocale())
	wc := p.parseWeatherCondition("+SHGSRA")
	if wc == nil {
		t.Fatal("expected non-nil")
	}
	assertEqual(t, *wc.Intensity, IntensityHeavy, "intensity")
}
func TestGeneralParseInvalidWeather(t *testing.T) { testGeneralParse(t, "+VFDR", false) }

// ============================================================
// MetarParser tests
// ============================================================

func TestMetarParserBasic(t *testing.T) {
	input := "LFPG 170830Z 00000KT 0350 R27L/0375N R09R/0175N R26R/0500D R08L/0400N R26L/0275D R08R/0250N R27R/0300N R09L/0200N FG SCT000 M01/M01 Q1026 NOSIG"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, result.Station, "LFPG", "station")
	assertIntPtr(t, result.Day, 17, "day")
	assertIntPtr(t, result.Hour, 8, "hour")
	assertIntPtr(t, result.Minute, 30, "minute")
	assertNotNil(t, result.Wind, "wind")
	assertEqual(t, result.Wind.Speed, 0, "wind.speed")
	assertEqual(t, result.Wind.Direction, "N", "wind.direction")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(350), "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitMeters, "vis.unit")
	assertEqual(t, len(result.RunwaysInfo), 8, "runwaysInfo length")
	if len(result.RunwaysInfo) > 0 {
		assertEqual(t, result.RunwaysInfo[0].Range.Name, "27L", "runway[0].name")
		if result.RunwaysInfo[0].Range != nil {
			assertEqual(t, result.RunwaysInfo[0].Range.MinRange, 375, "runway[0].minRange")
			if result.RunwaysInfo[0].Range.Trend != nil {
				assertEqual(t, string(*result.RunwaysInfo[0].Range.Trend), "N", "runway[0].trend")
			}
		}
	}
}

func TestMetarParserCanadian(t *testing.T) {
	input := "CYWG 172000Z 30015G25KT 3/4SM R36/4000FT/D -SN BLSN BKN008 OVC040 M05/M08 A2992 REFZRA WS RWY36 RMK SF5NS3 SLP134"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, result.Station, "CYWG", "station")
	assertIntPtr(t, result.Day, 17, "day")
	assertIntPtr(t, result.Hour, 20, "hour")
	assertIntPtr(t, result.Minute, 0, "minute")
	assertNotNil(t, result.Wind, "wind")
	assertEqual(t, result.Wind.Speed, 15, "wind.speed")
	assertIntPtr(t, result.Wind.Gust, 25, "wind.gust")
	assertEqual(t, result.Wind.Direction, "WNW", "wind.direction")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(0.75), "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitStatuteMiles, "vis.unit")
	assertEqual(t, len(result.RunwaysInfo), 1, "runwaysInfo length")
	if len(result.RunwaysInfo) > 0 {
		assertEqual(t, result.RunwaysInfo[0].Range.Name, "36", "runway[0].name")
		if result.RunwaysInfo[0].Range != nil {
			assertEqual(t, result.RunwaysInfo[0].Range.MinRange, 4000, "runway[0].minRange")
			assertEqual(t, result.RunwaysInfo[0].Range.Unit, RunwayInfoUnitFeet, "runway[0].unit")
			assertNotNil(t, result.RunwaysInfo[0].Range.Trend, "runway[0].trend")
			assertEqual(t, *result.RunwaysInfo[0].Range.Trend, RunwayInfoTrendDecreasing, "runway[0].trend")
		}
	}
}

func TestMetarParserAutoAsStation(t *testing.T) {
	input := "AUTO 061950Z 10002KT 9999NDV NCD 01/M00 Q1015 RMK="
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNil(t, result.Auto, "auto") // "AUTO" is station here, not flag
	assertEqual(t, result.Station, "AUTO", "station")
	assertIntPtr(t, result.Day, 6, "day")
	assertIntPtr(t, result.Hour, 19, "hour")
	assertIntPtr(t, result.Minute, 50, "minute")
}

func TestMetarParserAutoFlag(t *testing.T) {
	input := "AUTO LSZL 061950Z 10002KT 9999NDV NCD 01/M00 Q1015 RMK="
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertBoolPtr(t, result.Auto, true, "auto")
	assertEqual(t, result.Station, "LSZL", "station")
	assertIntPtr(t, result.Day, 6, "day")
	assertIntPtr(t, result.Hour, 19, "hour")
	assertIntPtr(t, result.Minute, 50, "minute")
}

func TestMetarParserTempo(t *testing.T) {
	input := "LFBG 081130Z AUTO 23012KT 9999 SCT022 BKN072 BKN090 22/16 Q1011 TEMPO 26015G25KT 3000 TSRA SCT025CB BKN050"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertBoolPtr(t, result.Auto, true, "auto")
	assertEqual(t, len(result.Clouds), 3, "clouds length")
	assertEqual(t, len(result.Trends), 1, "trends length")
	trend := result.Trends[0]
	assertEqual(t, trend.Type, WeatherChangeTypeTEMPO, "trend.type")
	assertNotNil(t, trend.Wind, "trend.wind")
	assertIntPtr(t, trend.Wind.Degrees, 260, "trend.wind.degrees")
	assertEqual(t, trend.Wind.Speed, 15, "trend.wind.speed")
	assertIntPtr(t, trend.Wind.Gust, 25, "trend.wind.gust")
	assertNotNil(t, trend.Visibility, "trend.visibility")
	assertEqual(t, trend.Visibility.Value, float64(3000), "trend.vis.value")
	assertEqual(t, trend.Visibility.Unit, DistanceUnitMeters, "trend.vis.unit")
	assertEqual(t, len(trend.WeatherConditions), 1, "trend.wc length")
	assertEqual(t, trend.WeatherConditions[0].Phenomena[0], PhenomenonRain, "trend.wc phenomenon")
	assertEqual(t, len(trend.Clouds), 2, "trend.clouds length")
	if len(trend.Clouds) >= 2 {
		assertEqual(t, trend.Clouds[0].Quantity, CloudQuantitySCT, "trend.clouds[0].quantity")
		assertIntPtr(t, trend.Clouds[0].Height, 2500, "trend.clouds[0].height")
		if trend.Clouds[0].Type != nil {
			assertEqual(t, *trend.Clouds[0].Type, CloudTypeCB, "trend.clouds[0].type")
		} else {
			t.Error("expected trend.clouds[0].type to be CB")
		}
		assertEqual(t, trend.Clouds[1].Quantity, CloudQuantityBKN, "trend.clouds[1].quantity")
		assertIntPtr(t, trend.Clouds[1].Height, 5000, "trend.clouds[1].height")
		assertNil(t, trend.Clouds[1].Type, "trend.clouds[1].type")
	}
	assertEqual(t, len(trend.Times), 0, "trend.times length")
	assertEqual(t, trend.Raw, "TEMPO 26015G25KT 3000 TSRA SCT025CB BKN050", "trend.raw")
}

func TestMetarParserTempoBecmg(t *testing.T) {
	input := "LFRM 081630Z AUTO 30007KT 260V360 9999 24/15 Q1008 TEMPO SHRA BECMG SKC"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.Trends), 2, "trends length")
	assertEqual(t, result.Trends[0].Type, WeatherChangeTypeTEMPO, "trend[0].type")
	assertEqual(t, len(result.Trends[0].WeatherConditions), 1, "trend[0].wc length")
	if len(result.Trends[0].WeatherConditions) > 0 {
		assertNotNil(t, result.Trends[0].WeatherConditions[0].Descriptive, "trend[0].wc.descriptive")
		assertEqual(t, *result.Trends[0].WeatherConditions[0].Descriptive, DescriptiveShowers, "trend[0].wc.descriptive")
		assertEqual(t, result.Trends[0].WeatherConditions[0].Phenomena[0], PhenomenonRain, "trend[0].wc.phenomenon")
	}
	assertEqual(t, result.Trends[0].Raw, "TEMPO SHRA", "trend[0].raw")
	assertEqual(t, result.Trends[1].Type, WeatherChangeTypeBECMG, "trend[1].type")
	assertEqual(t, result.Trends[1].Raw, "BECMG SKC", "trend[1].raw")
	assertEqual(t, len(result.Trends[1].Clouds), 1, "trend[1].clouds length")
}

func TestMetarParserTempoFM(t *testing.T) {
	input := "LFRM 081630Z AUTO 30007KT 260V360 9999 24/15 Q1008 TEMPO FM1830 SHRA"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.Trends), 1, "trends length")
	trend := result.Trends[0]
	assertEqual(t, trend.Type, WeatherChangeTypeTEMPO, "trend.type")
	assertEqual(t, len(trend.WeatherConditions), 1, "trend.wc length")
	if len(trend.Times) > 0 {
		assertEqual(t, trend.Times[0].Type, TimeIndicatorFM, "trend.time[0].type")
		assertIntPtr(t, trend.Times[0].Hour, 18, "trend.time[0].hour")
		assertIntPtr(t, trend.Times[0].Minute, 30, "trend.time[0].minute")
	}
}

func TestMetarParserTempoTL(t *testing.T) {
	input := "LFRM 081630Z AUTO 30007KT 260V360 9999 24/15 Q1008 TEMPO FM1700 TL1830 SHRA"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.Trends), 1, "trends length")
	trend := result.Trends[0]
	if len(trend.Times) >= 2 {
		assertEqual(t, trend.Times[0].Type, TimeIndicatorFM, "trend.time[0].type")
		assertIntPtr(t, trend.Times[0].Hour, 17, "trend.time[0].hour")
		assertEqual(t, trend.Times[1].Type, TimeIndicatorTL, "trend.time[1].type")
		assertIntPtr(t, trend.Times[1].Hour, 18, "trend.time[1].hour")
		assertIntPtr(t, trend.Times[1].Minute, 30, "trend.time[1].minute")
	}
}

func TestMetarParserMinVisibility(t *testing.T) {
	input := "LFPG 161430Z 24015G25KT 5000 1100w"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertIntPtr(t, result.Day, 16, "day")
	assertIntPtr(t, result.Hour, 14, "hour")
	assertIntPtr(t, result.Minute, 30, "minute")
	assertNotNil(t, result.Wind, "wind")
	assertIntPtr(t, result.Wind.Degrees, 240, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 15, "wind.speed")
	assertIntPtr(t, result.Wind.Gust, 25, "wind.gust")
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(5000), "vis.value")
	if result.Visibility.Min != nil {
		assertEqual(t, result.Visibility.Min.Value, 1100, "vis.min.value")
		assertEqual(t, result.Visibility.Min.Direction, "w", "vis.min.direction")
	} else {
		t.Error("expected visibility.min to be set")
	}
}

func TestMetarParserWind0000KT(t *testing.T) {
	input := "KATW 022045Z 00000KT 10SM SCT120 00/M08 A2996"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Wind, "wind")
	assertIntPtr(t, result.Wind.Degrees, 0, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 0, "wind.speed")
	assertEqual(t, result.Wind.Direction, "N", "wind.direction")
	if result.Visibility != nil && result.Visibility.Min != nil {
		t.Error("expected visibility.min to be nil")
	}
}

func TestMetarParserWindVariation(t *testing.T) {
	input := "LFPG 161430Z 24015G25KT 180V300"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Wind, "wind")
	assertIntPtr(t, result.Wind.Degrees, 240, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 15, "wind.speed")
	assertIntPtr(t, result.Wind.MinVariation, 180, "wind.minVariation")
	assertIntPtr(t, result.Wind.MaxVariation, 300, "wind.maxVariation")
}

func TestMetarParserVerticalVisibility(t *testing.T) {
	input := "LFLL 160730Z 28002KT 0350 FG VV002"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertIntPtr(t, result.Day, 16, "day")
	assertIntPtr(t, result.Hour, 7, "hour")
	assertIntPtr(t, result.Minute, 30, "minute")
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(350), "vis.value")
	assertIntPtr(t, result.VerticalVisibility, 200, "verticalVisibility")
	assertEqual(t, len(result.WeatherConditions), 1, "wc length")
	assertEqual(t, result.WeatherConditions[0].Phenomena[0], PhenomenonFog, "wc phenomenon")
}

func TestMetarParserNDV(t *testing.T) {
	input := "LSZL 300320Z AUTO 00000KT 9999NDV BKN060 OVC074 00/M04 Q1001 RMK="
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(9999), "vis.value")
	assertEqual(t, result.Visibility.Ndv, true, "vis.ndv")
}

func TestMetarParserCavok(t *testing.T) {
	input := "LFPG 212030Z 03003KT CAVOK 09/06 Q1031 NOSIG"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertBoolPtr(t, result.Cavok, true, "cavok")
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(9999), "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitMeters, "vis.unit")
	assertNotNil(t, result.Temperature, "temperature")
	assertEqual(t, *result.Temperature, float64(9), "temperature")
	assertNotNil(t, result.DewPoint, "dewPoint")
	assertEqual(t, *result.DewPoint, float64(6), "dewPoint")
	assertNotNil(t, result.Altimeter, "altimeter")
	assertEqual(t, result.Altimeter.Value, float64(1031), "altimeter.value")
	assertEqual(t, result.Altimeter.Unit, AltimeterUnitHPa, "altimeter.unit")
	assertBoolPtr(t, result.Nosig, true, "nosig")
}

func TestMetarParserAltimeterMercury(t *testing.T) {
	input := "KTTN 051853Z 04011KT 9999 VCTS SN FZFG BKN003 OVC010 M02/M02 A3006"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Altimeter, "altimeter")
	assertEqual(t, result.Altimeter.Value, float64(30.06), "altimeter.value")
	assertEqual(t, result.Altimeter.Unit, AltimeterUnitInHg, "altimeter.unit")
	assertEqual(t, len(result.WeatherConditions), 3, "wc length")
}

func TestMetarParserDescriptiveOnly(t *testing.T) {
	input := "AGGH 140340Z 05010KT 9999 TS FEW020 SCT021CB BKN300 32/26 Q1010"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.WeatherConditions), 1, "wc length")
	assertNotNil(t, result.WeatherConditions[0].Descriptive, "wc.descriptive")
	assertEqual(t, *result.WeatherConditions[0].Descriptive, DescriptiveThunderstorm, "wc.descriptive")
}

func TestMetarParserInvalidWeather(t *testing.T) {
	input := "ENLK 081350Z 26026G40 240V300 9999 VCMI"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.WeatherConditions), 0, "wc length")
}

func TestMetarParserRunwayDeposit(t *testing.T) {
	input := "UNAA 240830Z 34002MPS CAVOK M14/M18 Q1019 R02/190054 NOSIG RMK QFE741"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, result.Station, "UNAA", "station")
	assertNotNil(t, result.Wind, "wind")
	assertIntPtr(t, result.Wind.Degrees, 340, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 2, "wind.speed")
	assertEqual(t, result.Wind.Unit, SpeedUnitMetersPerSecond, "wind.unit")
	assertBoolPtr(t, result.Cavok, true, "cavok")
	assertBoolPtr(t, result.Nosig, true, "nosig")
	assertNotNil(t, result.Remark, "remark")
	if result.Remark != nil {
		if !strings.Contains(*result.Remark, "QFE741") {
			t.Errorf("remark should contain QFE741, got %q", *result.Remark)
		}
	}
	assertEqual(t, len(result.Remarks), 1, "remarks length")
}

func TestMetarParserTempoRmk(t *testing.T) {
	input := "EGLL 231250Z 14012G22KT 2000 +TSRA FG BKN008 SCT025CB OVC050 18/17 Q1010 TEMPO 1000 -SHRA RMK QFE998"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.Trends), 1, "trends length")
	assertEqual(t, result.Trends[0].Type, WeatherChangeTypeTEMPO, "trend[0].type")
	assertNotNil(t, result.Remark, "remark")
	if result.Remark != nil {
		if !strings.Contains(*result.Remark, "QFE998") {
			t.Errorf("remark should contain QFE998, got %q", *result.Remark)
		}
	}
}

func TestMetarParserNil(t *testing.T) {
	input := "SVMC 211703Z AUTO NIL"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertBoolPtr(t, result.Nil, true, "nil")
}

func TestMetarParserLessThanQuarterVis(t *testing.T) {
	input := "SUMU 070520Z M1/4SM"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Visibility, "visibility")
	assertNotNil(t, result.Visibility.Indicator, "vis.indicator")
	assertEqual(t, *result.Visibility.Indicator, ValueIndicatorLessThan, "vis.indicator")
	assertEqual(t, result.Visibility.Value, float64(0.25), "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitStatuteMiles, "vis.unit")
}

func TestMetarParserP6SM(t *testing.T) {
	input := "SUMU 070520Z P6SM"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Visibility, "visibility")
	assertNotNil(t, result.Visibility.Indicator, "vis.indicator")
	assertEqual(t, *result.Visibility.Indicator, ValueIndicatorGreaterThan, "vis.indicator")
	assertEqual(t, result.Visibility.Value, float64(6), "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitStatuteMiles, "vis.unit")
}

func TestMetarParser3QuarterVis(t *testing.T) {
	input := "SUMU 070520Z 3 1/4SM"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(3.25), "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitStatuteMiles, "vis.unit")
}

func TestMetarParserMetarType(t *testing.T) {
	tests := []struct {
		input string
		mtype *MetarType
	}{
		{"SUMU 070520Z 3 1/4SM", nil},
		{"METAR SUMU 070520Z 3 1/4SM", metarTypePtr(MetarTypeMETAR)},
		{"SPECI SUMU 070520Z 3 1/4SM", metarTypePtr(MetarTypeSPECI)},
	}
	for _, tt := range tests {
		result, err := ParseMetar(tt.input, nil)
		if err != nil {
			t.Fatalf("ParseMetar(%q) error: %v", tt.input, err)
		}
		if tt.mtype == nil {
			if result.Type != nil {
				t.Errorf("ParseMetar(%q): expected nil type, got %v", tt.input, *result.Type)
			}
		} else {
			if result.Type == nil {
				t.Errorf("ParseMetar(%q): expected type %v, got nil", tt.input, *tt.mtype)
			} else if *result.Type != *tt.mtype {
				t.Errorf("ParseMetar(%q): got %v, want %v", tt.input, *result.Type, *tt.mtype)
			}
		}
	}
}

func metarTypePtr(m MetarType) *MetarType { return &m }

// ---- Missing MetarParser tests ----

func TestMetarParserWind00000MPS(t *testing.T) {
	input := "KATL 270200Z 00000MPS"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Wind, "wind")
	assertEqual(t, *result.Wind.Degrees, 0, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 0, "wind.speed")
	assertEqual(t, result.Wind.Unit, SpeedUnitMetersPerSecond, "wind.unit")
	assertEqual(t, result.Wind.Direction, "N", "wind.direction")
}

func TestMetarParserVisNotWind(t *testing.T) {
	input := "VIDP 270200Z 00000MPS 0050"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Wind, "wind")
	assertEqual(t, *result.Wind.Degrees, 0, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 0, "wind.speed")
	assertEqual(t, result.Wind.Unit, SpeedUnitMetersPerSecond, "wind.unit")
	assertNotNil(t, result.Visibility, "vis")
	assertEqual(t, result.Visibility.Value, 50.0, "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitMeters, "vis.unit")
}

func TestMetarParserWindAltForm(t *testing.T) {
	input := "ENLK 081350Z 26026G40 240V300 9999 VCSH FEW025 BKN030 02/M01 Q0996"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Wind, "wind")
	assertEqual(t, *result.Wind.Degrees, 260, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 26, "wind.speed")
	assertIntPtr(t, result.Wind.Gust, 40, "wind.gust")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
	assertIntPtr(t, result.Wind.MinVariation, 240, "wind.minVar")
	assertIntPtr(t, result.Wind.MaxVariation, 300, "wind.maxVar")
	assertEqual(t, len(result.WeatherConditions), 1, "wc length")
	assertNotNil(t, result.WeatherConditions[0].Intensity, "wc[0].intensity")
	assertEqual(t, *result.WeatherConditions[0].Intensity, IntensityInVicinity, "wc[0].intensity")
	if result.WeatherConditions[0].Descriptive != nil {
		assertEqual(t, *result.WeatherConditions[0].Descriptive, DescriptiveShowers, "wc[0].descriptive")
	}
}

func TestMetarParserNoneDeposit(t *testing.T) {
	input := "UUWW 151030Z 34002MPS CAVOK 14/02 Q1026 R01/000070 NOSIG"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.RunwaysInfo), 1, "runwaysInfo length")
	assertEqual(t, result.RunwaysInfo[0].Deposit.Name, "01", "runway.name")
}

func TestMetarParserMinVisSUMU(t *testing.T) {
	input := "SUMU 070520Z 34025KT 8000 2000SW VCSH SCT013CB BKN026 00/M05 Q1012 TEMPO 2000 SHSN="
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Visibility, "vis")
	assertNotNil(t, result.Visibility.Min, "vis.min")
	assertEqual(t, result.Visibility.Min.Value, 2000, "vis.min.value")
	assertEqual(t, result.Visibility.Min.Direction, "SW", "vis.min.direction")
}

func TestMetarParserMoreThan1HalfVis(t *testing.T) {
	input := "SUMU 070520Z P1 1/2SM"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertNotNil(t, result.Visibility, "vis")
	assertNotNil(t, result.Visibility.Indicator, "vis.indicator")
	assertEqual(t, *result.Visibility.Indicator, ValueIndicatorGreaterThan, "vis.indicator")
	assertEqual(t, result.Visibility.Value, 1.5, "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitStatuteMiles, "vis.unit")
}

func TestMetarParserUnknownCloudTypes(t *testing.T) {
	input := "EKVG 291550Z AUTO 13009KT 9999 BKN037/// BKN048/// 07/06 Q1009 RMK FEW011/// FEW035/// WIND SKEID 13020KT"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.Clouds), 2, "clouds length")
}

func TestMetarParserInvalidCloudQty(t *testing.T) {
	input := "EKVG 291550Z AUTO 13009KT 9999 BKN037AAA"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.Clouds), 0, "clouds length")
}

func TestMetarParserInvalidCloudType(t *testing.T) {
	input := "EKVG 291550Z AUTO 13009KT 9999 AAA037"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.Clouds), 0, "clouds length")
}

func TestMetarParserThreeWeather(t *testing.T) {
	input := "CYVM 282100Z 36028G36KT 1SM -SN DRSN VCBLSN OVC008 M03/M04 A2935 RMK SN2ST8 LAST STFFD OBS/NXT 291200UTC SLP940"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.WeatherConditions), 3, "wc length")
	assertNotNil(t, result.WeatherConditions[0].Intensity, "wc[0].intensity")
	assertEqual(t, *result.WeatherConditions[0].Intensity, IntensityLight, "wc[0].intensity")
	if result.WeatherConditions[0].Descriptive != nil {
		t.Errorf("wc[0].descriptive should be nil, got %v", *result.WeatherConditions[0].Descriptive)
	}
	assertEqual(t, result.WeatherConditions[0].Phenomena[0], PhenomenonSnow, "wc[0].phenomenon")
	assertNotNil(t, result.WeatherConditions[1].Descriptive, "wc[1].descriptive")
	assertEqual(t, *result.WeatherConditions[1].Descriptive, DescriptiveDrifting, "wc[1].descriptive")
	assertEqual(t, result.WeatherConditions[1].Phenomena[0], PhenomenonSnow, "wc[1].phenomenon")
	assertNotNil(t, result.WeatherConditions[2].Intensity, "wc[2].intensity")
	assertEqual(t, *result.WeatherConditions[2].Intensity, IntensityInVicinity, "wc[2].intensity")
	if result.WeatherConditions[2].Descriptive != nil {
		assertEqual(t, *result.WeatherConditions[2].Descriptive, DescriptiveBlowing, "wc[2].descriptive")
	}
	assertEqual(t, result.WeatherConditions[2].Phenomena[0], PhenomenonSnow, "wc[2].phenomenon")
}

func ExampleParseMetar() {
	metar, err := ParseMetar("KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992 RMK AO2", nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(metar.Station)
	// Output: KLAX
}

func ExampleParseTAF() {
	taf, err := ParseTAF("KLAX 140520Z 1406/1512 05005KT P6SM FEW010", nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(taf.Station)
	// Output: KLAX
}

func ExampleParseMetarDated() {
	issued := time.Date(2024, 6, 14, 8, 53, 0, 0, time.UTC)
	metar, err := ParseMetarDated("KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992 RMK AO2", issued, nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(metar.Station)
	// Output: KLAX
}

func ExampleParseTAFAsForecast() {
	issued := time.Date(2024, 6, 14, 5, 20, 0, 0, time.UTC)
	fc, err := ParseTAFAsForecast("KLAX 140520Z 1406/1512 05005KT P6SM FEW010", issued, nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(fc.Station)
	// Output: KLAX
}

func ExampleParseTAFDated() {
	issued := time.Date(2024, 6, 14, 5, 20, 0, 0, time.UTC)
	taf, err := ParseTAFDated("KLAX 140520Z 1406/1512 05005KT P6SM FEW010", issued, nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(taf.Station)
	// Output: KLAX
}

func ExampleGetCompositeForecastForDate() {
	issued := time.Date(2024, 6, 14, 5, 20, 0, 0, time.UTC)
	fc, err := ParseTAFAsForecast("KLAX 140520Z 1406/1512 05005KT P6SM FEW010", issued, nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	date := time.Date(2024, 6, 14, 8, 0, 0, 0, time.UTC)
	cf, err := GetCompositeForecastForDate(date, fc)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(cf.Prevailing.Wind.Speed)
	// Output: 5
}

// ============================================================
// Error type tests
// ============================================================

func TestParseError(t *testing.T) {
	err := NewParseError("test error")
	assertEqual(t, err.Error(), "test error", "Error()")
}

func TestInvalidWeatherStatementError(t *testing.T) {
	err1 := NewInvalidWeatherStatementError(nil)
	assertEqual(t, err1.Error(), "Invalid weather string", "Error() with nil")

	err2 := NewInvalidWeatherStatementError("bad input")
	assertEqual(t, err2.Error(), "Invalid weather string: bad input", "Error() with string")

	err3 := NewInvalidWeatherStatementError(42)
	assertEqual(t, err3.Error(), "Invalid weather string", "Error() with non-string")
}

func TestInvalidWeatherStatementErrorUnwrap(t *testing.T) {
	err := NewInvalidWeatherStatementError("test")
	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Error("expected errors.As to find *ParseError")
	}
}

func TestUnexpectedParseError(t *testing.T) {
	err := NewUnexpectedParseError("unexpected")
	assertEqual(t, err.Error(), "unexpected", "Error()")
	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Error("expected errors.As to find *ParseError")
	}
}

func TestTimestampOutOfBoundsError(t *testing.T) {
	err := NewTimestampOutOfBoundsError("bounds")
	assertEqual(t, err.Error(), "bounds", "Error()")
}

func TestPartialWeatherStatementErrorUnwrap(t *testing.T) {
	err := NewPartialWeatherStatementError("partial", 1, 3)
	var invalidErr *InvalidWeatherStatementError
	if !errors.As(err, &invalidErr) {
		t.Error("expected errors.As to find *InvalidWeatherStatementError")
	}
	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Error("expected errors.As to find *ParseError")
	}
}

// ============================================================
// Runway max range variant test
// ============================================================

func TestMetarParserRunwayMaxRange(t *testing.T) {
	input := "KATL 022045Z 00000KT 10SM SCT120 00/M08 A2996 R12/1000V1200U"
	result, err := ParseMetar(input, nil)
	if err != nil {
		t.Fatalf("ParseMetar(%q) error: %v", input, err)
	}
	assertEqual(t, len(result.RunwaysInfo), 1, "runwaysInfo length")
	if len(result.RunwaysInfo) > 0 {
		ri := result.RunwaysInfo[0]
		if ri.Range == nil {
			t.Fatal("expected Range to be non-nil")
		}
		assertEqual(t, ri.Range.Name, "12", "range.name")
		assertEqual(t, ri.Range.MinRange, 1000, "range.minRange")
		assertIntPtr(t, ri.Range.MaxRange, 1200, "range.maxRange")
		if ri.Range.Trend != nil {
			assertEqual(t, string(*ri.Range.Trend), "U", "range.trend")
		}
	}
}

// ============================================================
// isValidPhenomenon coverage
// ============================================================

func TestIsValidPhenomenonAll(t *testing.T) {
	all := []string{"RA", "DZ", "SN", "SG", "PL", "IC", "GR", "GS", "UP",
		"FG", "VA", "BR", "HZ", "DU", "FU", "SA", "PY", "SQ", "PO",
		"TS", "DS", "SS", "FC", "NSW"}
	for _, p := range all {
		if !isValidPhenomenon(p) {
			t.Errorf("isValidPhenomenon(%q) should be true", p)
		}
	}
	if isValidPhenomenon("XX") {
		t.Error("isValidPhenomenon('XX') should be false")
	}
}

// ============================================================
// Locale edge case tests
// ============================================================

func TestLocaleGetNonStringValue(t *testing.T) {
	locale := DefaultLocale()
	r := localeGet("Converter.D", locale)
	if r == nil {
		t.Fatalf("expected non-nil for valid string key")
	}
	assertEqual(t, *r, "decreasing", "Converter.D")

	r = localeGet("Barometer", locale)
	if r != nil {
		t.Error("expected nil for non-string value (slice)")
	}
}

func TestLocaleGetMissingKey(t *testing.T) {
	locale := DefaultLocale()
	r := localeGet("NonExistent.Key", locale)
	if r != nil {
		t.Errorf("expected nil for missing key, got %q", *r)
	}
}

func TestFormatMsgNilMessage(t *testing.T) {
	r := formatMsg(nil, "arg1")
	if r != nil {
		t.Errorf("expected nil for nil message, got %q", *r)
	}
}

func TestFormatMsgNilArg(t *testing.T) {
	msg := "test {0}"
	r := formatMsg(&msg, nil)
	if r != nil {
		t.Errorf("expected nil for nil arg, got %q", *r)
	}
}

func TestFormatMsgBasic(t *testing.T) {
	msg := "hello {0} world {1}"
	r := formatMsg(&msg, "foo", 42)
	if r == nil {
		t.Fatal("expected non-nil")
	}
	assertEqual(t, *r, "hello foo world 42", "formatted")
}

func TestFormatMsgOutOfRange(t *testing.T) {
	msg := "hello {0} world {5}"
	r := formatMsg(&msg, "foo")
	if r == nil {
		t.Fatal("expected non-nil")
	}
	assertEqual(t, *r, "hello foo world {5}", "unmatched index preserved")
}
