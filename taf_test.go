package metartafparser

import (
	"errors"
	"testing"
)

// ============================================================
// parseValidity tests
// ============================================================

func TestParseValidity(t *testing.T) {
	v := parseValidity("3118/0124")
	assertEqual(t, v.StartDay, 31, "startDay")
	assertEqual(t, v.StartHour, 18, "startHour")
	assertEqual(t, v.EndDay, 1, "endDay")
	assertEqual(t, v.EndHour, 24, "endHour")
}

// ============================================================
// parseTemperature tests
// ============================================================

func TestParseTemperatureMax(t *testing.T) {
	tm, err := parseTemperature("TX15/0612Z")
	if err != nil {
		t.Fatalf("parseTemperature error: %v", err)
	}
	assertEqual(t, tm.Temperature, float64(15), "temperature")
	assertEqual(t, tm.Day, 6, "day")
	assertEqual(t, tm.Hour, 12, "hour")
}

func TestParseTemperatureMin(t *testing.T) {
	tm, err := parseTemperature("TNM02/0612Z")
	if err != nil {
		t.Fatalf("parseTemperature error: %v", err)
	}
	assertEqual(t, tm.Temperature, float64(-2), "temperature")
	assertEqual(t, tm.Day, 6, "day")
	assertEqual(t, tm.Hour, 12, "hour")
}

// ============================================================
// TAFParser tests
// ============================================================

func TestTAFParserBasic(t *testing.T) {
	input := `TAF TXFL 150500Z 1506/1612 17005KT 6000 SCT012
TEMPO 1506/1509 3000 BR BKN006 PROB40
TEMPO 1506/1508 0400 BCFG BKN002 PROB40
TEMPO 1512/1516 4000 -SHRA FEW030TCU BKN040
BECMG 1520/1522 CAVOK
TEMPO 1603/1608 3000 BR BKN006 PROB40
TEMPO 1604/1607 0400 BCFG BKN002 TX17/1512Z TN07/1605Z`
	result, err := ParseTAF(input, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, result.Station, "TXFL", "station")
	assertIntPtr(t, result.Day, 15, "day")
	assertIntPtr(t, result.Hour, 5, "hour")
	assertIntPtr(t, result.Minute, 0, "minute")
	assertEqual(t, result.Validity.StartDay, 15, "validity.startDay")
	assertEqual(t, result.Validity.StartHour, 6, "validity.startHour")
	assertEqual(t, result.Validity.EndDay, 16, "validity.endDay")
	assertEqual(t, result.Validity.EndHour, 12, "validity.endHour")
	assertNotNil(t, result.Wind, "wind")
	assertIntPtr(t, result.Wind.Degrees, 170, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 5, "wind.speed")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(6000), "vis.value")
	assertEqual(t, len(result.Clouds), 1, "clouds length")
	assertEqual(t, result.Clouds[0].Quantity, CloudQuantitySCT, "cloud[0].quantity")
	assertIntPtr(t, result.Clouds[0].Height, 1200, "cloud[0].height")
	assertEqual(t, len(result.Trends), 6, "trends length")
	assertEqual(t, result.Message, input, "message")
	assertEqual(t, result.InitialRaw, "TAF TXFL 150500Z 1506/1612 17005KT 6000 SCT012", "initialRaw")
}

