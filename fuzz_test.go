package metartafparser

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// FuzzParseMetar fuzzes ParseMetar with raw METAR strings.
func FuzzParseMetar(f *testing.F) {
	seeds := []string{
		"KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992 RMK AO2",
		"LFPG 170830Z 00000KT 0350 R27L/0375N R09R/0175N R26R/0500D R08L/0400N R26L/0275D R08R/0250N R27R/0300N R09L/0200N FG SCT000 M01/M01 Q1026 NOSIG",
		"CYWG 172000Z 30015G25KT 3/4SM R36/4000FT/D -SN BLSN BKN008 OVC040 M05/M08 A2992 REFZRA WS RWY36 RMK SF5NS3 SLP134",
		"AUTO LSZL 061950Z 10002KT 9999NDV NCD 01/M00 Q1015 RMK=",
		"LFBG 081130Z AUTO 23012KT 9999 SCT022 BKN072 BKN090 22/16 Q1011 TEMPO 26015G25KT 3000 TSRA SCT025CB BKN050",
		"LFPG 212030Z 03003KT CAVOK 09/06 Q1031 NOSIG",
		"KTTN 051853Z 04011KT 1 1/2SM VCTS SN FZFG BKN003 OVC010 M02/M02 A3006 RMK AO2 TSB40 SLP176 P0002 T10171017=",
		"SUMU 070520Z M1/4SM",
		"SUMU 070520Z P6SM",
		"SUMU 070520Z 3 1/4SM",
		"KATL 270200Z 00000MPS",
		"AGGH 140340Z 05010KT 9999 TS FEW020 SCT021CB BKN300 32/26 Q1010",
		"SVMC 211703Z AUTO NIL",
		"METAR SUMU 070520Z 3 1/4SM",
		"SPECI SUMU 070520Z 3 1/4SM",
		"UUWW 151030Z 34002MPS CAVOK 14/02 Q1026 R01/000070 NOSIG",
		"VIDP 270200Z 00000MPS 0050",
		"ENLK 081350Z 26026G40 240V300 9999 VCSH FEW025 BKN030 02/M01 Q0996",
		"EGLL 231250Z 14012G22KT 2000 +TSRA FG BKN008 SCT025CB OVC050 18/17 Q1010 TEMPO 1000 -SHRA RMK QFE998",
		"CYVM 282100Z 36028G36KT 1SM -SN DRSN VCBLSN OVC008 M03/M04 A2935 RMK SN2ST8 LAST STFFD OBS/NXT 291200UTC SLP940",
		"KATL 022045Z 00000KT 10SM SCT120 00/M08 A2996 R12/1000V1200U",
		"LFRM 081630Z AUTO 30007KT 260V360 9999 24/15 Q1008 TEMPO SHRA BECMG SKC",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, input string) {
		result, err := ParseMetar(input, nil)
		if err == nil {
			if result.Station == "" {
				t.Error("station must not be empty on successful parse")
			}
			_, jsonErr := json.Marshal(result)
			if jsonErr != nil {
				t.Errorf("json.Marshal failed: %v", jsonErr)
			}
		}
	})
}

