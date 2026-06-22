package metartafparser

// Container holds weather data fields used in METAR and TAF reports. Some fields
// (Icing, Turbulence) apply only to TAF reports.
type Container struct {
	// Surface wind conditions.
	Wind *Wind `json:"wind,omitempty"`
	// Prevailing visibility.
	Visibility *Visibility `json:"visibility,omitempty"`
	// Vertical visibility in feet (e.g. VV005).
	VerticalVisibility *int `json:"verticalVisibility,omitempty"`
	// Low-level wind shear.
	WindShear *WindShear `json:"windShear,omitempty"`
	// Ceiling and visibility OK (CAVOK).
	Cavok *bool `json:"cavok,omitempty"`
	// Raw remark string (the text after RMK).
	Remark *string `json:"remark,omitempty"`
	// Decoded remark entries.
	Remarks []Remark `json:"remarks"`
	// Cloud layers.
	Clouds []Cloud `json:"clouds"`
	// Present weather phenomena.
	WeatherConditions []WeatherCondition `json:"weatherConditions"`
	// Turbulence layers (TAF only).
	Turbulence []Turbulence `json:"turbulence,omitempty"`
	// Icing layers (TAF only).
	Icing []Icing `json:"icing,omitempty"`
}

// Wind represents surface wind direction, speed, gust, and directional variation.
type Wind struct {
	// Wind direction in degrees (0-360). Nil indicates variable (VRB).
	Degrees *int `json:"degrees,omitempty"`
	// Gust speed. Nil means no gust reported.
	Gust *int `json:"gust,omitempty"`
	// Minimum directional variation in degrees.
	MinVariation *int `json:"minVariation,omitempty"`
	// Maximum directional variation in degrees.
	MaxVariation *int `json:"maxVariation,omitempty"`
	// Cardinal direction (e.g. N, NE, VRB).
	Direction Direction `json:"direction"`
	// Unit of speed (knots, meters/second, etc.).
	Unit SpeedUnit `json:"unit"`
	// Sustained wind speed.
	Speed int `json:"speed"`
}

// WindShear represents low-level wind shear with direction, speed, gust, and height.
type WindShear struct {
	// Wind direction in degrees.
	Degrees *int `json:"degrees,omitempty"`
	// Gust speed.
	Gust *int `json:"gust,omitempty"`
	// Cardinal direction.
	Direction Direction `json:"direction"`
	// Unit of speed.
	Unit SpeedUnit `json:"unit"`
	// Sustained wind speed.
	Speed int `json:"speed"`
	// Height of wind shear in feet.
	Height int `json:"height"`
}

// Distance represents a measured distance with unit, value, and an optional comparison indicator.
type Distance struct {
	// Comparison indicator (less than, greater than).
	Indicator *ValueIndicator `json:"indicator,omitempty"`
	// Unit of measurement (meters, statute miles, etc.).
	Unit DistanceUnit `json:"unit"`
	// Numeric distance value.
	Value float64 `json:"value"`
	// No direction variation (NDV) flag.
	Ndv bool `json:"ndv,omitempty"`
}

// VisibilityMin captures the minimum visibility value and its associated direction.
type VisibilityMin struct {
	// Direction of minimum visibility.
	Direction Direction `json:"direction"`
	// Minimum visibility value.
	Value int `json:"value"`
}

// Visibility represents prevailing visibility, optionally with a minimum and direction.
type Visibility struct {
	Distance

	// Minimum visibility and its direction.
	Min *VisibilityMin `json:"min,omitempty"`
}

// WeatherCondition describes a weather phenomenon with intensity, descriptive qualifier, and phenomena.
type WeatherCondition struct {
	// Intensity (light, moderate, heavy, in vicinity).
	Intensity *Intensity `json:"intensity,omitempty"`
	// Descriptive qualifier (showers, thunderstorms, etc.).
	Descriptive *Descriptive `json:"descriptive,omitempty"`
	// Weather phenomena (rain, snow, drizzle, etc.).
	Phenomena []Phenomenon `json:"phenomena"`
}

