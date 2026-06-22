package main

import (
	"testing"

	metartafparser "github.com/ryansavara/metar-taf-parser"
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

func TestParseAuto_TAFPrefixFallbackToMetar(t *testing.T) {
	result, err := parseAuto("TAF LFPG 170830Z 00000KT 0350 FG SCT000 M01/M01 Q1026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.(*metartafparser.Metar); !ok {
		t.Errorf("expected *Metar, got %T", result)
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
