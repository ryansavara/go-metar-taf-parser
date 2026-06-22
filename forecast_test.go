package metartafparser

import (
	"errors"
	"testing"
	"time"
)

// ============================================================
// Forecast Tests
// ============================================================

func TestGetForecastFromTAFSimple(t *testing.T) {
	taf, err := ParseTAFDated("TAF KMSN 142325Z 1500/1524", time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	assertEqual(t, fc.Start, time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC).UnixMilli(), "start")
	assertEqual(t, fc.End, time.Date(2022, 4, 16, 0, 0, 0, 0, time.UTC).UnixMilli(), "end")
	assertEqual(t, len(fc.Forecast), 1, "forecast length")
}

func TestGetForecastFromTAFWithTrends(t *testing.T) {
	code := "TAF KMSN 142325Z 1500/1524 25014G30KT P6SM VCSH SCT035 BKN070 TEMPO 1500/1502 6SM -SHRASN BKN035 FM150100 25012G25KT P6SM VCSH SCT040 BKN070 FM150300 26011G21KT P6SM SCT080"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	assertEqual(t, len(fc.Forecast), 4, "forecast length")

	// Initial forecast
	f0 := fc.Forecast[0]
	assertEqual(t, f0.Start, time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC).UnixMilli(), "f0.start")
	assertEqual(t, f0.End, time.Date(2022, 4, 15, 1, 0, 0, 0, time.UTC).UnixMilli(), "f0.end")
	if f0.Type != nil {
		t.Errorf("f0.type should be nil (initial), got %v", *f0.Type)
	}
	assertEqual(t, f0.Raw, "TAF KMSN 142325Z 1500/1524 25014G30KT P6SM VCSH SCT035 BKN070", "f0.raw")

	// TEMPO (supplemental - keeps explicit end)
	f1 := fc.Forecast[1]
	if f1.Type == nil || *f1.Type != WeatherChangeTypeTEMPO {
		t.Errorf("f1.type should be TEMPO")
	}
	assertEqual(t, f1.Start, time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC).UnixMilli(), "f1.start")
	assertEqual(t, f1.End, time.Date(2022, 4, 15, 2, 0, 0, 0, time.UTC).UnixMilli(), "f1.end")
	assertEqual(t, f1.Raw, "TEMPO 1500/1502 6SM -SHRASN BKN035", "f1.raw")

	// FM
	f2 := fc.Forecast[2]
	if f2.Type == nil || *f2.Type != WeatherChangeTypeFM {
		t.Errorf("f2.type should be FM")
	}
	assertEqual(t, f2.Start, time.Date(2022, 4, 15, 1, 0, 0, 0, time.UTC).UnixMilli(), "f2.start")
	assertEqual(t, f2.End, time.Date(2022, 4, 15, 3, 0, 0, 0, time.UTC).UnixMilli(), "f2.end")

	// Second FM
	f3 := fc.Forecast[3]
	assertEqual(t, f3.Start, time.Date(2022, 4, 15, 3, 0, 0, 0, time.UTC).UnixMilli(), "f3.start")
	assertEqual(t, f3.End, time.Date(2022, 4, 16, 0, 0, 0, 0, time.UTC).UnixMilli(), "f3.end")
}

func TestGetForecastFromTAFBECMGWithBy(t *testing.T) {
	code := "TAF ESGG 260830Z 2609/2709 02009KT 3000 BR BKN003 BECMG 2609/2611 9999 NSW FEW015"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	assertEqual(t, len(fc.Forecast), 2, "forecast length")

	// BECMG has by
	f1 := fc.Forecast[1]
	if f1.Type == nil || *f1.Type != WeatherChangeTypeBECMG {
		t.Errorf("f1.type should be BECMG")
	}
	if f1.By == nil {
		t.Errorf("f1.by should be set for BECMG")
	} else {
		assertEqual(t, *f1.By, time.Date(2022, 4, 26, 11, 0, 0, 0, time.UTC).UnixMilli(), "f1.by")
	}
}