func TestTAFParserWithoutLineBreaks(t *testing.T) {
	input := "TAF LSZH 292025Z 2921/3103 VRB03KT 9999 FEW020 BKN080 TX20/3014Z TN06/3003Z PROB30 TEMPO 2921/2923 SHRA BECMG 3001/3004 4000 MIFG NSC PROB40 3003/3007 1500 BCFG SCT004 PROB30 3004/3007 0800 FG VV003 BECMG 3006/3009 9999 FEW030 PROB40 TEMPO 3012/3017 30008KT"
	result, err := ParseTAF(input, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertIntPtr(t, result.Day, 29, "day")
	assertIntPtr(t, result.Hour, 20, "hour")
	assertIntPtr(t, result.Minute, 25, "minute")
	assertEqual(t, result.Validity.StartDay, 29, "validity.startDay")
	assertEqual(t, result.Validity.StartHour, 21, "validity.startHour")
	assertEqual(t, result.Validity.EndDay, 31, "validity.endDay")
	assertEqual(t, result.Validity.EndHour, 3, "validity.endHour")
	assertNil(t, result.Wind.Degrees, "wind.degrees")
	assertEqual(t, result.Wind.Direction, "VRB", "wind.direction")
	assertEqual(t, result.Wind.Speed, 3, "wind.speed")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
	assertNil(t, result.Wind.Gust, "wind.gust")
	assertNotNil(t, result.Visibility, "visibility")
	assertEqual(t, result.Visibility.Value, float64(9999), "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitMeters, "vis.unit")
	assertNotNil(t, result.Visibility.Indicator, "vis.indicator")
	assertEqual(t, *result.Visibility.Indicator, ValueIndicatorGreaterThan, "vis.indicator")
	assertEqual(t, len(result.Clouds), 2, "clouds length")
	if len(result.Clouds) >= 2 {
		assertEqual(t, result.Clouds[0].Quantity, CloudQuantityFEW, "clouds[0].quantity")
		assertIntPtr(t, result.Clouds[0].Height, 2000, "clouds[0].height")
		assertNil(t, result.Clouds[0].Type, "clouds[0].type")
		assertEqual(t, result.Clouds[1].Quantity, CloudQuantityBKN, "clouds[1].quantity")
		assertIntPtr(t, result.Clouds[1].Height, 8000, "clouds[1].height")
		assertNil(t, result.Clouds[1].Type, "clouds[1].type")
	}
	if result.MaxTemperature != nil {
		assertEqual(t, result.MaxTemperature.Day, 30, "maxTemp.day")
		assertEqual(t, result.MaxTemperature.Hour, 14, "maxTemp.hour")
		assertEqual(t, result.MaxTemperature.Temperature, float64(20), "maxTemp.temp")
	}
	if result.MinTemperature != nil {
		assertEqual(t, result.MinTemperature.Day, 30, "minTemp.day")
		assertEqual(t, result.MinTemperature.Hour, 3, "minTemp.hour")
		assertEqual(t, result.MinTemperature.Temperature, float64(6), "minTemp.temp")
	}
	assertEqual(t, len(result.Trends), 6, "trends length")
	if len(result.Trends) < 6 {
		return
	}
	// Trend 0: PROB30 TEMPO 2921/2923 SHRA
	{
		tr := result.Trends[0]
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[0].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[0].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 29, "trend[0].validity.startDay")
			assertEqual(t, v.StartHour, 21, "trend[0].validity.startHour")
			assertEqual(t, v.EndDay, 29, "trend[0].validity.endDay")
			assertEqual(t, v.EndHour, 23, "trend[0].validity.endHour")
		}
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[0].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityModerate, "trend[0].wc[0].intensity")
			assertNotNil(t, tr.WeatherConditions[0].Descriptive, "trend[0].wc[0].descriptive")
			assertEqual(t, *tr.WeatherConditions[0].Descriptive, DescriptiveShowers, "trend[0].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[0].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonRain, "trend[0].wc[0].phenomenon[0]")
			}
		}
		assertIntPtr(t, tr.Probability, 30, "trend[0].prob")
		assertEqual(t, tr.Raw, "PROB30 TEMPO 2921/2923 SHRA", "trend[0].raw")
	}
	// Trend 1: BECMG 3001/3004 4000 MIFG NSC
	{
		tr := result.Trends[1]
		assertEqual(t, tr.Type, WeatherChangeTypeBECMG, "trend[1].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[1].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 30, "trend[1].validity.startDay")
			assertEqual(t, v.StartHour, 1, "trend[1].validity.startHour")
			assertEqual(t, v.EndDay, 30, "trend[1].validity.endDay")
			assertEqual(t, v.EndHour, 4, "trend[1].validity.endHour")
		}
		assertNotNil(t, tr.Visibility, "trend[1].vis")
		assertEqual(t, tr.Visibility.Value, 4000.0, "trend[1].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[1].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[1].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityModerate, "trend[1].wc[0].intensity")
			assertNotNil(t, tr.WeatherConditions[0].Descriptive, "trend[1].wc[0].descriptive")
			assertEqual(t, *tr.WeatherConditions[0].Descriptive, DescriptiveShallow, "trend[1].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[1].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonFog, "trend[1].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 1, "trend[1].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityNSC, "trend[1].clouds[0].quantity")
		}
		assertEqual(t, tr.Raw, "BECMG 3001/3004 4000 MIFG NSC", "trend[1].raw")
	}
	// Trend 2: PROB40 3003/3007 1500 BCFG SCT004
	{
		tr := result.Trends[2]
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[2].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 30, "trend[2].validity.startDay")
			assertEqual(t, v.StartHour, 3, "trend[2].validity.startHour")
			assertEqual(t, v.EndDay, 30, "trend[2].validity.endDay")
			assertEqual(t, v.EndHour, 7, "trend[2].validity.endHour")
		}
		assertNotNil(t, tr.Visibility, "trend[2].vis")
		assertEqual(t, tr.Visibility.Value, 1500.0, "trend[2].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[2].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[2].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityModerate, "trend[2].wc[0].intensity")
			assertNotNil(t, tr.WeatherConditions[0].Descriptive, "trend[2].wc[0].descriptive")
			assertEqual(t, *tr.WeatherConditions[0].Descriptive, DescriptivePatches, "trend[2].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[2].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonFog, "trend[2].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 1, "trend[2].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantitySCT, "trend[2].clouds[0].quantity")
			assertNil(t, tr.Clouds[0].Type, "trend[2].clouds[0].type")
		}
		assertIntPtr(t, tr.Probability, 40, "trend[2].prob")
		assertEqual(t, tr.Raw, "PROB40 3003/3007 1500 BCFG SCT004", "trend[2].raw")
	}
	// Trend 3: PROB30 3004/3007 0800 FG VV003
	{
		tr := result.Trends[3]
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[3].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 30, "trend[3].validity.startDay")
			assertEqual(t, v.StartHour, 4, "trend[3].validity.startHour")
			assertEqual(t, v.EndDay, 30, "trend[3].validity.endDay")
			assertEqual(t, v.EndHour, 7, "trend[3].validity.endHour")
		}
		assertNotNil(t, tr.Visibility, "trend[3].vis")
		assertEqual(t, tr.Visibility.Value, 800.0, "trend[3].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[3].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[3].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityModerate, "trend[3].wc[0].intensity")
			assertNil(t, tr.WeatherConditions[0].Descriptive, "trend[3].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[3].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonFog, "trend[3].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 0, "trend[3].clouds")
		assertNotNil(t, tr.VerticalVisibility, "trend[3].vertVis")
		assertIntPtr(t, tr.VerticalVisibility, 300, "trend[3].vertVis.value")
		assertIntPtr(t, tr.Probability, 30, "trend[3].prob")
		assertEqual(t, tr.Raw, "PROB30 3004/3007 0800 FG VV003", "trend[3].raw")
	}
	// Trend 4: BECMG 3006/3009 9999 FEW030
	{
		tr := result.Trends[4]
		assertEqual(t, tr.Type, WeatherChangeTypeBECMG, "trend[4].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[4].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 30, "trend[4].validity.startDay")
			assertEqual(t, v.StartHour, 6, "trend[4].validity.startHour")
			assertEqual(t, v.EndDay, 30, "trend[4].validity.endDay")
			assertEqual(t, v.EndHour, 9, "trend[4].validity.endHour")
		}
		assertNotNil(t, tr.Visibility, "trend[4].vis")
		assertEqual(t, tr.Visibility.Value, 9999.0, "trend[4].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[4].vis.unit")
		assertNotNil(t, tr.Visibility.Indicator, "trend[4].vis.indicator")
		assertEqual(t, *tr.Visibility.Indicator, ValueIndicatorGreaterThan, "trend[4].vis.indicator")
		assertEqual(t, len(tr.WeatherConditions), 0, "trend[4].wc")
		assertEqual(t, len(tr.Clouds), 1, "trend[4].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityFEW, "trend[4].clouds[0].quantity")
			assertIntPtr(t, tr.Clouds[0].Height, 3000, "trend[4].clouds[0].height")
			assertNil(t, tr.Clouds[0].Type, "trend[4].clouds[0].type")
		}
		assertEqual(t, tr.Raw, "BECMG 3006/3009 9999 FEW030", "trend[4].raw")
	}
	// Trend 5: PROB40 TEMPO 3012/3017 30008KT
	{
		tr := result.Trends[5]
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[5].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[5].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 30, "trend[5].validity.startDay")
			assertEqual(t, v.StartHour, 12, "trend[5].validity.startHour")
			assertEqual(t, v.EndDay, 30, "trend[5].validity.endDay")
			assertEqual(t, v.EndHour, 17, "trend[5].validity.endHour")
		}
		assertEqual(t, len(tr.WeatherConditions), 0, "trend[5].wc")
		assertNotNil(t, tr.Wind, "trend[5].wind")
		assertIntPtr(t, tr.Wind.Degrees, 300, "trend[5].wind.degrees")
		assertEqual(t, tr.Wind.Speed, 8, "trend[5].wind.speed")
		assertNil(t, tr.Wind.Gust, "trend[5].wind.gust")
		assertEqual(t, tr.Wind.Unit, SpeedUnitKnot, "trend[5].wind.unit")
		assertIntPtr(t, tr.Probability, 40, "trend[5].prob")
		assertEqual(t, tr.Raw, "PROB40 TEMPO 3012/3017 30008KT", "trend[5].raw")
	}
}

func TestTAFParserWindShear(t *testing.T) {
	input := "TAF KMKE 011530 0116/0218 WS020/24045KT FM010200 17005KT P6SM SKC WS020/23055KT"
	result, err := ParseTAF(input, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertNotNil(t, result.WindShear, "windShear")
	assertEqual(t, result.WindShear.Height, 2000, "windShear.height")
	assertIntPtr(t, result.WindShear.Degrees, 240, "windShear.degrees")
	assertEqual(t, result.WindShear.Speed, 45, "windShear.speed")
	assertEqual(t, len(result.Trends), 1, "trends length")
	if len(result.Trends) > 0 {
		assertNotNil(t, result.Trends[0].WindShear, "trend.windShear")
	}
}

func TestTAFParserTurbulence(t *testing.T) {
	input := "TAF KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 530009 QNH2941INS"
	result, err := ParseTAF(input, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Turbulence), 1, "turbulence length")
	if len(result.Turbulence) > 0 {
		assertEqual(t, result.Turbulence[0].Intensity, TurbulenceIntensityModerateClearAirFreq, "turb[0].intensity")
		assertEqual(t, result.Turbulence[0].BaseHeight, 0, "turb[0].baseHeight")
		assertEqual(t, result.Turbulence[0].Depth, 9000, "turb[0].depth")
	}
}