// FuzzParseTAF fuzzes ParseTAF with raw TAF strings.
func FuzzParseTAF(f *testing.F) {
	seeds := []string{
		"TAF KMSN 142325Z 1500/1524",
		"TAF AMD KMSN 152044Z 1521/1618 24009G16KT P6SM SCT100 TEMPO 1521/1523 BKN100 FM160000 27008KT P6SM SCT150 FM160300 30006KT P6SM FEW230 FM161700 30011G18KT P6SM BKN050",
		"TAF ESGG 260830Z 2609/2709 02009KT 3000 BR BKN003 BECMG 2609/2611 9999 NSW FEW015",
		"TAF EHAM 311720Z 3118/0124 12012KT CAVOK BECMG 3121/3124 17017KT",
		"TAF SBPJ 221450Z 2218/2318 21006KT 8000 SCT030 FEW040TCU TN25/2309Z TX34/2316Z BECMG 2221/2223 VRB03KT FEW030 BECMG 2302/2304 16003KT 5000 FU RMK PGU BECMG 2313/2315 23005KT SCT030 FEW040TCU",
		"TAF FMMI 082300Z 0900/1006 16006KT 9999 FEW017 BKN020 PROB30 TEMPO 0908/0916 4500 RADZ BECMG 0909/0911 10010KT BECMG 0918/0920 16006KT",
		"TAF CYVR 152340Z 1600/1706 29015KT P6SM FEW015 FM162200 28010KT P6SM SKC RMK NXT FCST BY 160300Z",
		"TAF KNBC 1215/1315 27010KT 9999 SCT010 BKN080 QNH2992INS TEMPO 1218/1300 25010G20KT 4800 TSRA BR BKN010CB BECMG 1300/1302 30015KT 6000 SHRA BR BKN015 QNH2998INS FM130430 04012KT 9999 NSW SCT020 BKN050 QNH2991INS T30/1219Z T22/1309Z",
		"TAF KMSN 302325Z 0100/0124",
		"TAF COR EDDS 201148Z 2012/2112 31010KT CAVOK BECMG 2018/2021 33004KT BECMG 2106/2109 07005KT",
		"TAF KMKE 011530 0116/0218 WS020/24045KT FM010200 17005KT P6SM SKC WS020/23055KT",
		"TAF KLSV 222300Z 2223/2405 21020G35KT 8000 BLDU BKN160 530009 QNH2941INS",
		"TAF VTBD 281000Z 2812/2912 CNL=",
		"TAF AMD CZBF 300939Z 3010/3022 VRB03KT 6SM -SN OVC015 TEMPO 3010/3012 11/2SM -SN OVC009 \nFM301200 10008KT 2SM -SN OVC010 TEMPO 3012/3022 3/4SM -SN VV007 RMK FCST BASED ON AUTO OBS. NXT FCST BY 301400Z",
		"TAF AMD KGWO 161553Z 1616/1712 21005KT 4SM -TSRA BR SCT010 OVC070CB",
		"TAF YWLM 270209Z 2703/2800 30014KT 9999 -SHRA NSC FM270400 28007KT 9999 -SHRA SCT040 FM270700 03010KT 9999 -SHRA SCT040 FM271200 30008KT CAVOK FM272100 29014KT CAVOK INTER 2703/2709 30018G30KT 5000 SHRA SCT015 BKN040 FEW040TCU PROB30 INTER 2704/2709 VRB25G45KT 2000 TSRAGR BKN010 SCT040CB",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, input string) {
		result, err := ParseTAF(input, nil)
		if err == nil {
			if result.Station == "" {
				t.Error("station must not be empty on successful parse")
			}
			_, jsonErr := json.Marshal(result)
			if jsonErr != nil {
				t.Errorf("json.Marshal failed: %v", jsonErr)
			}
		}
	})
}

// FuzzParseMetarDated fuzzes ParseMetarDated with raw METAR strings and issue timestamps.
//
//nolint:dupl
func FuzzParseMetarDated(f *testing.F) {
	type datedSeed struct {
		input  string
		issued int64
	}
	seeds := []datedSeed{
		{"LFPG 161430Z 24015G25KT 5000 1100w", time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992 RMK AO2", time.Date(2024, 6, 14, 8, 53, 0, 0, time.UTC).UnixMilli()},
		{"LFPG 170830Z 00000KT 0350 R27L/0375N NOSIG", time.Date(2022, 1, 17, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"SUMU 070520Z P6SM", time.Date(2024, 6, 7, 5, 20, 0, 0, time.UTC).UnixMilli()},
		{"SVMC 211703Z AUTO NIL", time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC).UnixMilli()},
	}
	for _, s := range seeds {
		f.Add(s.input, s.issued)
	}
	f.Fuzz(func(t *testing.T, input string, issued int64) {
		result, err := ParseMetarDated(input, time.UnixMilli(issued).UTC(), nil)
		if err != nil {
			return
		}
		if result.Station == "" {
			t.Error("station must not be empty on successful parse")
		}
		_, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			t.Errorf("json.Marshal failed: %v", jsonErr)
		}
	})
}

// FuzzParseTAFDated fuzzes ParseTAFDated with raw TAF strings and issue timestamps.
//
//nolint:dupl
func FuzzParseTAFDated(f *testing.F) {
	type tafDatedSeed struct {
		input  string
		issued int64
	}
	seeds := []tafDatedSeed{
		{"TAF AMD KMSN 152044Z 1521/1618 24009G16KT P6SM SCT100 TEMPO 1521/1523 BKN100", time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"TAF KMSN 302325Z 0100/0124", time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"TAF EHAM 311720Z 3118/0124 12012KT CAVOK BECMG 3121/3124 17017KT", time.Date(2022, 4, 15, 21, 46, 0, 0, time.UTC).UnixMilli()},
		{"TAF FMMI 082300Z 0900/1006 16006KT 9999 FEW017 BKN020 PROB30 TEMPO 0908/0916 4500 RADZ", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"TAF CYVR 152340Z 1600/1706 29015KT P6SM FEW015 FM162200 28010KT P6SM SKC RMK NXT FCST BY 160300Z", time.Date(2022, 10, 22, 0, 0, 0, 0, time.UTC).UnixMilli()},
	}
	for _, s := range seeds {
		f.Add(s.input, s.issued)
	}
	f.Fuzz(func(t *testing.T, input string, issued int64) {
		result, err := ParseTAFDated(input, time.UnixMilli(issued).UTC(), nil)
		if err != nil {
			return
		}
		if result.Station == "" {
			t.Error("station must not be empty on successful parse")
		}
		_, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			t.Errorf("json.Marshal failed: %v", jsonErr)
		}
	})
}