func TestGetForecastFromTAFBECMGContextInheritance(t *testing.T) {
	code := "TAF SBPJ 221450Z 2218/2318 21006KT 8000 SCT030 FEW040TCU TN25/2309Z TX34/2316Z BECMG 2221/2223 VRB03KT FEW030 BECMG 2302/2304 16003KT 5000 FU RMK PGU BECMG 2313/2315 23005KT SCT030 FEW040TCU"
	taf, err := ParseTAFDated(code, time.Date(2022, 10, 22, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	assertEqual(t, len(fc.Forecast), 4, "forecast length")

	f0 := fc.Forecast[0]
	assertEqual(t, len(f0.Clouds), 2, "f0.clouds")
	assertEqual(t, len(f0.WeatherConditions), 0, "f0.wc")
	if f0.Visibility != nil {
		assertEqual(t, f0.Visibility.Value, float64(8000), "f0.vis")
	}

	// BECMG 0: only wind changes, inherits vis/clouds from initial
	f1 := fc.Forecast[1]
	assertEqual(t, f1.Start, time.Date(2022, 10, 22, 21, 0, 0, 0, time.UTC).UnixMilli(), "f1.start")
	assertEqual(t, f1.End, time.Date(2022, 10, 23, 2, 0, 0, 0, time.UTC).UnixMilli(), "f1.end")
	assertEqual(t, len(f1.Clouds), 1, "f1.clouds")
	if f1.Visibility != nil {
		assertEqual(t, f1.Visibility.Value, float64(8000), "f1.vis (inherited)")
	}

	// BECMG 1: vis changes, inherits clouds
	f2 := fc.Forecast[2]
	assertEqual(t, len(f2.Clouds), 1, "f2.clouds")
	assertEqual(t, len(f2.WeatherConditions), 1, "f2.wc")
	assertEqual(t, len(f2.Remarks), 1, "f2.remarks")
	if f2.Visibility != nil {
		assertEqual(t, f2.Visibility.Value, float64(5000), "f2.vis")
	}

	// BECMG 2: clouds change, inherits vis
	f3 := fc.Forecast[3]
	assertEqual(t, len(f3.Clouds), 2, "f3.clouds")
	assertEqual(t, len(f3.WeatherConditions), 1, "f3.wc")
	assertEqual(t, len(f3.Remarks), 0, "f3.remarks")
	if f3.Visibility != nil {
		assertEqual(t, f3.Visibility.Value, float64(5000), "f3.vis")
	}
}

func TestGetForecastFromTAFCAVOKPropagation(t *testing.T) {
	code := "TAF TAF EHAM 311720Z 3118/0124 12012KT CAVOK BECMG 3121/3124 17017KT"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 15, 21, 46, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	// CAVOK should propagate to BECMG since only wind changes
	if fc.Forecast[0].Cavok == nil || !*fc.Forecast[0].Cavok {
		t.Error("f0.cavok should be true")
	}
	if fc.Forecast[1].Cavok == nil || !*fc.Forecast[1].Cavok {
		t.Error("f1.cavok should be true (inherited)")
	}
}

func TestGetForecastFromTAFCAVOKNotPropagated(t *testing.T) {
	code := "TAF TAF EHAM 311720Z 3118/0124 12012KT CAVOK TEMPO 3123/0102 5000 RADZ SCT014 BKN018"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 15, 21, 46, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	// TEMPO has explicit weather, should NOT have CAVOK
	assertEqual(t, len(fc.Forecast), 2, "forecast length")
	if fc.Forecast[0].Cavok == nil || !*fc.Forecast[0].Cavok {
		t.Error("f0.cavok should be true (initial)")
	}
	if fc.Forecast[1].Cavok != nil && *fc.Forecast[1].Cavok {
		t.Error("f1.cavok should NOT be true (TEMPO has its own weather)")
	}
}

func TestGetForecastFromTAFMaxMinTemperatures(t *testing.T) {
	code := "TAF AMD SBPJ 221450Z 2218/2318 TN25/2309Z TX34/2316Z"
	taf, err := ParseTAFDated(code, time.Date(2022, 10, 22, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}

	// Check temperature dates on the dated TAF
	if taf.MaxTemperature == nil {
		t.Fatal("expected maxTemperature")
	}
	assertEqual(t, taf.MaxTemperature.Temperature, float64(34), "maxTemp.Temperature")
	assertEqual(t, taf.MaxTemperature.Day, 23, "maxTemp.Day")
	assertEqual(t, taf.MaxTemperature.Hour, 16, "maxTemp.Hour")

	if taf.MinTemperature == nil {
		t.Fatal("expected minTemperature")
	}
	assertEqual(t, taf.MinTemperature.Temperature, float64(25), "minTemp.Temperature")

	fc := getForecastFromTAF(taf)
	if fc.MaxTemperature != nil {
		assertEqual(t, fc.MaxTemperature.Date, time.Date(2022, 10, 23, 16, 0, 0, 0, time.UTC).UnixMilli(), "maxTemp.date")
	}
	if fc.MinTemperature != nil {
		assertEqual(t, fc.MinTemperature.Date, time.Date(2022, 10, 23, 9, 0, 0, 0, time.UTC).UnixMilli(), "minTemp.date")
	}
	if fc.Amendment == nil || !*fc.Amendment {
		t.Error("amendment should be true")
	}
}

func TestGetForecastFromTAFMonthRollover(t *testing.T) {
	code := "TAF KMSN 302325Z 0100/0124"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	assertEqual(t, fc.Issued, time.Date(2022, 4, 30, 23, 25, 0, 0, time.UTC).UnixMilli(), "issued")
	assertEqual(t, fc.Start, time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC).UnixMilli(), "start")
	assertEqual(t, fc.End, time.Date(2022, 5, 2, 0, 0, 0, 0, time.UTC).UnixMilli(), "end")
}

// ============================================================
// GetCompositeForecastForDate Tests
// ============================================================

func TestGetCompositeForecastForDateFindsTEMPO(t *testing.T) {
	code := "TAF KMSN 142325Z 1500/1524 25014G30KT P6SM VCSH SCT035 BKN070 TEMPO 1500/1502 6SM -SHRASN BKN035 FM150100 25012G25KT P6SM VCSH SCT040 BKN070 FM150300 26011G21KT P6SM SCT080"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	composite, err := GetCompositeForecastForDate(time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC), fc)
	if err != nil {
		t.Fatalf("GetCompositeForecastForDate error: %v", err)
	}
	assertEqual(t, len(composite.Supplemental), 1, "supplemental count")
	assertEqual(t, composite.Prevailing.Start, fc.Forecast[0].Start, "prevailing.start")
}

func TestGetCompositeForecastForDateFindsFM(t *testing.T) {
	code := "TAF KMSN 142325Z 1500/1524 25014G30KT P6SM VCSH SCT035 BKN070 TEMPO 1500/1501 6SM -SHRASN BKN035 FM150100 25012G25KT P6SM VCSH SCT040 BKN070 FM150300 26011G21KT P6SM SCT080"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	// At 01:00, FM should be prevailing, no TEMPO (exclusive end)
	composite, err := GetCompositeForecastForDate(time.Date(2022, 4, 15, 1, 0, 0, 0, time.UTC), fc)
	if err != nil {
		t.Fatalf("GetCompositeForecastForDate error: %v", err)
	}
	assertEqual(t, len(composite.Supplemental), 0, "supplemental count")
	if composite.Prevailing.Type == nil || *composite.Prevailing.Type != WeatherChangeTypeFM {
		t.Error("prevailing should be FM")
	}
}

func TestGetCompositeForecastForDateOutOfBounds(t *testing.T) {
	code := "TAF KMSN 142325Z 1500/1524 25014G30KT P6SM VCSH SCT035 BKN070 TEMPO 1500/1502 6SM -SHRASN BKN035 FM150100 25012G25KT P6SM VCSH SCT040 BKN070 FM150300 26011G21KT P6SM SCT080"
	taf, err := ParseTAFDated(code, time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	fc := getForecastFromTAF(taf)

	// Before start
	_, err = GetCompositeForecastForDate(time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC), fc)
	var timestampErr *TimestampOutOfBoundsError
	if !errors.As(err, &timestampErr) {
		t.Errorf("expected TimestampOutOfBoundsError, got %T: %v", err, err)
	}

	// After end
	_, err = GetCompositeForecastForDate(time.Date(2022, 4, 16, 0, 0, 0, 0, time.UTC), fc)
	if !errors.As(err, &timestampErr) {
		t.Errorf("expected TimestampOutOfBoundsError, got %T: %v", err, err)
	}

	// Inclusive start should NOT throw
	_, err = GetCompositeForecastForDate(time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC), fc)
	if err != nil {
		t.Errorf("inclusive start should not throw: %v", err)
	}
}