func TestTAFParserIcing(t *testing.T) {
	input := "TAF KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 620304 QNH2941INS"
	result, err := ParseTAF(input, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Icing), 1, "icing length")
	if len(result.Icing) > 0 {
		assertEqual(t, result.Icing[0].Intensity, IcingIntensityLightRimeIcingCloud, "icing[0].intensity")
		assertEqual(t, result.Icing[0].BaseHeight, 3000, "icing[0].baseHeight")
		assertEqual(t, result.Icing[0].Depth, 4000, "icing[0].depth")
	}
}

func TestTAFParserCanceled(t *testing.T) {
	input := "TAF VTBD 281000Z 2812/2912 CNL="
	result, err := ParseTAF(input, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertBoolPtr(t, result.Canceled, true, "canceled")
	assertEqual(t, result.Station, "VTBD", "station")
	assertEqual(t, result.Validity.StartDay, 28, "validity.startDay")
	assertEqual(t, result.Validity.EndDay, 29, "validity.endDay")
}

func TestTAFParserCorrected(t *testing.T) {
	input := "TAF COR EDDS 201148Z 2012/2112 31010KT CAVOK BECMG 2018/2021 33004KT BECMG 2106/2109 07005KT"
	result, err := ParseTAF(input, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertBoolPtr(t, result.Corrected, true, "corrected")
	assertEqual(t, result.Station, "EDDS", "station")
	assertEqual(t, result.Validity.StartDay, 20, "validity.startDay")
	assertEqual(t, result.Validity.EndHour, 12, "validity.endHour")
}

// ---- Missing TAFParser tests ----

func TestTAFParserInvalidLineBreaks(t *testing.T) {
	code := "TAF TXFL 150500Z 1506/1612 17005KT 6000 SCT012 \n" +
		"TEMPO 1506/1509 3000 BR BKN006 PROB40 \n" +
		"TEMPO 1506/1508 0400 BCFG BKN002 PROB40 \n" +
		"TEMPO 1512/1516 4000 -SHRA FEW030TCU BKN040 \n" +
		"BECMG 1520/1522 CAVOK \n" +
		"TEMPO 1603/1608 3000 BR BKN006 PROB40 \n" +
		"TEMPO 1604/1607 0400 BCFG BKN002 TX17/1512Z TN07/1605Z"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, result.Station, "TXFL", "station")
	assertIntPtr(t, result.Day, 15, "day")
	assertIntPtr(t, result.Hour, 5, "hour")
	assertEqual(t, result.Validity.StartDay, 15, "valid.startDay")
	assertEqual(t, result.Validity.StartHour, 6, "valid.startHour")
	assertEqual(t, result.Validity.EndDay, 16, "valid.endDay")
	assertEqual(t, result.Validity.EndHour, 12, "valid.endHour")
	assertNotNil(t, result.Wind, "wind")
	assertEqual(t, *result.Wind.Degrees, 170, "wind.degrees")
	assertEqual(t, result.Wind.Speed, 5, "wind.speed")
	assertEqual(t, result.Visibility.Value, 6000.0, "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitMeters, "vis.unit")
	assertEqual(t, len(result.Clouds), 1, "clouds")
	assertEqual(t, result.Clouds[0].Quantity, CloudQuantitySCT, "cloud.quantity")
	assertEqual(t, len(result.Trends), 6, "trends")
	if len(result.Trends) < 6 {
		return
	}
	// Trend 0: TEMPO 1506/1509 3000 BR BKN006
	{
		tr := result.Trends[0]
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[0].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[0].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 15, "trend[0].validity.startDay")
			assertEqual(t, v.StartHour, 6, "trend[0].validity.startHour")
			assertEqual(t, v.EndDay, 15, "trend[0].validity.endDay")
			assertEqual(t, v.EndHour, 9, "trend[0].validity.endHour")
		}
		assertNotNil(t, tr.Visibility, "trend[0].vis")
		assertEqual(t, tr.Visibility.Value, 3000.0, "trend[0].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[0].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[0].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityModerate, "trend[0].wc[0].intensity")
			assertNil(t, tr.WeatherConditions[0].Descriptive, "trend[0].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[0].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonMist, "trend[0].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 1, "trend[0].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityBKN, "trend[0].clouds[0].quantity")
			assertNil(t, tr.Clouds[0].Type, "trend[0].clouds[0].type")
		}
		assertNil(t, tr.Probability, "trend[0].prob")
		assertEqual(t, tr.Raw, "TEMPO 1506/1509 3000 BR BKN006", "trend[0].raw")
	}
	// Trend 1: PROB40 TEMPO 1506/1508 0400 BCFG BKN002
	{
		tr := result.Trends[1]
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[1].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[1].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 15, "trend[1].validity.startDay")
			assertEqual(t, v.StartHour, 6, "trend[1].validity.startHour")
			assertEqual(t, v.EndDay, 15, "trend[1].validity.endDay")
			assertEqual(t, v.EndHour, 8, "trend[1].validity.endHour")
		}
		assertEqual(t, tr.Visibility.Value, 400.0, "trend[1].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[1].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[1].wc")
		if len(tr.WeatherConditions) > 0 {
			assertNotNil(t, tr.WeatherConditions[0].Descriptive, "trend[1].wc[0].descriptive")
			assertEqual(t, *tr.WeatherConditions[0].Descriptive, DescriptivePatches, "trend[1].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[1].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonFog, "trend[1].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 1, "trend[1].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityBKN, "trend[1].clouds[0].quantity")
			assertIntPtr(t, tr.Clouds[0].Height, 200, "trend[1].clouds[0].height")
		}
		assertIntPtr(t, tr.Probability, 40, "trend[1].prob")
		assertEqual(t, tr.Raw, "PROB40 TEMPO 1506/1508 0400 BCFG BKN002", "trend[1].raw")
	}
	// Trend 2: PROB40 TEMPO 1512/1516 4000 -SHRA FEW030TCU BKN040
	{
		tr := result.Trends[2]
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[2].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[2].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 15, "trend[2].validity.startDay")
			assertEqual(t, v.StartHour, 12, "trend[2].validity.startHour")
			assertEqual(t, v.EndDay, 15, "trend[2].validity.endDay")
			assertEqual(t, v.EndHour, 16, "trend[2].validity.endHour")
		}
		assertEqual(t, tr.Visibility.Value, 4000.0, "trend[2].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[2].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[2].wc")
		if len(tr.WeatherConditions) > 0 {
			assertNotNil(t, tr.WeatherConditions[0].Intensity, "trend[2].wc[0].intensity")
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityLight, "trend[2].wc[0].intensity")
			assertNotNil(t, tr.WeatherConditions[0].Descriptive, "trend[2].wc[0].descriptive")
			assertEqual(t, *tr.WeatherConditions[0].Descriptive, DescriptiveShowers, "trend[2].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[2].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonRain, "trend[2].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 2, "trend[2].clouds")
		if len(tr.Clouds) >= 2 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityFEW, "trend[2].clouds[0].quantity")
			assertIntPtr(t, tr.Clouds[0].Height, 3000, "trend[2].clouds[0].height")
			if tr.Clouds[0].Type != nil {
				assertEqual(t, *tr.Clouds[0].Type, CloudTypeTCU, "trend[2].clouds[0].type")
			}
			assertEqual(t, tr.Clouds[1].Quantity, CloudQuantityBKN, "trend[2].clouds[1].quantity")
			assertIntPtr(t, tr.Clouds[1].Height, 4000, "trend[2].clouds[1].height")
			assertNil(t, tr.Clouds[1].Type, "trend[2].clouds[1].type")
		}
		assertIntPtr(t, tr.Probability, 40, "trend[2].prob")
		assertEqual(t, tr.Raw, "PROB40 TEMPO 1512/1516 4000 -SHRA FEW030TCU BKN040", "trend[2].raw")
	}
	// Trend 3: BECMG 1520/1522 CAVOK
	{
		tr := result.Trends[3]
		assertEqual(t, tr.Type, WeatherChangeTypeBECMG, "trend[3].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[3].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 15, "trend[3].validity.startDay")
			assertEqual(t, v.StartHour, 20, "trend[3].validity.startHour")
			assertEqual(t, v.EndDay, 15, "trend[3].validity.endDay")
			assertEqual(t, v.EndHour, 22, "trend[3].validity.endHour")
		}
		assertEqual(t, tr.Raw, "BECMG 1520/1522 CAVOK", "trend[3].raw")
	}
	// Trend 4: TEMPO 1603/1608 3000 BR BKN006
	{
		tr := result.Trends[4]
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[4].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 16, "trend[4].validity.startDay")
			assertEqual(t, v.StartHour, 3, "trend[4].validity.startHour")
			assertEqual(t, v.EndDay, 16, "trend[4].validity.endDay")
			assertEqual(t, v.EndHour, 8, "trend[4].validity.endHour")
		}
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[4].type")
		assertEqual(t, tr.Visibility.Value, 3000.0, "trend[4].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[4].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[4].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityModerate, "trend[4].wc[0].intensity")
			assertNil(t, tr.WeatherConditions[0].Descriptive, "trend[4].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[4].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonMist, "trend[4].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 1, "trend[4].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityBKN, "trend[4].clouds[0].quantity")
			assertNil(t, tr.Clouds[0].Type, "trend[4].clouds[0].type")
		}
		assertNil(t, tr.Probability, "trend[4].prob")
		assertEqual(t, tr.Raw, "TEMPO 1603/1608 3000 BR BKN006", "trend[4].raw")
	}
	// Trend 5: PROB40 TEMPO 1604/1607 0400 BCFG BKN002 TX17/1512Z TN07/1605Z
	{
		tr := result.Trends[5]
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[5].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[5].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 16, "trend[5].validity.startDay")
			assertEqual(t, v.StartHour, 4, "trend[5].validity.startHour")
			assertEqual(t, v.EndDay, 16, "trend[5].validity.endDay")
			assertEqual(t, v.EndHour, 7, "trend[5].validity.endHour")
		}
		assertEqual(t, tr.Visibility.Value, 400.0, "trend[5].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[5].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[5].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, *tr.WeatherConditions[0].Intensity, IntensityModerate, "trend[5].wc[0].intensity")
			assertNotNil(t, tr.WeatherConditions[0].Descriptive, "trend[5].wc[0].descriptive")
			assertEqual(t, *tr.WeatherConditions[0].Descriptive, DescriptivePatches, "trend[5].wc[0].descriptive")
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 1, "trend[5].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) > 0 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonFog, "trend[5].wc[0].phenomenon[0]")
			}
		}
		assertEqual(t, len(tr.Clouds), 1, "trend[5].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityBKN, "trend[5].clouds[0].quantity")
			assertNil(t, tr.Clouds[0].Type, "trend[5].clouds[0].type")
		}
		assertIntPtr(t, tr.Probability, 40, "trend[5].prob")
		assertEqual(t, tr.Raw, "PROB40 TEMPO 1604/1607 0400 BCFG BKN002 TX17/1512Z TN07/1605Z", "trend[5].raw")
	}
}

