package metartafparser

import (
	"strings"
	"testing"
)

// RemarkParser tests (first block)
// ============================================================

func TestRemarkParserBasic(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("Token AO1 End of remark NXT FCST BY 160300Z")
	if len(remarks) < 4 {
		t.Fatalf("expected at least 4 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeAO1, "remarks[1].type")
	assertEqual(t, remarks[3].Type, RemarkTypeNextForecastBy, "remarks[3].type")
	if remarks[3].Type == RemarkTypeNextForecastBy {
		if remarks[3].Description != nil {
			if !strings.Contains(*remarks[3].Description, "next forecast") {
				t.Errorf("remarks[3].description should contain 'next forecast', got %q", *remarks[3].Description)
			}
		}
		assertIntPtr(t, remarks[3].Day, 16, "remarks[3].day")
		assertIntPtr(t, remarks[3].Hour, 3, "remarks[3].hour")
		assertIntPtr(t, remarks[3].Minute, 0, "remarks[3].minute")
	}
}

// ============================================================
// RemarkParser tests (extensive block)
// ============================================================

func TestRemarkParserAO1(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("Token AO1 End of remark")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeAO1, "type")
}

func TestRemarkParserAO2(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("Token AO2 End of remark")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeAO2, "type")
}

func TestRemarkParserPeakWind(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 PK WND 28045/15")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeWindPeak, "type")
	assertIntPtr(t, remarks[1].Speed, 45, "speed")
	assertIntPtr(t, remarks[1].Degrees, 280, "degrees")
	assertIntPtr(t, remarks[1].StartMinute, 15, "startMinute")
}

func TestRemarkParserPeakWindHour(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 PK WND 28045/1515")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeWindPeak, "type")
	assertIntPtr(t, remarks[1].Speed, 45, "speed")
	assertIntPtr(t, remarks[1].Degrees, 280, "degrees")
	assertIntPtr(t, remarks[1].StartHour, 15, "startHour")
	assertIntPtr(t, remarks[1].StartMinute, 15, "startMinute")
}

func TestRemarkParserWindShift(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 WSHFT 30")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeWindShift, "type")
	assertIntPtr(t, remarks[1].StartMinute, 30, "startMinute")
}

func TestRemarkParserWindShiftHour(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 WSHFT 1530")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeWindShift, "type")
	assertIntPtr(t, remarks[1].StartHour, 15, "startHour")
	assertIntPtr(t, remarks[1].StartMinute, 30, "startMinute")
}

func TestRemarkParserWindShiftFropa(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 WSHFT 1530 FROPA")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeWindShiftFropa, "type")
	assertIntPtr(t, remarks[1].StartHour, 15, "startHour")
	assertIntPtr(t, remarks[1].StartMinute, 30, "startMinute")
}

func TestRemarkParserTowerVisibility(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 TWR VIS 16 1/2")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeTowerVisibility, "type")
	assertFloatPtr(t, remarks[1].Value, 16.5, "value")
}

func TestRemarkParserSurfaceVisibility(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 SFC VIS 16 1/2")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSurfaceVisibility, "type")
	assertFloatPtr(t, remarks[1].Value, 16.5, "value")
}

func TestRemarkParserPrevailingVisibility(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 VIS 1/2V2")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypePrevailingVisibility, "type")
	assertFloatPtr(t, remarks[1].Min, 0.5, "min")
	assertFloatPtr(t, remarks[1].Max, 2, "max")
}

func TestRemarkParserSectorVisibility(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 VIS NE 2 1/2")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSectorVisibility, "type")
	assertFloatPtr(t, remarks[1].Value, 2.5, "value")
	if remarks[1].Direction != nil {
		assertEqual(t, *remarks[1].Direction, "NE", "direction")
	}
}

func TestRemarkParserSecondLocationVisibility(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 VIS 2 1/2 RWY11")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSecondLocationVisibility, "type")
	assertFloatPtr(t, remarks[1].Value, 2.5, "value")
	if remarks[1].Location != nil {
		assertEqual(t, *remarks[1].Location, "RWY11", "location")
	}
}