func isWeatherConditionValid(wc WeatherCondition) bool {
	if len(wc.Phenomena) != 0 {
		return true
	}
	if wc.Descriptive != nil && *wc.Descriptive == DescriptiveThunderstorm {
		return true
	}
	if wc.Intensity != nil && *wc.Intensity == IntensityInVicinity &&
		wc.Descriptive != nil && *wc.Descriptive == DescriptiveShowers {
		return true
	}
	return false
}

// Altimeter represents the atmospheric pressure adjusted to sea level.
type Altimeter struct {
	// Unit of pressure (inches of mercury, hectopascals).
	Unit AltimeterUnit `json:"unit"`
	// Pressure value.
	Value float64 `json:"value"`
}

// Temperature represents a temperature value at a given day and hour.
type Temperature struct {
	// Temperature in degrees Celsius.
	Temperature float64 `json:"temperature"`
	// Day of the month (1-31).
	Day int `json:"day"`
	// Hour in UTC (0-23).
	Hour int `json:"hour"`
}

// Cloud describes a cloud layer with quantity, height, and optional type.
type Cloud struct {
	// Cloud base height above airport in feet.
	Height *int `json:"height,omitempty"`
	// Cloud type (cumulonimbus, towering cumulus, etc.).
	Type *CloudType `json:"type,omitempty"`
	// Secondary cloud type.
	SecondaryType *CloudType `json:"secondaryType,omitempty"`
	// Cloud coverage quantity (FEW, SCT, BKN, OVC, etc.).
	Quantity CloudQuantity `json:"quantity"`
}

// RunwayInfoRange describes runway visual range (RVR) including runway designator, minimum, maximum, trend, and unit.
type RunwayInfoRange struct {
	// Maximum range value.
	MaxRange *int `json:"maxRange,omitempty"`
	// Comparison indicator (less than, greater than).
	Indicator *ValueIndicator `json:"indicator,omitempty"`
	// Trend (upward, downward, no change).
	Trend *RunwayInfoTrend `json:"trend,omitempty"`
	// Runway designator (e.g. "07L").
	Name string `json:"name"`
	// Unit of measurement.
	Unit RunwayInfoUnit `json:"unit"`
	// Minimum range value.
	MinRange int `json:"minRange"`
}

// RunwayInfoDeposit describes runway surface condition: deposit type, coverage, thickness, and braking capacity.
type RunwayInfoDeposit struct {
	// Runway designator.
	Name string `json:"name"`
	// Type of deposit (dry, wet, snow, etc.).
	DepositType *DepositType `json:"depositType,omitempty"`
	// Extent of deposit coverage.
	Coverage *DepositCoverage `json:"coverage,omitempty"`
	// Thickness of deposit in mm.
	Thickness string `json:"thickness,omitempty"`
	// Braking capacity/code.
	BrakingCapacity string `json:"brakingCapacity,omitempty"`
}

// RunwayInfo aggregates runway visual range and/or surface deposit information for a single runway.
type RunwayInfo struct {
	// Runway visual range information.
	Range *RunwayInfoRange `json:"range,omitempty"`
	// Runway surface deposit information.
	Deposit *RunwayInfoDeposit `json:"deposit,omitempty"`
}

// Icing describes an icing layer with intensity, base height, and depth.
type Icing struct {
	// Intensity of icing (none, light, moderate, severe).
	Intensity IcingIntensity `json:"intensity"`
	// Base height of icing layer in feet.
	BaseHeight int `json:"baseHeight"`
	// Depth of icing layer in feet.
	Depth int `json:"depth"`
}

// Turbulence describes a turbulence layer with intensity, base height, and depth.
type Turbulence struct {
	// Intensity of turbulence (none, light, moderate, severe).
	Intensity TurbulenceIntensity `json:"intensity"`
	// Base height of turbulence layer in feet.
	BaseHeight int `json:"baseHeight"`
	// Depth of turbulence layer in feet.
	Depth int `json:"depth"`
}