func TestTAFParserNoLineBreaksEndTemp(t *testing.T) {
	code := "TAF KLSV 120700Z 1207/1313 VRB06KT 9999 SCT250 QNH2992INS BECMG 1217/1218 10010G15KT 9999 SCT250 QNH2980INS BECMG 1303/1304 VRB06KT 9999 FEW250 QNH2979INS TX42/1223Z TN24/1213Z"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertIntPtr(t, result.Day, 12, "day")
	assertIntPtr(t, result.Hour, 7, "hour")
	assertEqual(t, result.Wind.Direction, "VRB", "wind.direction")
	assertEqual(t, result.Wind.Speed, 6, "wind.speed")
	assertNil(t, result.Wind.Gust, "wind.gust")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
	assertNotNil(t, result.MaxTemperature, "maxTemp")
	assertEqual(t, result.MaxTemperature.Temperature, 42.0, "maxTemp.temp")
	assertEqual(t, result.MaxTemperature.Hour, 23, "maxTemp.hour")
	assertNotNil(t, result.MinTemperature, "minTemp")
	assertEqual(t, result.MinTemperature.Temperature, 24.0, "minTemp.temp")
	assertEqual(t, result.MinTemperature.Hour, 13, "minTemp.hour")
	assertEqual(t, len(result.Trends), 2, "trends")
	if len(result.Trends) < 2 {
		return
	}
	// Trend 0: BECMG 1217/1218 10010G15KT 9999 SCT250
	{
		tr := result.Trends[0]
		assertEqual(t, tr.Type, WeatherChangeTypeBECMG, "trend[0].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[0].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 12, "trend[0].validity.startDay")
			assertEqual(t, v.StartHour, 17, "trend[0].validity.startHour")
			assertEqual(t, v.EndDay, 12, "trend[0].validity.endDay")
			assertEqual(t, v.EndHour, 18, "trend[0].validity.endHour")
		}
		assertNotNil(t, tr.Visibility, "trend[0].vis")
		assertEqual(t, tr.Visibility.Value, 9999.0, "trend[0].vis.value")
		assertNotNil(t, tr.Visibility.Indicator, "trend[0].vis.indicator")
		assertEqual(t, *tr.Visibility.Indicator, ValueIndicatorGreaterThan, "trend[0].vis.indicator")
		assertNotNil(t, tr.Wind, "trend[0].wind")
		assertIntPtr(t, tr.Wind.Degrees, 100, "trend[0].wind.degrees")
		assertEqual(t, tr.Wind.Speed, 10, "trend[0].wind.speed")
		assertIntPtr(t, tr.Wind.Gust, 15, "trend[0].wind.gust")
		assertEqual(t, tr.Wind.Unit, SpeedUnitKnot, "trend[0].wind.unit")
		assertEqual(t, len(tr.WeatherConditions), 0, "trend[0].wc")
		assertEqual(t, len(tr.Clouds), 1, "trend[0].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantitySCT, "trend[0].clouds[0].quantity")
			assertIntPtr(t, tr.Clouds[0].Height, 25000, "trend[0].clouds[0].height")
			assertNil(t, tr.Clouds[0].Type, "trend[0].clouds[0].type")
		}
	}
	// Trend 1: BECMG 1303/1304 VRB06KT 9999 FEW250
	{
		tr := result.Trends[1]
		assertEqual(t, tr.Type, WeatherChangeTypeBECMG, "trend[1].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[1].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 13, "trend[1].validity.startDay")
			assertEqual(t, v.StartHour, 3, "trend[1].validity.startHour")
			assertEqual(t, v.EndDay, 13, "trend[1].validity.endDay")
			assertEqual(t, v.EndHour, 4, "trend[1].validity.endHour")
		}
		assertNotNil(t, tr.Visibility, "trend[1].vis")
		assertEqual(t, tr.Visibility.Value, 9999.0, "trend[1].vis.value")
		assertNotNil(t, tr.Visibility.Indicator, "trend[1].vis.indicator")
		assertEqual(t, *tr.Visibility.Indicator, ValueIndicatorGreaterThan, "trend[1].vis.indicator")
		assertNotNil(t, tr.Wind, "trend[1].wind")
		assertNil(t, tr.Wind.Degrees, "trend[1].wind.degrees")
		assertEqual(t, tr.Wind.Direction, "VRB", "trend[1].wind.direction")
		assertEqual(t, tr.Wind.Speed, 6, "trend[1].wind.speed")
		assertNil(t, tr.Wind.Gust, "trend[1].wind.gust")
		assertEqual(t, tr.Wind.Unit, SpeedUnitKnot, "trend[1].wind.unit")
		assertEqual(t, len(tr.WeatherConditions), 0, "trend[1].wc")
		assertEqual(t, len(tr.Clouds), 1, "trend[1].clouds")
		if len(tr.Clouds) > 0 {
			assertEqual(t, tr.Clouds[0].Quantity, CloudQuantityFEW, "trend[1].clouds[0].quantity")
			assertIntPtr(t, tr.Clouds[0].Height, 25000, "trend[1].clouds[0].height")
			assertNil(t, tr.Clouds[0].Type, "trend[1].clouds[0].type")
		}
	}
}