func TestRemarkParserTornadoBeg(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 TORNADO B13 6 NE")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeTornadicActivityBeg, "type")
	assertIntPtr(t, remarks[1].StartMinute, 13, "startMinute")
	assertFloatPtr(t, remarks[1].Value, 6, "value")
	if remarks[1].Direction != nil {
		assertEqual(t, *remarks[1].Direction, "NE", "direction")
	}
}

func TestRemarkParserTornadoBegHour(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 TORNADO B1513 6 NE")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeTornadicActivityBeg, "type")
	assertIntPtr(t, remarks[1].StartHour, 15, "startHour")
	assertIntPtr(t, remarks[1].StartMinute, 13, "startMinute")
	assertFloatPtr(t, remarks[1].Value, 6, "value")
	if remarks[1].Direction != nil {
		assertEqual(t, *remarks[1].Direction, "NE", "direction")
	}
}

func TestRemarkParserFunnelCloudBegEnd(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 FUNNEL CLOUD B1513E1630 6 NE")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeTornadicActivityBegEnd, "type")
	assertIntPtr(t, remarks[1].StartHour, 15, "startHour")
	assertIntPtr(t, remarks[1].StartMinute, 13, "startMinute")
	assertIntPtr(t, remarks[1].EndHour, 16, "endHour")
	assertIntPtr(t, remarks[1].EndMinute, 30, "endMinute")
	assertFloatPtr(t, remarks[1].Value, 6, "value")
	if remarks[1].Direction != nil {
		assertEqual(t, *remarks[1].Direction, "NE", "direction")
	}
}

func TestRemarkParserWaterspoutEnd(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 WATERSPOUT E16 12 NE")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeTornadicActivityEnd, "type")
	assertIntPtr(t, remarks[1].EndMinute, 16, "endMinute")
	assertFloatPtr(t, remarks[1].Value, 12, "value")
	if remarks[1].Direction != nil {
		assertEqual(t, *remarks[1].Direction, "NE", "direction")
	}
}

func TestRemarkParserWaterspoutEndHour(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 WATERSPOUT E1516 12 NE")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeTornadicActivityEnd, "type")
	assertIntPtr(t, remarks[1].EndHour, 15, "endHour")
	assertIntPtr(t, remarks[1].EndMinute, 16, "endMinute")
	assertFloatPtr(t, remarks[1].Value, 12, "value")
	if remarks[1].Direction != nil {
		assertEqual(t, *remarks[1].Direction, "NE", "direction")
	}
}