// ============================================================
// Dated API Tests
// ============================================================

func TestParseMetarDated(t *testing.T) {
	result, err := ParseMetarDated("LFPG 161430Z 24015G25KT 5000 1100w", time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseMetarDated error: %v", err)
	}
	assertEqual(t, result.Station, "LFPG", "station")
	assertEqual(t, result.Issued, time.Date(2022, 1, 16, 14, 30, 0, 0, time.UTC).UnixMilli(), "issued")
}

func TestParseTAFDated(t *testing.T) {
	code := "TAF AMD KMSN 152044Z 1521/1618 24009G16KT P6SM SCT100 TEMPO 1521/1523 BKN100 FM160000 27008KT P6SM SCT150 FM160300 30006KT P6SM FEW230 FM161700 30011G18KT P6SM BKN050"
	result, err := ParseTAFDated(code, time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}

	assertEqual(t, result.Station, "KMSN", "station")
	assertEqual(t, result.Issued, time.Date(2022, 1, 15, 20, 44, 0, 0, time.UTC).UnixMilli(), "issued")
	assertEqual(t, result.ValidityDated.Start, time.Date(2022, 1, 15, 21, 0, 0, 0, time.UTC).UnixMilli(), "validity.start")
	assertEqual(t, result.ValidityDated.End, time.Date(2022, 1, 16, 18, 0, 0, 0, time.UTC).UnixMilli(), "validity.end")
	assertEqual(t, len(result.Trends), 4, "trends length")

	if result.Amendment == nil || !*result.Amendment {
		t.Error("amendment should be true")
	}

	// Trend validity dates
	trend0 := result.Trends[0]
	if trend0.ValidityStartMs == nil {
		t.Fatal("trend0 ValidityStartMs should be set")
	}
	assertEqual(t, *trend0.ValidityStartMs, time.Date(2022, 1, 15, 21, 0, 0, 0, time.UTC).UnixMilli(), "trend0.start")
	if trend0.ValidityEndMs == nil {
		t.Fatal("trend0 ValidityEndMs should be set")
	}
	assertEqual(t, *trend0.ValidityEndMs, time.Date(2022, 1, 15, 23, 0, 0, 0, time.UTC).UnixMilli(), "trend0.end")

	// FM trend
	trend1 := result.Trends[1]
	if trend1.Type != WeatherChangeTypeFM {
		t.Error("trend1 should be FM")
	}
	if trend1.ValidityStartMs == nil {
		t.Fatal("trend1 ValidityStartMs should be set")
	}
	assertEqual(t, *trend1.ValidityStartMs, time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC).UnixMilli(), "trend1.start")
	// FM should NOT have ValidityEndMs
	if trend1.ValidityEndMs != nil {
		t.Error("trend1 ValidityEndMs should be nil for FM")
	}
}