// MetarTrendTime represents a time associated with a METAR trend, with an indicator (AT, FM, TL).
type MetarTrendTime struct {
	// Hour of the trend time.
	Hour *int `json:"hour,omitempty"`
	// Minute of the trend time.
	Minute *int `json:"minute,omitempty"`
	// Time indicator type (AT, FM, TL).
	Type TimeIndicator `json:"type"`
}

// MetarTrend represents a change in weather conditions within a METAR (BECMG, TEMPO, etc.).
type MetarTrend struct {
	Container

	// Type of weather change (BECMG, TEMPO, etc.).
	Type WeatherChangeType `json:"type"`
	// Raw trend string as reported.
	Raw string `json:"raw"`
	// Trend time markers.
	Times []MetarTrendTime `json:"times"`
}

// Validity defines the start and end period for a TAF forecast segment.
type Validity struct {
	// Start day of month.
	StartDay int `json:"startDay"`
	// Start hour in UTC.
	StartHour int `json:"startHour"`
	// End hour in UTC.
	EndHour int `json:"endHour"`
	// End day of month.
	EndDay int `json:"endDay"`
}

// FMValidity marks the exact start time (day, hour, minute) of a FROM (FM) TAF trend.
type FMValidity struct {
	// Start day of month.
	StartDay int `json:"startDay"`
	// Start hour in UTC.
	StartHour int `json:"startHour"`
	// Start minute.
	StartMinutes int `json:"startMinutes"`
}

// ValidityUnion is a tagged union of Validity and FMValidity.
type ValidityUnion interface {
	validityUnion()
}

func (Validity) validityUnion()   {}
func (FMValidity) validityUnion() {}

// BaseTAFTrend is the common structure shared by all TAF trend types.
type BaseTAFTrend struct {
	Container

	// Validity period (Validity or FMValidity).
	Validity ValidityUnion `json:"validity"`
	// Probability of occurrence (for PROB groups).
	Probability *int `json:"probability,omitempty"`
	// Absolute start time in Unix milliseconds (after date hydration).
	ValidityStartMs *int64 `json:"validityStartMs,omitempty"`
	// Absolute end time in Unix milliseconds (after date hydration).
	ValidityEndMs *int64 `json:"validityEndMs,omitempty"`
	// Weather change type (FM, BECMG, TEMPO, PROB, etc.).
	Type WeatherChangeType `json:"type"`
	// Raw trend string as reported.
	Raw string `json:"raw"`
}

// TAFTrend is a concrete TAF trend segment embedding BaseTAFTrend.
type TAFTrend struct {
	BaseTAFTrend
}

// Metar represents a fully parsed METAR (Meteorological Aerodrome Report) message.
type Metar struct {
	Container

	// Type of METAR report (METAR, SPECI).
	Type *MetarType `json:"type,omitempty"`
	// ICAO station identifier (e.g. KLAX).
	Station string `json:"station"`
	// Day of month (1-31).
	Day *int `json:"day,omitempty"`
	// Hour in UTC (0-23).
	Hour *int `json:"hour,omitempty"`
	// Minute (0-59).
	Minute *int `json:"minute,omitempty"`
	// Raw report message text.
	Message string `json:"message"`
	// Temperature in degrees Celsius.
	Temperature *float64 `json:"temperature,omitempty"`
	// Dew point in degrees Celsius.
	DewPoint *float64 `json:"dewPoint,omitempty"`
	// Altimeter (barometric pressure).
	Altimeter *Altimeter `json:"altimeter,omitempty"`
	// No significant weather changes (NOSIG).
	Nosig *bool `json:"nosig,omitempty"`
	// Runway visual range and deposit information.
	RunwaysInfo []RunwayInfo `json:"runwaysInfo"`
	// METAR trends (BECMG, TEMPO groups).
	Trends []MetarTrend `json:"trends"`
	// Report is an amendment (AMD).
	Amendment *bool `json:"amendment,omitempty"`
	// Report is fully automated (AO2 or AO1).
	Auto *bool `json:"auto,omitempty"`
	// Report is canceled.
	Canceled *bool `json:"canceled,omitempty"`
	// Report is a correction (COR).
	Corrected *bool `json:"corrected,omitempty"`
	// Report is nil (missing).
	Nil *bool `json:"nil,omitempty"`
}