func TestRemarkParserPrecipBegEnd(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 RAB05E30SNB1520E1655")
	if len(remarks) < 3 {
		t.Fatalf("expected at least 3 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypePrecipitationBegEnd, "type")
	if remarks[1].Phenomenon != nil {
		assertEqual(t, *remarks[1].Phenomenon, PhenomenonRain, "remarks[1].phenomenon")
	}
	assertIntPtr(t, remarks[1].StartMinute, 5, "remarks[1].startMinute")
	assertIntPtr(t, remarks[1].EndMinute, 30, "remarks[1].endMinute")

	assertEqual(t, remarks[2].Type, RemarkTypePrecipitationBegEnd, "remarks[2].type")
	if remarks[2].Phenomenon != nil {
		assertEqual(t, *remarks[2].Phenomenon, PhenomenonSnow, "remarks[2].phenomenon")
	}
	assertIntPtr(t, remarks[2].StartHour, 15, "remarks[2].startHour")
	assertIntPtr(t, remarks[2].StartMinute, 20, "remarks[2].startMinute")
	assertIntPtr(t, remarks[2].EndHour, 16, "remarks[2].endHour")
	assertIntPtr(t, remarks[2].EndMinute, 55, "remarks[2].endMinute")
}

func TestRemarkParserPrecipBegEndDescriptive(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 SHRAB05E30SHSNB20E55")
	if len(remarks) < 3 {
		t.Fatalf("expected at least 3 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypePrecipitationBegEnd, "type")
	if remarks[1].Descriptive != nil {
		assertEqual(t, *remarks[1].Descriptive, DescriptiveShowers, "remarks[1].descriptive")
	}
	if remarks[1].Phenomenon != nil {
		assertEqual(t, *remarks[1].Phenomenon, PhenomenonRain, "remarks[1].phenomenon")
	}
	assertEqual(t, remarks[2].Type, RemarkTypePrecipitationBegEnd, "remarks[2].type")
	if remarks[2].Descriptive != nil {
		assertEqual(t, *remarks[2].Descriptive, DescriptiveShowers, "remarks[2].descriptive")
	}
	if remarks[2].Phenomenon != nil {
		assertEqual(t, *remarks[2].Phenomenon, PhenomenonSnow, "remarks[2].phenomenon")
	}
}

func TestRemarkParserPrecipBeg(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 SHRAB05SHSNB0220")
	if len(remarks) < 3 {
		t.Fatalf("expected at least 3 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypePrecipitationBeg, "remarks[1].type")
	if remarks[1].Descriptive != nil {
		assertEqual(t, *remarks[1].Descriptive, DescriptiveShowers, "remarks[1].descriptive")
	}
	if remarks[1].Phenomenon != nil {
		assertEqual(t, *remarks[1].Phenomenon, PhenomenonRain, "remarks[1].phenomenon")
	}
	assertIntPtr(t, remarks[1].StartMinute, 5, "remarks[1].startMinute")
	assertEqual(t, remarks[2].Type, RemarkTypePrecipitationBeg, "remarks[2].type")
	if remarks[2].Descriptive != nil {
		assertEqual(t, *remarks[2].Descriptive, DescriptiveShowers, "remarks[2].descriptive")
	}
	if remarks[2].Phenomenon != nil {
		assertEqual(t, *remarks[2].Phenomenon, PhenomenonSnow, "remarks[2].phenomenon")
	}
	assertIntPtr(t, remarks[2].StartHour, 2, "remarks[2].startHour")
	assertIntPtr(t, remarks[2].StartMinute, 20, "remarks[2].startMinute")
}

func TestRemarkParserPrecipEnd(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 SHRAE05SHSNE0120")
	if len(remarks) < 3 {
		t.Fatalf("expected at least 3 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypePrecipitationEnd, "remarks[1].type")
	if remarks[1].Descriptive != nil {
		assertEqual(t, *remarks[1].Descriptive, DescriptiveShowers, "remarks[1].descriptive")
	}
	if remarks[1].Phenomenon != nil {
		assertEqual(t, *remarks[1].Phenomenon, PhenomenonRain, "remarks[1].phenomenon")
	}
	assertIntPtr(t, remarks[1].EndMinute, 5, "remarks[1].endMinute")
	assertEqual(t, remarks[2].Type, RemarkTypePrecipitationEnd, "remarks[2].type")
	if remarks[2].Descriptive != nil {
		assertEqual(t, *remarks[2].Descriptive, DescriptiveShowers, "remarks[2].descriptive")
	}
	if remarks[2].Phenomenon != nil {
		assertEqual(t, *remarks[2].Phenomenon, PhenomenonSnow, "remarks[2].phenomenon")
	}
	assertIntPtr(t, remarks[2].EndHour, 1, "remarks[2].endHour")
	assertIntPtr(t, remarks[2].EndMinute, 20, "remarks[2].endMinute")
}

func TestRemarkParserThunderstormBegEnd(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 TSB0159E30")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypePrecipitationBegEnd, "type")
	if remarks[1].Phenomenon != nil {
		assertEqual(t, *remarks[1].Phenomenon, PhenomenonThunderstorm, "phenomenon")
	}
	assertIntPtr(t, remarks[1].StartHour, 1, "startHour")
	assertIntPtr(t, remarks[1].StartMinute, 59, "startMinute")
	assertIntPtr(t, remarks[1].EndMinute, 30, "endMinute")
}

func TestRemarkParserThunderstormLocation(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 TS SE")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeThunderStormLocation, "type")
	if remarks[1].Location != nil {
		assertEqual(t, *remarks[1].Location, "SE", "location")
	}
}

func TestRemarkParserThunderstormLocationMoving(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 TS SE MOV NE")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeThunderStormLocationMoving, "type")
	if remarks[1].Location != nil {
		assertEqual(t, *remarks[1].Location, "SE", "location")
	}
	if remarks[1].Moving != nil {
		assertEqual(t, *remarks[1].Moving, "NE", "moving")
	}
}

func TestRemarkParserHailSize(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 GR 1 3/4")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeHailSize, "type")
	assertFloatPtr(t, remarks[1].Value, 1.75, "value")
}

func TestRemarkParserSmallHailSize(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 GR LESS THAN 1/4")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSmallHailSize, "type")
	assertFloatPtr(t, remarks[1].Value, 0.25, "value")
}

func TestRemarkParserSnowPellets(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 GS MOD")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSnowPellets, "type")
}

func TestRemarkParserVirgaDirection(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 VIRGA SW")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeVirgaDirection, "type")
	if remarks[1].Direction != nil {
		assertEqual(t, *remarks[1].Direction, "SW", "direction")
	}
}

func TestRemarkParserVirga(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 VIRGA")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeVIRGA, "type")
}

func TestRemarkParserCeilingHeight(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 CIG 005V010")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeCeilingHeight, "type")
	assertFloatPtr(t, remarks[1].Min, 500, "min")
	assertFloatPtr(t, remarks[1].Max, 1000, "max")
}

func TestRemarkParserObscuration(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 FU BKN020")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeObscuration, "type")
	assertFloatPtr(t, remarks[1].Min, 2000, "height")
}

func TestRemarkParserVariableSky(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("BKN V OVC")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeVariableSky, "type")
}

func TestRemarkParserVariableSkyHeight(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("BKN014 V OVC")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeVariableSkyHeight, "type")
	assertFloatPtr(t, remarks[0].Min, 1400, "height")
}

func TestRemarkParserCeilingSecondLocation(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("CIG 002 RWY11")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeCeilingSecondLocation, "type")
	assertFloatPtr(t, remarks[0].Min, 200, "height")
	if remarks[0].Location != nil {
		assertEqual(t, *remarks[0].Location, "RWY11", "location")
	}
}

func TestRemarkParserSeaLevelPressure(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 SLP134")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSeaLevelPressure, "type")
	assertFloatPtr(t, remarks[1].Value, 1013.4, "value")
}

func TestRemarkParserSnowIncrease(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 SNINCR 2/10")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSnowIncrease, "type")
	assertIntPtr(t, remarks[1].InchesLastHour, 2, "inchesLastHour")
	assertIntPtr(t, remarks[1].TotalDepth, 10, "totalDepth")
}

func TestRemarkParserHourlyMaxMinTemp(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 401020020")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeHourlyMaximumMinimumTemperature, "type")
	if remarks[1].Max != nil {
		assertEqual(t, *remarks[1].Max, 10.2, "max")
	}
	if remarks[1].Min != nil {
		assertEqual(t, *remarks[1].Min, 2.0, "min")
	}
}

func TestRemarkParserHourlyMaxMinTempTSVariant(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("401001015")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyMaximumMinimumTemperature, "type")
	assertFloatPtr(t, remarks[0].Max, 10.0, "max")
	assertFloatPtr(t, remarks[0].Min, -1.5, "min")
}

func TestRemarkParserHourlyTempDewPoint(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 T10171017")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeHourlyTemperatureDewPoint, "type")
	assertEqual(t, remarks[1].Raw, "T10171017", "raw")
	assertFloatPtr(t, remarks[1].Temperature, -1.7, "temperature")
	assertFloatPtr(t, remarks[1].DewPoint, -1.7, "dewPoint")
}

func TestRemarkParserNextForecastBy(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("NXT FCST BY 160300Z")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeNextForecastBy, "type")
	assertIntPtr(t, remarks[0].Day, 16, "day")
	assertIntPtr(t, remarks[0].Hour, 3, "hour")
	assertIntPtr(t, remarks[0].Minute, 0, "minute")
}

// ---- Missing RemarkParser tests ----

func TestRemarkParserWindShiftFropaHour(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 WSHFT 30 FROPA")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeWindShiftFropa, "type")
	assertIntPtr(t, remarks[1].StartMinute, 30, "startMinute")
}

func TestRemarkParserSeaLevelPressureLower(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("AO1 SLP982")
	if len(remarks) < 2 {
		t.Fatalf("expected at least 2 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSeaLevelPressure, "type")
	assertFloatPtr(t, remarks[1].Value, 998.2, "value")
}

func TestRemarkParserRmkSLP(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("CF1AC8 CF TR SLP091 DENSITY ALT 200FT")
	if len(remarks) < 3 {
		t.Fatalf("expected at least 3 remarks, got %d", len(remarks))
	}
	assertEqual(t, remarks[1].Type, RemarkTypeSeaLevelPressure, "type")
	assertEqual(t, remarks[1].Raw, "SLP091", "raw")
	assertFloatPtr(t, remarks[1].Value, 1009.1, "value")
}

func TestRemarkParserHourlyMaxBelowZero(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("11021 AO1")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyMaximumTemperature, "type")
	assertFloatPtr(t, remarks[0].Max, -2.1, "max")
}

func TestRemarkParserHourlyMaxAboveZero(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("10142")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyMaximumTemperature, "type")
	assertFloatPtr(t, remarks[0].Max, 14.2, "max")
}

func TestRemarkParserHourlyMinNegative(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("21001")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyMinimumTemperature, "type")
	assertFloatPtr(t, remarks[0].Min, -0.1, "min")
}

func TestRemarkParserHourlyMinPositive(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("20012")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyMinimumTemperature, "type")
	assertFloatPtr(t, remarks[0].Min, 1.2, "min")
}

func TestRemarkParserHourlyPressure(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("52032")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyPressure, "type")
}

