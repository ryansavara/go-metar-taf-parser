// Package metartafparser parses METAR (Meteorological Aerodrome Report) and TAF
// (Terminal Aerodrome Forecast) aviation weather reports.
//
// METAR and TAF are standardized formats defined by the World Meteorological
// Organization (WMO) and the International Civil Aviation Organization (ICAO)
// for reporting and forecasting airport weather conditions.
//
// Basic usage:
//
//	// Parse a METAR
//	metar, err := metartafparser.ParseMetar("KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992 RMK AO2 SLP132 T01440117", nil)
//
//	// Parse a TAF
//	taf, err := metartafparser.ParseTAF("KLAX 140520Z 1406/1512 05005KT P6SM FEW010", nil)
//
//	// Parse with date hydration (converts relative day/hour to absolute timestamps)
//	issued := time.Date(2024, 6, 14, 8, 53, 0, 0, time.UTC)
//	metarDated, err := metartafparser.ParseMetarDated("KLAX 140853Z ...", issued, nil)
//
//	// Parse TAF as a forecast container
//	fc, err := metartafparser.ParseTAFAsForecast("KLAX 140520Z ...", issued, nil)
//
//	// Get composite forecast for a specific time
//	cf, err := metartafparser.GetCompositeForecastForDate(someTime, fc)
//
//	// Use a custom locale for translated output
//	opts := &metartafparser.ParseOptions{Locale: metartafparser.DefaultLocale()}
//	metar, err := metartafparser.ParseMetar("...", opts)
package metartafparser