// MetarDated wraps a Metar with the issue timestamp.
type MetarDated struct {
	Metar

	// Report issue time in Unix milliseconds.
	Issued int64 `json:"issued"`
}

// TAF represents a fully parsed TAF (Terminal Aerodrome Forecast) message.
type TAF struct {
	Container

	// ICAO station identifier.
	Station string `json:"station"`
	// Day of month.
	Day *int `json:"day,omitempty"`
	// Hour in UTC.
	Hour *int `json:"hour,omitempty"`
	// Minute.
	Minute *int `json:"minute,omitempty"`
	// Raw report message text.
	Message string `json:"message"`
	// Overall validity period of the TAF.
	Validity Validity `json:"validity"`
	// Forecast maximum temperature.
	MaxTemperature *Temperature `json:"maxTemperature,omitempty"`
	// Forecast minimum temperature.
	MinTemperature *Temperature `json:"minTemperature,omitempty"`
	// Raw text of the initial/base forecast segment.
	InitialRaw string `json:"initialRaw"`
	// TAF trend segments (FM, BECMG, TEMPO, PROB).
	Trends []TAFTrend `json:"trends"`
	// Report is an amendment.
	Amendment *bool `json:"amendment,omitempty"`
	// Report is automated.
	Auto *bool `json:"auto,omitempty"`
	// Report is canceled.
	Canceled *bool `json:"canceled,omitempty"`
	// Report is a correction.
	Corrected *bool `json:"corrected,omitempty"`
	// Report is nil.
	Nil *bool `json:"nil,omitempty"`
}

// ValidityDated is a dated variant of Validity with absolute Unix timestamps for the start and end.
type ValidityDated struct {
	// Start day of month.
	StartDay int `json:"startDay"`
	// Start hour in UTC.
	StartHour int `json:"startHour"`
	// End day of month.
	EndDay int `json:"endDay"`
	// End hour in UTC.
	EndHour int `json:"endHour"`
	// Absolute start time in Unix milliseconds.
	Start int64 `json:"start"`
	// Absolute end time in Unix milliseconds.
	End int64 `json:"end"`
}

// TAFDated wraps a TAF with the issue timestamp and dated validity.
type TAFDated struct {
	TAF

	// Report issue time in Unix milliseconds.
	Issued int64 `json:"issued"`
	// Dated validity with absolute timestamps.
	ValidityDated ValidityDated `json:"validity"`
}

// Remark represents a single decoded remark from a METAR/TAF report.
type Remark struct {
	// Hour associated with the remark.
	Hour *int `json:"hour,omitempty"`
	// Descriptive qualifier (e.g. light, moderate).
	Descriptive *Descriptive `json:"descriptive,omitempty"`
	// Movement direction (e.g. "E", "NE").
	Moving *string `json:"moving,omitempty"`
	// Minimum temperature or value.
	Min *float64 `json:"min,omitempty"`
	// Maximum temperature or value.
	Max *float64 `json:"max,omitempty"`
	// Wind speed in knots.
	Speed *int `json:"speed,omitempty"`
	// Wind direction in degrees.
	Degrees *int `json:"degrees,omitempty"`
	// Start hour of a time interval.
	StartHour *int `json:"startHour,omitempty"`
	// Start minute of a time interval.
	StartMinute *int `json:"startMinute,omitempty"`
	// End hour of a time interval.
	EndHour *int `json:"endHour,omitempty"`
	// End minute of a time interval.
	EndMinute *int `json:"endMinute,omitempty"`
	// Day of month.
	Day *int `json:"day,omitempty"`
	// Text description of the remark.
	Description *string `json:"description,omitempty"`
	// Location associated with the remark.
	Location *string `json:"location,omitempty"`
	// Precipitation depth in inches over the last hour.
	InchesLastHour *int `json:"inchesLastHour,omitempty"`
	// Generic numeric value.
	Value *float64 `json:"value,omitempty"`
	// Precipitation amount.
	Amount *float64 `json:"amount,omitempty"`
	// Related weather phenomenon.
	Phenomenon *Phenomenon `json:"phenomenon,omitempty"`
	// Timestamp (Unix milliseconds) associated with the remark.
	Date *int64 `json:"date,omitempty"`
	// Total snow/ice depth on the ground in inches.
	TotalDepth *int `json:"totalDepth,omitempty"`
	// Temperature value in degrees Celsius.
	Temperature *float64 `json:"temperature,omitempty"`
	// Dew point in degrees Celsius.
	DewPoint *float64 `json:"dewPoint,omitempty"`
	// Minute associated with the remark.
	Minute *int `json:"minute,omitempty"`
	// Direction string (e.g. "NE", "SE").
	Direction *string `json:"direction,omitempty"`
	// Type of remark (see RemarkType).
	Type RemarkType `json:"type"`
	// Raw remark string as reported.
	Raw string `json:"raw"`
}