func TestParseTAFDatedNextForecastByRemark(t *testing.T) {
	code := "TAF CYVR 152340Z 1600/1706 29015KT P6SM FEW015 FM162200 28010KT P6SM SKC RMK NXT FCST BY 160300Z"
	result, err := ParseTAFDated(code, time.Date(2022, 10, 22, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}

	if len(result.Trends) > 0 && len(result.Trends[0].Remarks) > 0 {
		remark := result.Trends[0].Remarks[0]
		assertEqual(t, remark.Type, RemarkTypeNextForecastBy, "remark.type")
		if remark.Date != nil {
			assertEqual(t, *remark.Date, time.Date(2022, 10, 16, 3, 0, 0, 0, time.UTC).UnixMilli(), "remark.date")
		} else {
			t.Error("remark.date should be set")
		}
	} else {
		t.Error("expected trends[0].remarks[0]")
	}
}

func TestParseTAFAsForecast(t *testing.T) {
	code := "TAF AMD KMSN 152044Z 1521/1618 24009G16KT P6SM SCT100 TEMPO 1521/1523 BKN100 FM160000 27008KT P6SM SCT150 FM160300 30006KT P6SM FEW230 FM161700 30011G18KT P6SM BKN050"
	fc, err := ParseTAFAsForecast(code, time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFAsForecast error: %v", err)
	}

	assertEqual(t, len(fc.Forecast), 5, "forecast length")
	assertEqual(t, fc.Station, "KMSN", "station")
	assertEqual(t, fc.Issued, time.Date(2022, 1, 15, 20, 44, 0, 0, time.UTC).UnixMilli(), "issued")
}

func TestParseTAFAsForecastWithoutDeliveryTime(t *testing.T) {
	code := "TAF KNBC 1215/1315 27010KT 9999 SCT010 BKN080 QNH2992INS TEMPO 1218/1300 25010G20KT 4800 TSRA BR BKN010CB BECMG 1300/1302 30015KT 6000 SHRA BR BKN015 QNH2998INS FM130430 04012KT 9999 NSW SCT020 BKN050 QNH2991INS T30/1219Z T22/1309Z"
	fc, err := ParseTAFAsForecast(code, time.Date(2022, 8, 12, 14, 57, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFAsForecast error: %v", err)
	}

	// No delivery time means issued should equal the passed date
	assertEqual(t, fc.Issued, time.Date(2022, 8, 12, 14, 57, 0, 0, time.UTC).UnixMilli(), "issued")
	assertEqual(t, fc.Station, "KNBC", "station")
	assertEqual(t, len(fc.Forecast), 4, "forecast length")
}

func TestParseTAFStationFM(t *testing.T) {
	code := "TAF FMMI 082300Z 0900/1006 16006KT 9999 FEW017 BKN020 PROB30 TEMPO 0908/0916 4500 RADZ BECMG 0909/0911 10010KT BECMG 0918/0920 16006KT"
	result, err := ParseTAFDated(code, time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatalf("ParseTAFDated error: %v", err)
	}
	assertEqual(t, result.Station, "FMMI", "station")
}