func TestRemarkParserPrecipAmount24Hour(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("70125")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypePrecipitationAmount24Hour, "type")
	assertFloatPtr(t, remarks[0].Amount, 1.25, "amount")
}

func TestRemarkParserSnowDepth(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("4/021")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeSnowDepth, "type")
	assertFloatPtr(t, remarks[0].Value, 21, "depth")
}

func TestRemarkParserSunshineDuration(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("98096")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeSunshineDuration, "type")
	assertFloatPtr(t, remarks[0].Min, 96, "duration")
}

func TestRemarkParserWaterEquivalentSnow(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("933036")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeWaterEquivalentSnow, "type")
	assertFloatPtr(t, remarks[0].Amount, 3.6, "amount")
}

func TestRemarkParserIceAccretion(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("l1004 AO1")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeIceAccretion, "type")
	assertFloatPtr(t, remarks[0].Amount, 0.04, "amount")
}

func TestRemarkParserHourlyTemperature(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("T0026 AO1")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyTemperatureDewPoint, "type")
	assertFloatPtr(t, remarks[0].Temperature, 2.6, "temperature")
	assertNil(t, remarks[0].DewPoint, "dewPoint")
}

func TestRemarkParserHourlyTemperatureDewPoint(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("T00261015 AO1")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyTemperatureDewPoint, "type")
	assertEqual(t, remarks[0].Raw, "T00261015", "raw")
	assertFloatPtr(t, remarks[0].Temperature, 2.6, "temperature")
	assertFloatPtr(t, remarks[0].DewPoint, -1.5, "dewPoint")
}

func TestRemarkParserPrecipAmount3Hours(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("30217")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypePrecipitationAmount36Hour, "type")
	assertFloatPtr(t, remarks[0].Amount, 2.17, "amount")
}

func TestRemarkParserPrecipAmount6Hours(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("60217")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypePrecipitationAmount36Hour, "type")
	assertFloatPtr(t, remarks[0].Amount, 2.17, "amount")
}

// ============================================================
// HourlyPrecipitationAmountCommand test
// ============================================================

func TestRemarkParserHourlyPrecipitationAmount(t *testing.T) {
	parser := newRemarkParser(DefaultLocale())
	remarks := parser.Parse("P0002")
	if len(remarks) < 1 {
		t.Fatalf("expected at least 1 remark, got %d", len(remarks))
	}
	assertEqual(t, remarks[0].Type, RemarkTypeHourlyPrecipitationAmount, "type")
	assertFloatPtr(t, remarks[0].Amount, 0.02, "amount")
	if remarks[0].Description != nil {
		if !strings.Contains(*remarks[0].Description, "0.02") {
			t.Errorf("expected description to contain amount, got %q", *remarks[0].Description)
		}
	}
}