// TemperatureDated extends Temperature with an absolute Unix timestamp.
type TemperatureDated struct {
	// Temperature in degrees Celsius.
	Temperature float64 `json:"temperature"`
	// Day of month.
	Day int `json:"day"`
	// Hour in UTC.
	Hour int `json:"hour"`
	// Absolute timestamp in Unix milliseconds.
	Date int64 `json:"date"`
}

// Forecast represents a single forecast period with weather conditions and timestamps.
type Forecast struct {
	// Ceiling and visibility OK.
	Cavok *bool `json:"cavok,omitempty"`
	// Raw remark string.
	Remark *string `json:"remark,omitempty"`
	// Weather change type (FM, BECMG, TEMPO, etc.).
	Type *WeatherChangeType `json:"type,omitempty"`
	// Time by which the change is expected in Unix milliseconds.
	By *int64 `json:"by,omitempty"`
	// Wind conditions during this forecast period.
	Wind *Wind `json:"wind,omitempty"`
	// Visibility during this forecast period.
	Visibility *Visibility `json:"visibility,omitempty"`
	// Low-level wind shear.
	WindShear *WindShear `json:"windShear,omitempty"`
	// Vertical visibility in feet.
	VerticalVisibility *int `json:"verticalVisibility,omitempty"`
	// Raw forecast segment string.
	Raw string `json:"raw"`
	// Decoded remarks.
	Remarks []Remark `json:"remarks"`
	// Cloud layers.
	Clouds []Cloud `json:"clouds"`
	// Weather phenomena.
	WeatherConditions []WeatherCondition `json:"weatherConditions"`
	// Turbulence layers.
	Turbulence []Turbulence `json:"turbulence,omitempty"`
	// Icing layers.
	Icing []Icing `json:"icing,omitempty"`
	// Start time in Unix milliseconds.
	Start int64 `json:"start"`
	// End time in Unix milliseconds.
	End int64 `json:"end"`
}

// ForecastContainer wraps a collection of Forecast segments with station and issue metadata.
type ForecastContainer struct {
	// Report is a correction.
	Corrected *bool `json:"corrected,omitempty"`
	// Report is an amendment.
	Amendment *bool `json:"amendment,omitempty"`
	// Report is automated.
	Auto *bool `json:"auto,omitempty"`
	// Report is canceled.
	Canceled *bool `json:"canceled,omitempty"`
	// Report is nil.
	Nil *bool `json:"nil,omitempty"`
	// Maximum temperature forecast.
	MaxTemperature *TemperatureDated `json:"maxTemperature,omitempty"`
	// Minimum temperature forecast.
	MinTemperature *TemperatureDated `json:"minTemperature,omitempty"`
	// Raw report message text.
	Message string `json:"message"`
	// ICAO station identifier.
	Station string `json:"station"`
	// Forecast segments.
	Forecast []Forecast `json:"forecast"`
	// Issue time in Unix milliseconds.
	Issued int64 `json:"issued"`
	// TAF start time in Unix milliseconds.
	Start int64 `json:"start"`
	// TAF end time in Unix milliseconds.
	End int64 `json:"end"`
}
