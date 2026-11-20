package main

import (
	"strings"
	"testing"

	metartafparser "github.com/ryansavara/go-metar-taf-parser"
)

func TestParseAuto_METAR(t *testing.T) {
	result, err := parseAuto("KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*metartafparser.Metar); !ok {
		t.Errorf("expected *Metar, got %T", result)
	}
}

func TestParseAuto_TAF(t *testing.T) {
	result, err := parseAuto("TAF KLAX 140520Z 1406/1512 05005KT P6SM FEW010")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*metartafparser.TAF); !ok {
		t.Errorf("expected *TAF, got %T", result)
	}
}

func TestParseAuto_RejectsLowercaseInput(t *testing.T) {
	cases := []string{
		"taf KJFK 121720Z 1218/1324 27008KT P6SM FEW020",
		"Taf KJFK 121720Z 1218/1324 27008KT P6SM FEW020",
		"kjfk 121720Z 27008KT P6SM FEW020",
	}

	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			_, err := parseAuto(input)
			if err == nil {
				t.Fatal("expected error for lowercase input")
			}
		})
	}
}

func TestParseAuto_InvalidTAFDoesNotFallbackToMetar(t *testing.T) {
	_, err := parseAuto("TAF LFPG 170830Z 00000KT 0350 FG SCT000 M01/M01 Q1026")
	if err == nil {
		t.Fatal("expected error for invalid TAF")
	}
}

func TestParseAuto_MetarStationStartingWithTAF(t *testing.T) {
	result, err := parseAuto("TAFB 121720Z 18005KT 10SM FEW020 10/08 A2992")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	metar, ok := result.(*metartafparser.Metar)
	if !ok {
		t.Fatalf("expected *Metar, got %T", result)
	}
	if metar.Station != "TAFB" {
		t.Errorf("expected station TAFB, got %s", metar.Station)
	}
}

func TestParseAuto_PrefixlessTAF(t *testing.T) {
	result, err := parseAuto("KAAA 121720Z 1218/1324 18005KT P6SM FEW020")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	taf, ok := result.(*metartafparser.TAF)
	if !ok {
		t.Fatalf("expected *TAF, got %T", result)
	}
	if taf.Station != "KAAA" {
		t.Errorf("expected station KAAA, got %s", taf.Station)
	}
}

func TestParseAuto_TAFWithAMD(t *testing.T) {
	result, err := parseAuto("TAF AMD KJFK 121720Z 1218/1324 27008KT P6SM FEW020")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	taf, ok := result.(*metartafparser.TAF)
	if !ok {
		t.Fatalf("expected *TAF, got %T", result)
	}
	if taf.Station != "KJFK" {
		t.Errorf("expected station KJFK, got %s", taf.Station)
	}
}

func TestParseAuto_TAFWithCOR(t *testing.T) {
	result, err := parseAuto("TAF COR KJFK 121720Z 1218/1324 27008KT P6SM FEW020")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	taf, ok := result.(*metartafparser.TAF)
	if !ok {
		t.Fatalf("expected *TAF, got %T", result)
	}
	if taf.Station != "KJFK" {
		t.Errorf("expected station KJFK, got %s", taf.Station)
	}
}

func TestParseAuto_SPECI(t *testing.T) {
	result, err := parseAuto("SPECI KJFK 121720Z 27008KT P6SM FEW020")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	metar, ok := result.(*metartafparser.Metar)
	if !ok {
		t.Fatalf("expected *Metar, got %T", result)
	}
	if metar.Station != "KJFK" {
		t.Errorf("expected station KJFK, got %s", metar.Station)
	}
	if metar.Type == nil || *metar.Type != metartafparser.MetarTypeSPECI {
		t.Errorf("expected type SPECI, got %v", metar.Type)
	}
}

func TestParseAuto_Invalid(t *testing.T) {
	_, err := parseAuto("BLAH BLAH BLAH")
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
}

func TestParseAuto_Empty(t *testing.T) {
	_, err := parseAuto("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestIsMetarValid_nil(t *testing.T) {
	if isMetarValid(nil) {
		t.Error("expected false for nil")
	}
}

func TestIsMetarValid_wind(t *testing.T) {
	if !isMetarValid(&metartafparser.Metar{Container: metartafparser.Container{Wind: &metartafparser.Wind{}}}) {
		t.Error("expected true with wind")
	}
}

func TestIsMetarValid_visibility(t *testing.T) {
	if !isMetarValid(&metartafparser.Metar{Container: metartafparser.Container{Visibility: &metartafparser.Visibility{}}}) {
		t.Error("expected true with visibility")
	}
}

func TestIsMetarValid_clouds(t *testing.T) {
	if !isMetarValid(&metartafparser.Metar{Container: metartafparser.Container{Clouds: []metartafparser.Cloud{{}}}}) {
		t.Error("expected true with clouds")
	}
}

func TestIsMetarValid_weather(t *testing.T) {
	if !isMetarValid(&metartafparser.Metar{Container: metartafparser.Container{WeatherConditions: []metartafparser.WeatherCondition{{Phenomena: []metartafparser.Phenomenon{metartafparser.PhenomenonRain}}}}}) {
		t.Error("expected true with weather")
	}
}

func TestIsMetarValid_empty(t *testing.T) {
	if isMetarValid(&metartafparser.Metar{}) {
		t.Error("expected false for empty Metar")
	}
}

func TestParseAuto_NilMetar(t *testing.T) {
	result, err := parseAuto("METAR KJFK 121751Z NIL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	metar, ok := result.(*metartafparser.Metar)
	if !ok {
		t.Fatalf("expected *Metar, got %T", result)
	}
	if metar.Nil == nil || !*metar.Nil {
		t.Error("expected Metar.Nil to be true")
	}
}

func TestParseAuto_NilMetarWithoutPrefix(t *testing.T) {
	result, err := parseAuto("SVMC 211703Z AUTO NIL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	metar, ok := result.(*metartafparser.Metar)
	if !ok {
		t.Fatalf("expected *Metar, got %T", result)
	}
	if metar.Nil == nil || !*metar.Nil {
		t.Error("expected Metar.Nil to be true")
	}
}

func TestParseAuto_StationAndTimeOnlyStillInvalid(t *testing.T) {
	_, err := parseAuto("KJFK 121751Z")
	if err == nil {
		t.Fatal("expected error for station and time only")
	}
}

func TestIsMetarValid_nilFlag(t *testing.T) {
	if !isMetarValid(&metartafparser.Metar{Nil: boolPtr(true)}) {
		t.Error("expected true for Metar with Nil=true")
	}
}

func TestIsMetarValid_canceledFlag(t *testing.T) {
	if !isMetarValid(&metartafparser.Metar{Canceled: boolPtr(true)}) {
		t.Error("expected true for Metar with Canceled=true")
	}
}

func boolPtr(b bool) *bool {
	return &b
}

// FuzzParseAuto verifies that CLI auto-detection rejects reports containing
// lowercase characters.
func FuzzParseAuto(f *testing.F) {
	f.Add("KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992")
	f.Add("TAF KLAX 140520Z 1406/1512 05005KT P6SM FEW010")
	f.Add("taf KJFK 121720Z 1218/1324 27008KT P6SM FEW020")
	f.Add("kjfk 121720Z 27008KT P6SM FEW020")

	f.Fuzz(func(t *testing.T, input string) {
		trimmed := strings.TrimSpace(input)
		_, err := parseAuto(input)
		if trimmed != "" && trimmed != strings.ToUpper(trimmed) && err == nil {
			t.Errorf("parseAuto accepted input containing lowercase characters: %q", input)
		}
	})
}