func TestTAFParserDoubleTAF(t *testing.T) {
	code := "TAF TAF LFPG 191100Z 1912/2018 02010KT 9999 FEW040 PROB30 1217/1218"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertNotNil(t, result, "result")
	assertEqual(t, len(result.Trends), 1, "trends")
	assertIntPtr(t, result.Trends[0].Probability, 30, "trend[0].prob")
	assertEqual(t, result.InitialRaw, "TAF TAF LFPG 191100Z 1912/2018 02010KT 9999 FEW040", "initialRaw")
}

func TestTAFParserNauticalMilesVis(t *testing.T) {
	code := "TAF AMD CZBF 300939Z 3010/3022 VRB03KT 6SM -SN OVC015 TEMPO 3010/3012 11/2SM -SN OVC009 \nFM301200 10008KT 2SM -SN OVC010 TEMPO 3012/3022 3/4SM -SN VV007 RMK FCST BASED ON AUTO OBS. NXT FCST BY 301400Z"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertNotNil(t, result.Visibility, "vis")
	assertEqual(t, result.Visibility.Value, 6.0, "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitStatuteMiles, "vis.unit")
	assertBoolPtr(t, result.Amendment, true, "amendment")
	assertEqual(t, len(result.Trends), 3, "trends")
	if len(result.Trends) >= 1 {
		assertNotNil(t, result.Trends[0].Visibility, "trend[0].vis")
		assertEqual(t, result.Trends[0].Visibility.Value, 5.5, "trend[0].vis.value")
		assertEqual(t, result.Trends[0].Visibility.Unit, DistanceUnitStatuteMiles, "trend[0].vis.unit")
	}
	if len(result.Trends) >= 2 {
		assertNotNil(t, result.Trends[1].Visibility, "trend[1].vis")
		assertEqual(t, result.Trends[1].Visibility.Value, 2.0, "trend[1].vis.value")
		assertEqual(t, result.Trends[1].Visibility.Unit, DistanceUnitStatuteMiles, "trend[1].vis.unit")
	}
	if len(result.Trends) >= 3 {
		assertNotNil(t, result.Trends[2].Visibility, "trend[2].vis")
		assertEqual(t, result.Trends[2].Visibility.Value, 0.75, "trend[2].vis.value")
		assertEqual(t, result.Trends[2].Visibility.Unit, DistanceUnitStatuteMiles, "trend[2].vis.unit")
	}
}

func TestTAFParserThunderstorms(t *testing.T) {
	code := "TAF AMD KGWO 161553Z 1616/1712 21005KT 4SM -TSRA BR SCT010 OVC070CB"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.WeatherConditions), 2, "wc")
	if len(result.WeatherConditions) > 0 {
		assertEqual(t, *result.WeatherConditions[0].Descriptive, DescriptiveThunderstorm, "wc.descriptive")
		assertNotNil(t, result.WeatherConditions[0].Intensity, "wc.intensity")
		assertEqual(t, *result.WeatherConditions[0].Intensity, IntensityLight, "wc.intensity")
		assertEqual(t, result.WeatherConditions[0].Phenomena[0], PhenomenonRain, "wc.phenomenon")
	}
}

func TestTAFParserRemark(t *testing.T) {
	code := "TAF CZBF 300939Z 3010/3022 VRB03KT 6SM -SN OVC015 RMK FCST BASED ON AUTO OBS. NXT FCST BY 301400Z\n TEMPO 3010/3012 11/2SM -SN OVC009 FM301200 10008KT 2SM -SN OVC010 \nTEMPO 3012/3022 3/4SM -SN VV007"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertNotNil(t, result.Remark, "remark")
	assertEqual(t, len(result.Remarks), 2, "remarks len")
}

func TestTAFParserTrendRemark(t *testing.T) {
	code := "TAF CZBF 300939Z 3010/3022 VRB03KT 6SM -SN OVC015\n TEMPO 3010/3012 11/2SM -SN OVC009 FM301200 10008KT 2SM -SN OVC010 TEMPO 3012/3022 3/4SM -SN VV007 RMK FCST BASED ON AUTO OBS. NXT FCST BY 301400Z"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 3, "trends")
	if len(result.Trends) >= 3 {
		assertNotNil(t, result.Trends[2].Remark, "trend[2].remark")
		assertEqual(t, len(result.Trends[2].Remarks), 2, "trend[2].remarks")
	}
}

func TestTAFParserInter(t *testing.T) {
	code := "TAF TAF AMD YWLM 270723Z 2707/2806 19020G30KT 9999 -SHRA SCT015 BKN020 BECMG 2708/2710 19014KT 9999 -SHRA SCT010 BKN015 BECMG 2800/2802 18015G25KT 9999 -SHRA SCT015 BKN020 TEMPO 2707/2712 3000 SHRA SCT005 BKN010 INTER 2712/2802 4000 SHRA SCT005 BKN010"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 4, "trends")
	if len(result.Trends) >= 4 {
		assertEqual(t, result.Trends[3].Type, WeatherChangeTypeINTER, "trend[3].type")
	}
}

func TestTAFParserInterProbability(t *testing.T) {
	code := "TAF YWLM 270209Z 2703/2800 30014KT 9999 -SHRA NSC FM270400 28007KT 9999 -SHRA SCT040 FM270700 03010KT 9999 -SHRA SCT040 FM271200 30008KT CAVOK FM272100 29014KT CAVOK INTER 2703/2709 30018G30KT 5000 SHRA SCT015 BKN040 FEW040TCU PROB30 INTER 2704/2709 VRB25G45KT 2000 TSRAGR BKN010 SCT040CB"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 6, "trends")
	if len(result.Trends) >= 6 {
		assertEqual(t, result.Trends[4].Type, WeatherChangeTypeINTER, "trend[4].type")
		assertNil(t, result.Trends[4].Probability, "trend[4].prob")
		assertEqual(t, result.Trends[5].Type, WeatherChangeTypeINTER, "trend[5].type")
		assertIntPtr(t, result.Trends[5].Probability, 30, "trend[5].prob")
	}
}

func TestTAFParserStopsWCRemark(t *testing.T) {
	code := "TAF CYTL 121940Z 1220/1308 RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z TEMPO 1303/1308 2SM -SN"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 1, "trends")
	assertEqual(t, len(result.WeatherConditions), 0, "wc")
}

