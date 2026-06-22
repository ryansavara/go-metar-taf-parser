# go-metar-taf-parser

[![Go Reference](https://pkg.go.dev/badge/github.com/ryansavara/metar-taf-parser.svg)](https://pkg.go.dev/github.com/ryansavara/metar-taf-parser)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ryansavara/metar-taf-parser)](https://github.com/ryansavara/metar-taf-parser)
[![GitHub Release](https://img.shields.io/github/v/release/ryansavara/metar-taf-parser)](https://github.com/ryansavara/metar-taf-parser/releases)
[![License](https://img.shields.io/github/license/ryansavara/metar-taf-parser)](LICENSE)

A Go library (Go 1.25.7) for parsing METAR and TAF aviation weather reports. This is a Go conversion of the [TypeScript metar-taf-parser](https://github.com/aeharding/metar-taf-parser).

```
go get github.com/ryansavara/metar-taf-parser
```

```go
import metartafparser "github.com/ryansavara/metar-taf-parser"
```

---

## Quick Start

```go
import metartafparser "github.com/ryansavara/metar-taf-parser"

// Parse a METAR
metar, err := metartafparser.ParseMetar("KLAX 140853Z 00000KT 10SM FEW010 14/12 A2992 RMK AO2 SLP132 T01440117", nil)

// Parse a TAF
taf, err := metartafparser.ParseTAF("KLAX 140520Z 1406/1512 05005KT P6SM FEW010", nil)

// Parse with date hydration (resolves day/hour to absolute timestamps)
issued := time.Date(2024, 6, 14, 8, 53, 0, 0, time.UTC)
metarDated, err := metartafparser.ParseMetarDated("KLAX 140853Z ...", issued, nil)

// Parse TAF with date hydration
tafDated, err := metartafparser.ParseTAFDated("KLAX 140520Z ...", issued, nil)

// Parse TAF as forecast segments
fc, err := metartafparser.ParseTAFAsForecast("KLAX 140520Z ...", issued, nil)

// Get prevailing + supplemental forecast for a specific time
cf, err := metartafparser.GetCompositeForecastForDate(someTime, fc)

// Use localized output
opts := &metartafparser.ParseOptions{Locale: metartafparser.DefaultLocale()}
metar, err := metartafparser.ParseMetar("...", opts)
```

---

## Parse Functions

| Function | Returns | Description |
|---|---|---|
| `ParseMetar(raw string, opts *ParseOptions)` | `(*Metar, error)` | Parse a raw METAR string |
| `ParseMetarDated(raw string, issued time.Time, opts *ParseOptions)` | `(*MetarDated, error)` | Parse METAR and resolve day/hour to absolute Unix timestamps |
| `ParseTAF(raw string, opts *ParseOptions)` | `(*TAF, error)` | Parse a raw TAF string |
| `ParseTAFDated(raw string, issued time.Time, opts *ParseOptions)` | `(*TAFDated, error)` | Parse TAF and resolve day/hour to absolute Unix timestamps |
| `ParseTAFAsForecast(raw string, issued time.Time, opts *ParseOptions)` | `(*ForecastContainer, error)` | Parse TAF into time-bound forecast segments |
| `GetCompositeForecastForDate(date time.Time, fc *ForecastContainer)` | `(*CompositeForecast, error)` | Get prevailing + supplemental forecasts for a specific timestamp |

---

## Data Model

### Top-Level Reports

```
Metar ──────────────────────────────────────── MetarDated
  ├─ Container (embedded)                       └─ Issued (int64 ms)
  ├─ Type             *MetarType
  ├─ Station          string
  ├─ Day/Hour/Minute  *int
  ├─ Temperature      *float64 °C
  ├─ DewPoint         *float64 °C
  ├─ Altimeter        *Altimeter
  ├─ RunwaysInfo      []RunwayInfo
  ├─ Trends           []MetarTrend
  ├─ Nosig            *bool
  └─ Amendment/Auto/Canceled/Corrected/Nil *bool

TAF ────────────────────────────────────────── TAFDated
  ├─ Container (embedded)                       ├─ Issued (int64 ms)
  ├─ Station          string                     └─ ValidityDated (absolute ms)
  ├─ Day/Hour/Minute  *int
  ├─ Validity         Validity
  ├─ MaxTemperature   *Temperature
  ├─ MinTemperature   *Temperature
  ├─ InitialRaw       string
  ├─ Trends           []TAFTrend
  └─ Amendment/Auto/Canceled/Corrected/Nil *bool
```

**`Metar`** (`types.go:279`) — Fully parsed METAR report.

**`MetarDated`** (`types.go:319`) — Wraps `Metar` with `Issued` (Unix ms) resolved from relative day/hour/minute.

**`TAF`** (`types.go:327`) — Fully parsed TAF report.

**`TAFDated`** (`types.go:379`) — Wraps `TAF` with `Issued` and `ValidityDated` (absolute Unix ms).

**`ForecastContainer`** (`types.go:493`) — TAF converted to time-bound forecast segments with absolute timestamps.

```
ForecastContainer
├─ Station       string
├─ Issued        int64 ms
├─ Start         int64 ms
├─ End           int64 ms
├─ Message       string
├─ Forecast      []Forecast
├─ MaxTemperature *TemperatureDated
├─ MinTemperature *TemperatureDated
└─ Amendment/Auto/Canceled/Corrected/Nil *bool
```

**`Forecast`** (`types.go:457`) — Single forecast period with weather conditions and time boundaries.

```
Forecast
├─ Type              *WeatherChangeType
├─ Start             int64 ms
├─ End               int64 ms
├─ By                *int64 ms    (BECMG deadline)
├─ Wind              *Wind
├─ Visibility        *Visibility
├─ WindShear         *WindShear
├─ VerticalVisibility *int
├─ Cavok             *bool
├─ Clouds            []Cloud
├─ WeatherConditions []WeatherCondition
├─ Turbulence        []Turbulence
├─ Icing             []Icing
├─ Remarks           []Remark
├─ Remark            *string
└─ Raw               string
```

**`CompositeForecast`** (`metartafparser.go:406`) — Prevailing forecast plus any supplemental (TEMPO/INTER) segments for a given timestamp.

```
CompositeForecast
├─ Prevailing   Forecast
└─ Supplemental []Forecast
```

**`ParseOptions`** (`metartafparser.go:9`) — Configuration with a single field:

```go
type ParseOptions struct {
    Locale Locale  // Use DefaultLocale() for English translations
}
```

### Weather Components

**`Container`** (`types.go:4`) — Shared weather fields embedded in `Metar`, `TAF`, and `BaseTAFTrend`.

```
Container
├─ Wind              *Wind
├─ Visibility        *Visibility
├─ VerticalVisibility *int
├─ WindShear         *WindShear
├─ Cavok             *bool
├─ Clouds            []Cloud
├─ WeatherConditions []WeatherCondition
├─ Turbulence        []Turbulence    (TAF only)
├─ Icing             []Icing         (TAF only)
├─ Remarks           []Remark
└─ Remark            *string
```

**`Wind`** (`types.go:30`) — Surface wind direction, speed, gust, and directional variation.

```go
type Wind struct {
    Degrees      *int          // 0–360 degrees (nil = VRB)
    Speed        int           // Sustained wind speed
    Gust         *int          // Gust speed (nil = no gust)
    MinVariation *int          // Directional variation start
    MaxVariation *int          // Directional variation end
    Direction    Direction     // Cardinal direction (N, NE, VRB, etc.)
    Unit         SpeedUnit     // KT, MPS, KM/H
}
```

**`WindShear`** (`types.go:48`) — Low-level wind shear with height.

```go
type WindShear struct {
    Degrees   *int
    Speed     int
    Gust      *int
    Height    int           // Feet
    Direction Direction
    Unit      SpeedUnit
}
```

**`Visibility`** (`types.go:84`) — Prevailing visibility, embedding `Distance`.

```go
type Visibility struct {
    Distance              // embedded Value, Unit, Indicator, Ndv
    Min *VisibilityMin    // Minimum visibility + direction
}
```

**`Distance`** (`types.go:64`) — Measured distance with unit and optional comparison indicator.

```go
type Distance struct {
    Value     float64
    Unit      DistanceUnit    // m or SM
    Indicator *ValueIndicator // P (>), M (<)
    Ndv       bool            // No directional variation
}
```

**`VisibilityMin`** (`types.go:76`) — Minimum visibility value and its direction.

**`Cloud`** (`types.go:134`) — Cloud layer with coverage quantity, base height, and type.

```go
type Cloud struct {
    Quantity      CloudQuantity // SKC, FEW, SCT, BKN, OVC, NSC
    Height        *int          // Base height in feet (nil = unknown)
    Type          *CloudType    // CB, TCU, CI, etc.
    SecondaryType *CloudType
}
```

**`WeatherCondition`** (`types.go:92`) — Weather phenomenon with intensity, descriptor, and phenomena list.

```go
type WeatherCondition struct {
    Intensity   *Intensity
    Descriptive *Descriptive
    Phenomena   []Phenomenon
}
```

**`Altimeter`** (`types.go:116`) — Barometric pressure adjusted to sea level.

```go
type Altimeter struct {
    Value float64
    Unit  AltimeterUnit    // inHg or hPa
}
```

### TAF-Specific Types

**`Validity`** (`types.go:226`) — Period start/end for a TAF forecast segment.

```go
type Validity struct {
    StartDay, StartHour, EndDay, EndHour int
}
```

**`ValidityDated`** (`types.go:363`) — `Validity` with absolute Unix millisecond timestamps.

```go
type ValidityDated struct {
    StartDay, StartHour, EndDay, EndHour int
    Start, End                           int64
}
```

**`FMValidity`** (`types.go:238`) — Exact start time (day, hour, minute) for a FROM trend.

```go
type FMValidity struct {
    StartDay, StartHour, StartMinutes int
}
```

**`ValidityUnion`** (`types.go:248`) — Interface satisfied by both `Validity` and `FMValidity`.

**`BaseTAFTrend`** (`types.go:256`) — Common trend structure with `Container`, validity, probability, and timestamps.

```go
type BaseTAFTrend struct {
    Container                        // Weather conditions for this trend
    Validity         ValidityUnion   // Validity or FMValidity
    Type             WeatherChangeType
    Probability      *int            // For PROB groups
    ValidityStartMs  *int64
    ValidityEndMs    *int64
    Raw              string
}
```

**`TAFTrend`** (`types.go:274`) — Concrete trend segment embedding `BaseTAFTrend`.

**`MetarTrend`** (`types.go:214`) — METAR weather change (BECMG, TEMPO) with time markers.

```go
type MetarTrend struct {
    Container
    Type  WeatherChangeType
    Times []MetarTrendTime
    Raw   string
}
```

**`MetarTrendTime`** (`types.go:204`) — Trend time marker (AT, FM, TL) with hour/minute.

**`Temperature`** (`types.go:124`) — Temperature value at a specific day and hour.

```go
type Temperature struct {
    Temperature float64
    Day         int
    Hour        int
}
```

**`TemperatureDated`** (`types.go:445`) — `Temperature` with an absolute Unix timestamp.

**`Icing`** (`types.go:184`) — In-flight icing layer.

```go
type Icing struct {
    Intensity  IcingIntensity  // 0–9 scale
    BaseHeight int             // Feet
    Depth      int             // Feet
}
```

**`Turbulence`** (`types.go:194`) — Turbulence layer.

```go
type Turbulence struct {
    Intensity  TurbulenceIntensity  // 0–9 scale or X (extreme)
    BaseHeight int                  // Feet
    Depth      int                  // Feet
}
```

### Runway Information

**`RunwayInfo`** (`types.go:176`) — Aggregates RVR and/or surface deposit data for a runway.

```go
type RunwayInfo struct {
    Range   *RunwayInfoRange
    Deposit *RunwayInfoDeposit
}
```

**`RunwayInfoRange`** (`types.go:146`) — Runway visual range with min, max, trend, and unit.

```go
type RunwayInfoRange struct {
    Name      string
    MinRange  int
    MaxRange  *int
    Unit      RunwayInfoUnit   // FT or m
    Indicator *ValueIndicator  // P (>) or M (<)
    Trend     *RunwayInfoTrend // U, D, N
}
```

**`RunwayInfoDeposit`** (`types.go:162`) — Runway surface condition: deposit type, coverage, thickness, braking capacity.

### Remarks

**`Remark`** (`types.go:389`) — A single decoded remark from the RMK section.

```go
type Remark struct {
    Type          RemarkType
    Raw           string
    Description   *string
    Value         *float64
    Amount        *float64
    Temperature   *float64
    DewPoint      *float64
    Hour          *int
    Minute        *int
    Day           *int
    Date          *int64        // Unix ms (for Next Forecast)
    Direction     *string
    Location      *string
    Moving        *string
    Speed         *int
    Degrees       *int
    StartHour     *int
    StartMinute   *int
    EndHour       *int
    EndMinute     *int
    // ... plus specialized fields for precipitation, snow, etc.
}
```

### Locale

**`Locale`** (`locale.go:9`) — Nested map-based i18n system.

```go
type Locale map[string]any
```

**`DefaultLocale()`** (`locale.go:47`) — Returns a `Locale` pre-populated with full English translations for all report fields. To customize, modify the returned map or build your own `Locale` from scratch.

---

## Enums

### Report Type

| Constant | Value |
|---|---|
| `MetarTypeMETAR` | `"METAR"` |
| `MetarTypeSPECI` | `"SPECI"` |

### Cloud Coverage (`CloudQuantity`)

| Constant | Value |
|---|---|
| `CloudQuantitySKC` | `"SKC"` — Sky clear |
| `CloudQuantityFEW` | `"FEW"` — Few (1–2 oktas) |
| `CloudQuantitySCT` | `"SCT"` — Scattered (3–4 oktas) |
| `CloudQuantityBKN` | `"BKN"` — Broken (5–7 oktas) |
| `CloudQuantityOVC` | `"OVC"` — Overcast (8 oktas) |
| `CloudQuantityNSC` | `"NSC"` — No significant clouds |

### Cloud Type (`CloudType`)

`CB` (Cumulonimbus), `TCU` (Towering cumulus), `CI` (Cirrus), `CC` (Cirrocumulus), `CS` (Cirrostratus), `AC` (Altocumulus), `AS` (Altostratus), `ST` (Stratus), `CU` (Cumulus), `NS` (Nimbostratus), `SC` (Stratocumulus).

### Weather Intensity (`Intensity`)

`IntensityLight`, `IntensityModerate`, `IntensityHeavy`, `IntensityInVicinity`.

### Weather Descriptor (`Descriptive`)

`DescriptiveShowers` (`"SH"`), `DescriptiveThunderstorm` (`"TS"`), `DescriptiveFreezing` (`"FZ"`), `DescriptiveShallow` (`"MI"`), `DescriptivePatches` (`"BC"`), `DescriptivePartial` (`"PR"`), `DescriptiveDrifting` (`"DR"`), `DescriptiveBlowing` (`"BL"`).

### Weather Phenomena (`Phenomenon`)

| Category | Values |
|---|---|
| Precipitation | `RA`, `DZ`, `SN`, `SG`, `PL`, `IC`, `GR`, `GS`, `UP` |
| Obscuration | `FG`, `BR`, `HZ`, `DU`, `FU`, `SA`, `PY`, `VA` |
| Other | `SQ`, `PO`, `DS`, `SS`, `FC`, `TS`, `NSW` |

### Time Indicators (`TimeIndicator`)

`TimeIndicatorAT` ("at"), `TimeIndicatorFM` ("from"), `TimeIndicatorTL` ("until").

### Weather Change Types (`WeatherChangeType`)

`WeatherChangeTypeFM`, `WeatherChangeTypeBECMG`, `WeatherChangeTypeTEMPO`, `WeatherChangeTypeINTER`, `WeatherChangeTypePROB`.

### Direction (`Direction`)

16 cardinal/intercardinal points (`N`, `NNE`, `NE`, `ENE`, `E`, `ESE`, `SE`, `SSE`, `S`, `SSW`, `SW`, `WSW`, `W`, `WNW`, `NW`, `NNW`) plus `VRB` (variable).

### Distance Unit (`DistanceUnit`)

`DistanceUnitMeters` (`"m"`), `DistanceUnitStatuteMiles` (`"SM"`).

### Speed Unit (`SpeedUnit`)

`SpeedUnitKnot` (`"KT"`), `SpeedUnitMetersPerSecond` (`"MPS"`), `SpeedUnitKilometersPerHour` (`"KM/H"`).

### Value Indicator (`ValueIndicator`)

`ValueIndicatorGreaterThan` (`"P"`), `ValueIndicatorLessThan` (`"M"`).

### Runway Visual Range Trend (`RunwayInfoTrend`)

`RunwayInfoTrendUprising` (`"U"`), `RunwayInfoTrendDecreasing` (`"D"`), `RunwayInfoTrendNoSignificantChange` (`"N"`).

### Runway Visual Range Unit (`RunwayInfoUnit`)

`RunwayInfoUnitFeet` (`"FT"`), `RunwayInfoUnitMeters` (`"m"`).

### Altimeter Unit (`AltimeterUnit`)

`AltimeterUnitInHg` (`"inHg"`), `AltimeterUnitHPa` (`"hPa"`).

### Icing Intensity (`IcingIntensity`)

Numeric scale 0–9: `None` (0) → `Light` (1) → `Moderate` (4–6) → `Severe` (7–9), with descriptors for rime/clear ice and cloud/precipitation.

### Turbulence Intensity (`TurbulenceIntensity`)

Numeric scale 0–9 plus `Extreme` (`"X"`): `None` (0) → `Light` (1) → `Moderate` (2–5) → `Severe` (6–9), with descriptors for clear air vs cloud and occasional vs frequent.

### Runway Deposit Type (`DepositType`)

`NotReported`, `ClearDry`, `Damp`, `WetWaterPatches`, `RimeFrostCovered`, `DrySnow`, `WetSnow`, `Slush`, `Ice`, `CompactedSnow`, `FrozenRidges`.

### Runway Deposit Coverage (`DepositCoverage`)

`None`, `NotReported`, `Less10`, `From11To25`, `From26To50`, `From51To100`.

### Remark Types (`RemarkType`)

50+ constants covering all standard remark categories: automation (`AO1`, `AO2`), pressure (`PRESFR`, `PRESRR`), tornadic activity, precipitation timing, hail size, snow, ceiling, visibility, temperature extremes, sea-level pressure, ice accretion, sunshine duration, and more.

---

## Error Handling

All parsing errors implement or embed the base `ParseError` type.

```go
// Base error type
type ParseError struct { ... }

// Input could not be parsed as a valid weather report
type InvalidWeatherStatementError struct { ... }

// Input appears to be a partial or incomplete TAF
type PartialWeatherStatementError struct {
    Part, Total int
    ...
}

// Unexpected error during parsing
type UnexpectedParseError struct { ... }

// Timestamp falls outside the report validity period
type TimestampOutOfBoundsError struct { ... }
```

**Constructors:**

| Function | Returns |
|---|---|
| `NewParseError(message string)` | `*ParseError` |
| `NewInvalidWeatherStatementError(cause any)` | `*InvalidWeatherStatementError` |
| `NewPartialWeatherStatementError(partialMessage string, part, total int)` | `*PartialWeatherStatementError` |
| `NewUnexpectedParseError(message string)` | `*UnexpectedParseError` |
| `NewTimestampOutOfBoundsError(message string)` | `*TimestampOutOfBoundsError` |

Error types form an unwrap chain: `PartialWeatherStatementError` → `InvalidWeatherStatementError` → `ParseError`. Use `errors.As` or `errors.Is` to check for specific error types.

---

## CLI

A command-line tool is available at `cmd/metar-taf-parser`:

```text
# Install
go install github.com/ryansavara/metar-taf-parser/cmd/metar-taf-parser@latest

# Parse as argument
metar-taf-parser "METAR KLAX 121653Z 27008KT 10SM FEW020 18/13 A2992"
metar-taf-parser "TAF KLAX 121720Z 1218/1324 27008KT 6SM HZ BKN020"

# Pipe via stdin
echo "KLAX 121653Z 27008KT 10SM FEW020 18/13 A2992" | metar-taf-parser
```

The CLI auto-detects METAR vs TAF, outputs indented JSON, and exits non-zero on parse failure.

### Cross-Platform Builds

Pre-built binaries for Linux, macOS, and Windows are attached to each [release](https://github.com/ryansavara/metar-taf-parser/releases). To build from source:

```text
git clone https://github.com/ryansavara/metar-taf-parser.git
cd metar-taf-parser
go build ./cmd/metar-taf-parser/
```

---

## Date Hydration

Reports use relative day/hour/minute values (e.g., `140853Z` = day 14, 08:53 UTC). The `Dated` variants resolve these to absolute Unix millisecond timestamps by searching ±1 month from the `issued` time and picking the closest match. This logic is used by:

- `ParseMetarDated` — resolves the report issue time
- `ParseTAFDated` — resolves issue time, validity period, and all trend timestamps
- `ParseTAFAsForecast` — same as `ParseTAFDated`, then converts to `Forecast` segments

BECMG forecasts inherit unset fields from the preceding forecast context when converted to segments.