// FuzzParseTAFAsForecast fuzzes ParseTAFAsForecast with raw TAF strings and issue timestamps.
func FuzzParseTAFAsForecast(f *testing.F) {
	type seed struct {
		input  string
		issued int64
	}
	seeds := []seed{
		{"TAF AMD KMSN 152044Z 1521/1618 24009G16KT P6SM SCT100 TEMPO 1521/1523 BKN100 FM160000 27008KT P6SM SCT150 FM160300 30006KT P6SM FEW230 FM161700 30011G18KT P6SM BKN050", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"TAF KNBC 1215/1315 27010KT 9999 SCT010 BKN080 QNH2992INS TEMPO 1218/1300 25010G20KT 4800 TSRA BR BKN010CB BECMG 1300/1302 30015KT 6000 SHRA BR BKN015 QNH2998INS FM130430 04012KT 9999 NSW SCT020 BKN050 QNH2991INS T30/1219Z T22/1309Z", time.Date(2022, 8, 12, 14, 57, 0, 0, time.UTC).UnixMilli()},
		{"TAF KMSN 142325Z 1500/1524 25014G30KT P6SM VCSH SCT035 BKN070 TEMPO 1500/1502 6SM -SHRASN BKN035 FM150100 25012G25KT P6SM VCSH SCT040 BKN070 FM150300 26011G21KT P6SM SCT080", time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"TAF ESGG 260830Z 2609/2709 02009KT 3000 BR BKN003 BECMG 2609/2611 9999 NSW FEW015", time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC).UnixMilli()},
		{"TAF SBPJ 221450Z 2218/2318 21006KT 8000 SCT030 FEW040TCU TN25/2309Z TX34/2316Z BECMG 2221/2223 VRB03KT FEW030 BECMG 2302/2304 16003KT 5000 FU RMK PGU BECMG 2313/2315 23005KT SCT030 FEW040TCU", time.Date(2022, 10, 22, 0, 0, 0, 0, time.UTC).UnixMilli()},
	}
	for _, s := range seeds {
		f.Add(s.input, s.issued)
	}
	f.Fuzz(func(t *testing.T, input string, issued int64) {
		result, err := ParseTAFAsForecast(input, time.UnixMilli(issued).UTC(), nil)
		if err == nil {
			if result.Station == "" {
				t.Error("station must not be empty on successful parse")
			}
			_, jsonErr := json.Marshal(result)
			if jsonErr != nil {
				t.Errorf("json.Marshal failed: %v", jsonErr)
			}
		}
	})
}

// FuzzGetCompositeForecastForDate fuzzes GetCompositeForecastForDate with
// raw TAF strings, issue timestamps, and target timestamps.
func FuzzGetCompositeForecastForDate(f *testing.F) {
	type seed struct {
		input  string
		issued int64
		target int64
	}
	seeds := []seed{
		{
			"TAF KMSN 142325Z 1500/1524 25014G30KT P6SM VCSH SCT035 BKN070 TEMPO 1500/1502 6SM -SHRASN BKN035 FM150100 25012G25KT P6SM VCSH SCT040 BKN070 FM150300 26011G21KT P6SM SCT080",
			time.Date(2022, 4, 29, 0, 0, 0, 0, time.UTC).UnixMilli(),
			time.Date(2022, 4, 15, 0, 30, 0, 0, time.UTC).UnixMilli(),
		},
		{
			"TAF KMSN 142325Z 1500/1524",
			time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC).UnixMilli(),
			time.Date(2022, 4, 15, 12, 0, 0, 0, time.UTC).UnixMilli(),
		},
		{
			"TAF EHAM 311720Z 3118/0124 12012KT CAVOK BECMG 3121/3124 17017KT",
			time.Date(2022, 4, 15, 21, 46, 0, 0, time.UTC).UnixMilli(),
			time.Date(2022, 4, 15, 22, 0, 0, 0, time.UTC).UnixMilli(),
		},
	}
	for _, s := range seeds {
		f.Add(s.input, s.issued, s.target)
	}
	f.Fuzz(func(t *testing.T, input string, issued int64, target int64) {
		fc, err := ParseTAFAsForecast(input, time.UnixMilli(issued).UTC(), nil)
		if err != nil {
			return
		}
		_, err = GetCompositeForecastForDate(time.UnixMilli(target).UTC(), fc)
		if err != nil {
			var unexpectedErr *UnexpectedParseError
			if errors.As(err, &unexpectedErr) {
				t.Errorf("unexpected parse error: %v", err)
			}
			var tsErr *TimestampOutOfBoundsError
			if !errors.As(err, &tsErr) && !errors.As(err, &unexpectedErr) {
				t.Errorf("unexpected error type: %T: %v", err, err)
			}
		}
	})
}