func TestTAFParserStopsWCTrendRemark(t *testing.T) {
	const tafCode = "TAF CYTL 121940Z 1220/1308 TEMPO 1303/1308 2SM -SN RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z"
	result, err := ParseTAF(tafCode, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 1, "trends")
	if len(result.Trends) > 0 {
		assertEqual(t, len(result.Trends[0].WeatherConditions), 1, "trend[0].wc")
		assertNotNil(t, result.Trends[0].WeatherConditions[0].Intensity, "trend[0].wc.intensity")
		assertEqual(t, *result.Trends[0].WeatherConditions[0].Intensity, IntensityLight, "trend[0].wc.intensity")
		assertEqual(t, result.Trends[0].WeatherConditions[0].Phenomena[0], PhenomenonSnow, "trend[0].wc.phenomenon")
	}
}

func TestTAFParserSameWithoutTAF(t *testing.T) {
	withTAF := "TAF CYTL 121940Z 1220/1308 TEMPO 1303/1308 2SM -SN RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z"
	withoutTAF := "CYTL 121940Z 1220/1308 TEMPO 1303/1308 2SM -SN RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z"
	result1, _ := ParseTAF(withoutTAF, nil)
	result2, _ := ParseTAF(withTAF, nil)

	// Message and InitialRaw differ, but all parsed fields should match
	compareTAFResults(t, result1, result2)
}

func compareTAFResults(t *testing.T, r1, r2 *TAF) {
	t.Helper()
	assertEqual(t, r1.Station, r2.Station, "station")
	assertEqual(t, r1.Validity, r2.Validity, "validity")
	if r1.Day != nil && r2.Day != nil {
		assertIntPtr(t, r1.Day, *r2.Day, "day")
	} else if r1.Day != nil || r2.Day != nil {
		t.Error("day: one is nil, the other is not")
	}
	if r1.Hour != nil && r2.Hour != nil {
		assertIntPtr(t, r1.Hour, *r2.Hour, "hour")
	} else if r1.Hour != nil || r2.Hour != nil {
		t.Error("hour: one is nil, the other is not")
	}
	if r1.Minute != nil && r2.Minute != nil {
		assertIntPtr(t, r1.Minute, *r2.Minute, "minute")
	} else if r1.Minute != nil || r2.Minute != nil {
		t.Error("minute: one is nil, the other is not")
	}
	if r1.Wind != nil && r2.Wind != nil {
		assertEqual(t, r1.Wind.Speed, r2.Wind.Speed, "wind.speed")
		assertEqual(t, r1.Wind.Direction, r2.Wind.Direction, "wind.direction")
		assertEqual(t, r1.Wind.Unit, r2.Wind.Unit, "wind.unit")
	}
	assertEqual(t, len(r1.Clouds), len(r2.Clouds), "clouds len")
	assertEqual(t, len(r1.WeatherConditions), len(r2.WeatherConditions), "wc len")
	assertEqual(t, len(r1.Trends), len(r2.Trends), "trends len")
	if len(r1.Trends) > 0 && len(r2.Trends) > 0 {
		assertEqual(t, r1.Trends[0].Type, r2.Trends[0].Type, "trend[0].type")
		assertEqual(t, len(r1.Trends[0].WeatherConditions), len(r2.Trends[0].WeatherConditions), "trend[0].wc len")
	}
}

func TestTAFParserSameWithDoubleTAF(t *testing.T) {
	withTAF := "TAF CYTL 121940Z 1220/1308 TEMPO 1303/1308 2SM -SN RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z"
	doubleTAF := "TAF TAF CYTL 121940Z 1220/1308 TEMPO 1303/1308 2SM -SN RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z"
	result1, _ := ParseTAF(withTAF, nil)
	result2, _ := ParseTAF(doubleTAF, nil)
	compareTAFResults(t, result1, result2)
}

func TestTAFParserSameWithRandomPrefix(t *testing.T) {
	withTAF := "TAF CYTL 121940Z 1220/1308 TEMPO 1303/1308 2SM -SN RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z"
	doubleTAF := "somethingrandom TAF CYTL 121940Z 1220/1308 TEMPO 1303/1308 2SM -SN RMK FCST BASED ON AUTO OBS. FCST BASED ON OBS BY OTHER SRCS. WIND SENSOR INOP. NXT FCST BY 130200Z"
	result1, err1 := ParseTAF(withTAF, nil)
	if err1 != nil {
		t.Fatalf("ParseTAF error on normal input: %v", err1)
	}
	result2, err2 := ParseTAF(doubleTAF, nil)
	if err2 != nil {
		t.Fatalf("ParseTAF error on random prefix: %v", err2)
	}
	compareTAFResults(t, result1, result2)
}

func TestTAFParserPROBXXTEMPO(t *testing.T) {
	code := "TAF EDDS 281100Z 2812/2912 04008KT 9999 BKN035 BECMG 2818/2821 33005KT PROB30 TEMPO 2818/2824 4000 TSRA BKN025CB TEMPO 2918/2824 BKN025CB"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 3, "trends")
	if len(result.Trends) >= 1 {
		assertNil(t, result.Trends[0].Probability, "trend[0].prob")
	}
	if len(result.Trends) >= 2 {
		assertEqual(t, *result.Trends[1].Probability, 30, "trend[1].prob")
		assertEqual(t, len(result.Trends[1].Clouds), 1, "trend[1].clouds")
	}
	if len(result.Trends) >= 3 {
		assertNil(t, result.Trends[2].Probability, "trend[2].prob")
		assertEqual(t, len(result.Trends[2].Clouds), 1, "trend[2].clouds")
	}
}

func TestTAFParserNSW(t *testing.T) {
	code := "TAF FYWB 222200Z 2300/2400 36007KT 0700 FG OVC009 BECMG 2307/2309 9999 NSW BKN012"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 1, "trends")
	if len(result.Trends) > 0 {
		assertEqual(t, result.Trends[0].Type, WeatherChangeTypeBECMG, "trend[0].type")
		assertEqual(t, len(result.Trends[0].WeatherConditions), 1, "trend[0].wc")
		if len(result.Trends[0].WeatherConditions) > 0 {
			assertEqual(t, result.Trends[0].WeatherConditions[0].Phenomena[0], PhenomenonNoSignificantWeather, "trend[0].wc.phenomenon")
		}
	}
}

func TestTAFParserBogusWeather(t *testing.T) {
	code := "TAF AMD KVOK 232200Z 2322/2423 15015G20KT 9999 BKN035 BKN050 510008 QNH2966INS BECMG 2411/2412 17012KT 9000 -SHRA FEW006 BKN014 OVC021 510065 QNH2965INS TX24/2322Z TN16/2413Z LAST NO AMDS AFT 2322 NEXT 2409"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 1, "trends")
	if len(result.Trends) > 0 {
		assertEqual(t, len(result.Trends[0].WeatherConditions), 1, "trend[0].wc")
		if len(result.Trends[0].WeatherConditions) > 0 {
			assertNotNil(t, result.Trends[0].WeatherConditions[0].Intensity, "trend[0].wc.intensity")
			assertEqual(t, *result.Trends[0].WeatherConditions[0].Intensity, IntensityLight, "trend[0].wc.intensity")
			assertEqual(t, result.Trends[0].WeatherConditions[0].Phenomena[0], PhenomenonRain, "trend[0].wc.phenomenon")
		}
	}
}

func TestTAFParserSlashWeather(t *testing.T) {
	code := "VOMM 281700Z 2818/2924 05005KT 4000 -RA/BR SCT020 BKN100 TEMPO 2818/2824 3000 TSRA/RA SCT018 FEW025TCU/CB BKN080 BECMG 2905/2906 05015KT 6000 BECMG 2915/2916 05005KT 4000 -DZ/BR"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	if len(result.Trends) >= 3 && len(result.Trends[2].WeatherConditions) > 0 {
		wc := result.Trends[2].WeatherConditions[0]
		assertNotNil(t, wc.Intensity, "trend[2].wc.intensity")
		assertEqual(t, *wc.Intensity, IntensityLight, "trend[2].wc.intensity")
		assertEqual(t, wc.Phenomena[0], PhenomenonDrizzle, "trend[2].wc.phenomenon[0]")
		assertEqual(t, wc.Phenomena[1], PhenomenonMist, "trend[2].wc.phenomenon[1]")
		assertNil(t, wc.Descriptive, "trend[2].wc.descriptive")
	}
}

func TestTAFParserTAFAMDTAFAMD(t *testing.T) {
	code := "TAF AMD TAF AMD CYRB 290006Z 2900/2924 06008KT P6SM SKC TEMPO 2900/2909 4SM IC PROB30 2900/2909 2SM IC BR BKN003 FM290900 02015KT P6SM FEW010 RMK NXT FCST BY 290600Z"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 3, "trends")
}

func TestTAFParserTurbulenceWind(t *testing.T) {
	code := "TAF KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 530009 QNH2941INS"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, result.Wind.Speed, 20, "wind.speed")
	assertEqual(t, *result.Wind.Degrees, 210, "wind.degrees")
	assertIntPtr(t, result.Wind.Gust, 35, "wind.gust")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
}

func TestTAFParserTurbulenceBaseMulti(t *testing.T) {
	code := "TAF KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 530009 5X0304 QNH2941INS"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Turbulence), 2, "turbulence")
	if len(result.Turbulence) >= 2 {
		assertEqual(t, result.Turbulence[0].Intensity, TurbulenceIntensityModerateClearAirFreq, "turb[0].intensity")
		assertEqual(t, result.Turbulence[1].Intensity, TurbulenceIntensityExtreme, "turb[1].intensity")
		assertEqual(t, result.Turbulence[1].BaseHeight, 3000, "turb[1].base")
		assertEqual(t, result.Turbulence[1].Depth, 4000, "turb[1].depth")
	}
}

func TestTAFParserTurbulenceTrendMulti(t *testing.T) {
	code := "TAF AMD CYRB 290006Z 2900/2924 06008KT P6SM SKC TEMPO 2900/2909 4SM IC PROB30 2900/2909 2SM IC BR BKN003 530009 5X0304"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 2, "trends")
	if len(result.Trends) >= 2 {
		assertEqual(t, len(result.Trends[1].Turbulence), 2, "trend[1].turbulence")
		if len(result.Trends[1].Turbulence) >= 2 {
			assertEqual(t, result.Trends[1].Turbulence[0].Intensity, TurbulenceIntensityModerateClearAirFreq, "turb[0].intensity")
			assertEqual(t, result.Trends[1].Turbulence[1].Intensity, TurbulenceIntensityExtreme, "turb[1].intensity")
		}
	}
}

func TestTAFParserIcingWind(t *testing.T) {
	code := "TAF KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 620304 QNH2941INS"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, result.Wind.Speed, 20, "wind.speed")
	assertEqual(t, *result.Wind.Degrees, 210, "wind.degrees")
	assertIntPtr(t, result.Wind.Gust, 35, "wind.gust")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
}

func TestTAFParserIcingBaseMulti(t *testing.T) {
	code := "TAF KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 620304 670009 QNH2941INS"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Icing), 2, "icing")
	if len(result.Icing) >= 2 {
		assertEqual(t, result.Icing[0].Intensity, IcingIntensityLightRimeIcingCloud, "icing[0].intensity")
		assertEqual(t, result.Icing[1].Intensity, IcingIntensitySevereMixedIcing, "icing[1].intensity")
	}
}

func TestTAFParserIcingTrendMulti(t *testing.T) {
	code := "TAF AMD CYRB 290006Z 2900/2924 06008KT P6SM SKC TEMPO 2900/2909 4SM IC PROB30 2900/2909 2SM IC BR BKN003 620304 670009"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, len(result.Trends), 2, "trends")
	if len(result.Trends) >= 2 {
		assertEqual(t, len(result.Trends[1].Icing), 2, "trend[1].icing")
		if len(result.Trends[1].Icing) >= 2 {
			assertEqual(t, result.Trends[1].Icing[0].Intensity, IcingIntensityLightRimeIcingCloud, "icing[0].intensity")
			assertEqual(t, result.Trends[1].Icing[1].Intensity, IcingIntensitySevereMixedIcing, "icing[1].intensity")
		}
	}
}

func TestTAFParserFMStation(t *testing.T) {
	code := "TAF FMMI 082300Z 0900/1006 16006KT 9999 FEW017 BKN020 PROB30 TEMPO 0908/0916 4500 RADZ BECMG 0909/0911 10010KT BECMG 0918/0920 16006KT"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertEqual(t, result.Station, "FMMI", "station")
	assertIntPtr(t, result.Day, 8, "day")
	assertIntPtr(t, result.Hour, 23, "hour")
	assertIntPtr(t, result.Minute, 0, "minute")
	assertEqual(t, result.Validity.StartDay, 9, "validity.startDay")
	assertEqual(t, result.Validity.StartHour, 0, "validity.startHour")
	assertEqual(t, result.Validity.EndDay, 10, "validity.endDay")
	assertEqual(t, result.Validity.EndHour, 6, "validity.endHour")
	assertNotNil(t, result.Wind, "wind")
	assertIntPtr(t, result.Wind.Degrees, 160, "wind.degrees")
	assertEqual(t, result.Wind.Direction, "SSE", "wind.direction")
	assertEqual(t, result.Wind.Speed, 6, "wind.speed")
	assertNil(t, result.Wind.Gust, "wind.gust")
	assertEqual(t, result.Wind.Unit, SpeedUnitKnot, "wind.unit")
	assertNotNil(t, result.Visibility, "vis")
	assertEqual(t, result.Visibility.Value, 9999.0, "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitMeters, "vis.unit")
	assertNotNil(t, result.Visibility.Indicator, "vis.indicator")
	assertEqual(t, *result.Visibility.Indicator, ValueIndicatorGreaterThan, "vis.indicator")
	assertEqual(t, len(result.Clouds), 2, "clouds")
	if len(result.Clouds) >= 2 {
		assertEqual(t, result.Clouds[0].Quantity, CloudQuantityFEW, "clouds[0].quantity")
		assertIntPtr(t, result.Clouds[0].Height, 1700, "clouds[0].height")
		assertNil(t, result.Clouds[0].Type, "clouds[0].type")
		assertEqual(t, result.Clouds[1].Quantity, CloudQuantityBKN, "clouds[1].quantity")
		assertIntPtr(t, result.Clouds[1].Height, 2000, "clouds[1].height")
		assertNil(t, result.Clouds[1].Type, "clouds[1].type")
	}
	assertEqual(t, len(result.WeatherConditions), 0, "wc")
	assertNil(t, result.MaxTemperature, "maxTemp")
	assertNil(t, result.MinTemperature, "minTemp")
	assertEqual(t, len(result.Trends), 3, "trends")
	if len(result.Trends) < 3 {
		return
	}
	// Trend 0: PROB30 TEMPO 0908/0916 4500 RADZ
	{
		tr := result.Trends[0]
		assertEqual(t, tr.Type, WeatherChangeTypeTEMPO, "trend[0].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[0].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 9, "trend[0].validity.startDay")
			assertEqual(t, v.StartHour, 8, "trend[0].validity.startHour")
			assertEqual(t, v.EndDay, 9, "trend[0].validity.endDay")
			assertEqual(t, v.EndHour, 16, "trend[0].validity.endHour")
		}
		assertEqual(t, len(tr.Clouds), 0, "trend[0].clouds")
		assertIntPtr(t, tr.Probability, 30, "trend[0].prob")
		assertNotNil(t, tr.Visibility, "trend[0].vis")
		assertEqual(t, tr.Visibility.Value, 4500.0, "trend[0].vis.value")
		assertEqual(t, tr.Visibility.Unit, DistanceUnitMeters, "trend[0].vis.unit")
		assertEqual(t, len(tr.WeatherConditions), 1, "trend[0].wc")
		if len(tr.WeatherConditions) > 0 {
			assertEqual(t, len(tr.WeatherConditions[0].Phenomena), 2, "trend[0].wc[0].phenomenons")
			if len(tr.WeatherConditions[0].Phenomena) >= 2 {
				assertEqual(t, tr.WeatherConditions[0].Phenomena[0], PhenomenonRain, "trend[0].wc[0].phenomenon[0]")
				assertEqual(t, tr.WeatherConditions[0].Phenomena[1], PhenomenonDrizzle, "trend[0].wc[0].phenomenon[1]")
			}
		}
		assertNil(t, tr.Wind, "trend[0].wind")
		assertEqual(t, tr.Raw, "PROB30 TEMPO 0908/0916 4500 RADZ", "trend[0].raw")
	}
	// Trend 1: BECMG 0909/0911 10010KT
	{
		tr := result.Trends[1]
		assertEqual(t, tr.Type, WeatherChangeTypeBECMG, "trend[1].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[1].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 9, "trend[1].validity.startDay")
			assertEqual(t, v.StartHour, 9, "trend[1].validity.startHour")
			assertEqual(t, v.EndDay, 9, "trend[1].validity.endDay")
			assertEqual(t, v.EndHour, 11, "trend[1].validity.endHour")
		}
		assertEqual(t, len(tr.Clouds), 0, "trend[1].clouds")
		assertNil(t, tr.Probability, "trend[1].prob")
		assertNil(t, tr.Visibility, "trend[1].vis")
		assertEqual(t, len(tr.WeatherConditions), 0, "trend[1].wc")
		assertNotNil(t, tr.Wind, "trend[1].wind")
		assertIntPtr(t, tr.Wind.Degrees, 100, "trend[1].wind.degrees")
		assertEqual(t, tr.Wind.Direction, "E", "trend[1].wind.direction")
		assertEqual(t, tr.Wind.Speed, 10, "trend[1].wind.speed")
		assertEqual(t, tr.Wind.Unit, SpeedUnitKnot, "trend[1].wind.unit")
		assertEqual(t, tr.Raw, "BECMG 0909/0911 10010KT", "trend[1].raw")
	}
	// Trend 2: BECMG 0918/0920 16006KT
	{
		tr := result.Trends[2]
		assertEqual(t, tr.Type, WeatherChangeTypeBECMG, "trend[2].type")
		v, ok := tr.Validity.(Validity)
		if !ok {
			t.Error("trend[2].validity not Validity type")
		} else {
			assertEqual(t, v.StartDay, 9, "trend[2].validity.startDay")
			assertEqual(t, v.StartHour, 18, "trend[2].validity.startHour")
			assertEqual(t, v.EndDay, 9, "trend[2].validity.endDay")
			assertEqual(t, v.EndHour, 20, "trend[2].validity.endHour")
		}
		assertEqual(t, len(tr.Clouds), 0, "trend[2].clouds")
		assertNil(t, tr.Probability, "trend[2].prob")
		assertNil(t, tr.Visibility, "trend[2].vis")
		assertEqual(t, len(tr.WeatherConditions), 0, "trend[2].wc")
		assertNotNil(t, tr.Wind, "trend[2].wind")
		assertIntPtr(t, tr.Wind.Degrees, 160, "trend[2].wind.degrees")
		assertEqual(t, tr.Wind.Direction, "SSE", "trend[2].wind.direction")
		assertEqual(t, tr.Wind.Speed, 6, "trend[2].wind.speed")
		assertEqual(t, tr.Wind.Unit, SpeedUnitKnot, "trend[2].wind.unit")
		assertEqual(t, tr.Raw, "BECMG 0918/0920 16006KT", "trend[2].raw")
	}
}
func TestTAFParserPartialStatement(t *testing.T) {
	input := "PART 1 OF 3 TAF SBGL 082150Z 0900/1006 09007KT CAVOK TN21/0909Z TX30/0917Z BECMG 0903/0905 34005KT SCT020 PROB40 0909/0912 4000 BR SCT010 BKN020 BECMG 0912/0914 01005KT FEW023 BECMG 0917/0919 23017KT SCT020 BECMG 0921/0923 27008KT BKN025 BECMG 1004/1006 5000 BR BKN010 RMK PHP"
	_, err := ParseTAF(input, nil)
	if err == nil {
		t.Fatal("expected error for PART x OF y, got nil")
	}
	var partialErr *PartialWeatherStatementError
	if !errors.As(err, &partialErr) {
		t.Fatalf("expected *PartialWeatherStatementError, got %T: %v", err, err)
	}
	assertEqual(t, partialErr.Part, 1, "part")
	assertEqual(t, partialErr.Total, 3, "total")
}

// ============================================================
// TAF parser: second findFlags block (TAF AMD TAF AMD pattern)
// ============================================================

func TestTAFParserFlagsAfterSecondTAF(t *testing.T) {
	code := "TAF AMD TAF AMD KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 QNH2941INS"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertBoolPtr(t, result.Amendment, true, "amendment")
	assertEqual(t, result.Station, "KLSV", "station")
}

// ============================================================
// TAF parser: error paths
// ============================================================

func TestTAFParserEmpty(t *testing.T) {
	_, err := ParseTAF("", nil)
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestTAFParserNoStation(t *testing.T) {
	_, err := ParseTAF("TAF", nil)
	if err == nil {
		t.Fatal("expected error for TAF without station")
	}
}

func TestTAFParserInvalidValidity(t *testing.T) {
	_, err := ParseTAF("TAF KLSV 222300Z INVALID 21020G35KT", nil)
	if err == nil {
		t.Fatal("expected error for invalid validity")
	}
}

// ============================================================
// TAF parser: visibility prefix merge (P + 6SM)
// ============================================================

func TestTAFParserVisPrefixMerge(t *testing.T) {
	code := "TAF KLAX 140520Z 1406/1512 P 6SM FEW010"
	result, err := ParseTAF(code, nil)
	if err != nil {
		t.Fatalf("ParseTAF error: %v", err)
	}
	assertNotNil(t, result.Visibility, "visibility")
	assertNotNil(t, result.Visibility.Indicator, "vis.indicator")
	assertEqual(t, *result.Visibility.Indicator, ValueIndicatorGreaterThan, "vis.indicator")
	assertEqual(t, result.Visibility.Value, 6.0, "vis.value")
	assertEqual(t, result.Visibility.Unit, DistanceUnitStatuteMiles, "vis.unit")
}
